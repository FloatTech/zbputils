// Package control 控制插件的启用与优先级等
package control

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"

	ctrl "github.com/FloatTech/zbpctrl"

	"github.com/FloatTech/zbputils/img"
	"github.com/FloatTech/zbputils/process"

	// 图片输出
	"image"

	"github.com/Coloured-glaze/gg"
	"github.com/FloatTech/zbputils/file"
	"github.com/FloatTech/zbputils/img/text"
	"github.com/FloatTech/zbputils/img/writer"
)

var (
	// managers 每个插件对应的管理
	managers = ctrl.NewManager[*zero.Ctx]("data/control/plugins.db", 10*time.Second)
)

type mpic struct {
	kanban string      // 看板娘
	stat   string      // 启用状态
	stat2  bool        // 启用状态
	plugin string      // 插件名
	ttf1   string      // 字体1
	ttf2   string      // 字体2
	info   []string    // 插件信息
	im     image.Image // 图片
}

func newctrl(service string, o *ctrl.Options[*zero.Ctx]) zero.Rule {
	c := managers.NewControl(service, o)
	return func(ctx *zero.Ctx) bool {
		ctx.State["manager"] = c
		return c.Handler(uintptr(unsafe.Pointer(ctx)), ctx.Event.GroupID, ctx.Event.UserID)
	}
}

// Lookup 查找服务
func Lookup(service string) (*ctrl.Control[*zero.Ctx], bool) {
	return managers.Lookup(service)
}

func init() {
	process.NewCustomOnce(&managers).Do(func() {
		err := os.MkdirAll("data/Control", 0755)
		if err != nil {
			panic(err)
		}

		zero.OnCommandGroup([]string{
			"响应", "response", "沉默", "silence",
		}, zero.UserOrGrpAdmin).SetBlock(true).SecondPriority().Handle(func(ctx *zero.Ctx) {
			grp := ctx.Event.GroupID
			if grp == 0 {
				// 个人用户
				grp = -ctx.Event.UserID
			}
			var msg message.MessageSegment
			switch ctx.State["command"] {
			case "响应", "response":
				err := managers.Response(grp)
				if err == nil {
					msg = message.Text(zero.BotConfig.NickName[0], "将开始在此工作啦~")
				} else {
					msg = message.Text("ERROR: ", err)
				}
			case "沉默", "silence":
				err := managers.Silence(grp)
				if err == nil {
					msg = message.Text(zero.BotConfig.NickName[0], "将开始休息啦~")
				} else {
					msg = message.Text("ERROR: ", err)
				}
			default:
				msg = message.Text("ERROR: bad command\"", ctx.State["command"], "\"")
			}
			ctx.SendChain(msg)
		})

		zero.OnCommandGroup([]string{
			"全局响应", "allresponse", "全局沉默", "allsilence",
		}, zero.SuperUserPermission).SetBlock(true).SecondPriority().Handle(func(ctx *zero.Ctx) {
			var msg message.MessageSegment
			cmd := ctx.State["command"].(string)
			switch {
			case strings.Contains(cmd, "响应") || strings.Contains(cmd, "response"):
				err := managers.Response(0)
				if err == nil {
					msg = message.Text(zero.BotConfig.NickName[0], "将开始在全部位置工作啦~")
				} else {
					msg = message.Text("ERROR: ", err)
				}
			case strings.Contains(cmd, "沉默") || strings.Contains(cmd, "silence"):
				err := managers.Silence(0)
				if err == nil {
					msg = message.Text(zero.BotConfig.NickName[0], "将开始在未显式启用的位置休息啦~")
				} else {
					msg = message.Text("ERROR: ", err)
				}
			default:
				msg = message.Text("ERROR: bad command\"", cmd, "\"")
			}
			ctx.SendChain(msg)
		})

		zero.OnCommandGroup([]string{
			"启用", "enable", "禁用", "disable",
		}, zero.UserOrGrpAdmin).SetBlock(true).SecondPriority().Handle(func(ctx *zero.Ctx) {
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			service, ok := Lookup(model.Args)
			if !ok {
				ctx.SendChain(message.Text("没有找到指定服务!"))
				return
			}
			grp := ctx.Event.GroupID
			if grp == 0 {
				// 个人用户
				grp = -ctx.Event.UserID
			}
			if strings.Contains(model.Command, "启用") || strings.Contains(model.Command, "enable") {
				service.Enable(grp)
				if service.Options.OnEnable != nil {
					service.Options.OnEnable(ctx)
				} else {
					ctx.SendChain(message.Text("已启用服务: " + model.Args))
				}
			} else {
				service.Disable(grp)
				if service.Options.OnDisable != nil {
					service.Options.OnDisable(ctx)
				} else {
					ctx.SendChain(message.Text("已禁用服务: " + model.Args))
				}
			}
		})

		zero.OnCommandGroup([]string{
			"全局启用", "allenable", "全局禁用", "alldisable",
		}, zero.OnlyToMe, zero.SuperUserPermission).SetBlock(true).SecondPriority().Handle(func(ctx *zero.Ctx) {
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			service, ok := Lookup(model.Args)
			if !ok {
				ctx.SendChain(message.Text("没有找到指定服务!"))
				return
			}
			if strings.Contains(model.Command, "启用") || strings.Contains(model.Command, "enable") {
				service.Enable(0)
				ctx.SendChain(message.Text("已全局启用服务: " + model.Args))
			} else {
				service.Disable(0)
				ctx.SendChain(message.Text("已全局禁用服务: " + model.Args))
			}
		})

		zero.OnCommandGroup([]string{"还原", "reset"}, zero.UserOrGrpAdmin).SetBlock(true).SecondPriority().Handle(func(ctx *zero.Ctx) {
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			service, ok := Lookup(model.Args)
			if !ok {
				ctx.SendChain(message.Text("没有找到指定服务!"))
				return
			}
			grp := ctx.Event.GroupID
			if grp == 0 {
				// 个人用户
				grp = -ctx.Event.UserID
			}
			service.Reset(grp)
			ctx.SendChain(message.Text("已还原服务的默认启用状态: " + model.Args))
		})

		zero.OnCommandGroup([]string{
			"禁止", "ban", "允许", "permit",
		}, zero.AdminPermission).SetBlock(true).SecondPriority().Handle(func(ctx *zero.Ctx) {
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			args := strings.Split(model.Args, " ")
			if len(args) >= 2 {
				service, ok := Lookup(args[0])
				if !ok {
					ctx.SendChain(message.Text("没有找到指定服务!"))
					return
				}
				grp := ctx.Event.GroupID
				if grp == 0 {
					grp = -ctx.Event.UserID
				}
				msg := "**" + args[0] + "报告**"
				var members map[int64]struct{}
				issu := zero.SuperUserPermission(ctx)
				if !issu {
					lst := ctx.GetGroupMemberList(ctx.Event.GroupID).Array()
					members = make(map[int64]struct{}, len(lst))
					for _, m := range lst {
						members[m.Get("user_id").Int()] = struct{}{}
					}
				}
				if strings.Contains(model.Command, "允许") || strings.Contains(model.Command, "permit") {
					for _, usr := range args[1:] {
						uid, err := strconv.ParseInt(usr, 10, 64)
						if err == nil {
							if issu {
								service.Permit(uid, grp)
								msg += "\n+ 已允许" + usr
							} else {
								_, ok := members[uid]
								if ok {
									service.Permit(uid, grp)
									msg += "\n+ 已允许" + usr
								} else {
									msg += "\nx " + usr + " 不在本群"
								}
							}
						}
					}
				} else {
					for _, usr := range args[1:] {
						uid, err := strconv.ParseInt(usr, 10, 64)
						if err == nil {
							if issu {
								service.Ban(uid, grp)
								msg += "\n- 已禁止" + usr
							} else {
								_, ok := members[uid]
								if ok {
									service.Ban(uid, grp)
									msg += "\n- 已禁止" + usr
								} else {
									msg += "\nx " + usr + " 不在本群"
								}
							}
						}
					}
				}
				ctx.SendChain(message.Text(msg))
				return
			}
			ctx.SendChain(message.Text("参数错误!"))
		})

		zero.OnCommandGroup([]string{
			"全局禁止", "allban", "全局允许", "allpermit",
		}, zero.SuperUserPermission).SetBlock(true).SecondPriority().Handle(func(ctx *zero.Ctx) {
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			args := strings.Split(model.Args, " ")
			if len(args) >= 2 {
				service, ok := Lookup(args[0])
				if !ok {
					ctx.SendChain(message.Text("没有找到指定服务!"))
					return
				}
				msg := "**" + args[0] + "全局报告**"
				if strings.Contains(model.Command, "允许") || strings.Contains(model.Command, "permit") {
					for _, usr := range args[1:] {
						uid, err := strconv.ParseInt(usr, 10, 64)
						if err == nil {
							service.Permit(uid, 0)
							msg += "\n+ 已允许" + usr
						}
					}
				} else {
					for _, usr := range args[1:] {
						uid, err := strconv.ParseInt(usr, 10, 64)
						if err == nil {
							service.Ban(uid, 0)
							msg += "\n- 已禁止" + usr
						}
					}
				}
				ctx.SendChain(message.Text(msg))
				return
			}
			ctx.SendChain(message.Text("参数错误!"))
		})

		zero.OnCommandGroup([]string{
			"封禁", "block", "解封", "unblock",
		}, zero.SuperUserPermission).SetBlock(true).SecondPriority().Handle(func(ctx *zero.Ctx) {
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			args := strings.Split(model.Args, " ")
			if len(args) >= 1 {
				msg := "**报告**"
				if strings.Contains(model.Command, "解") || strings.Contains(model.Command, "un") {
					for _, usr := range args {
						uid, err := strconv.ParseInt(usr, 10, 64)
						if err == nil {
							if managers.DoUnblock(uid) == nil {
								msg += "\n- 已解封" + usr
							}
						}
					}
				} else {
					for _, usr := range args {
						uid, err := strconv.ParseInt(usr, 10, 64)
						if err == nil {
							if managers.DoBlock(uid) == nil {
								msg += "\n+ 已封禁" + usr
							}
						}
					}
				}
				ctx.SendChain(message.Text(msg))
				return
			}
			ctx.SendChain(message.Text("参数错误!"))
		})

		zero.OnCommandGroup([]string{
			"改变默认启用状态", "allflip",
		}, zero.SuperUserPermission).SetBlock(true).SecondPriority().Handle(func(ctx *zero.Ctx) {
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			service, ok := Lookup(model.Args)
			if !ok {
				ctx.SendChain(message.Text("没有找到指定服务!"))
				return
			}
			err := service.Flip()
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			ctx.SendChain(message.Text("已改变全局默认启用状态: " + model.Args))
		})

		zero.OnCommandGroup([]string{"用法", "usage"}, zero.UserOrGrpAdmin).SetBlock(true).SecondPriority().
			Handle(func(ctx *zero.Ctx) {
				model := extension.CommandModel{}
				_ = ctx.Parse(&model)
				service, ok := Lookup(model.Args)
				if !ok {
					ctx.SendChain(message.Text("没有找到指定服务!"))
					return
				}
				if service.Options.Help == "" {
					ctx.SendChain(message.Text("该服务无帮助!"))
					return
				}
				gid := ctx.Event.GroupID
				if gid == 0 {
					gid = -ctx.Event.UserID
				}
				serviceinfo := strings.Split(strings.Trim(service.String(), "\n"), "\n")
				var menu = mpic{
					kanban: "data/Control/kanban.png", // 看板娘图片
					stat:   "○未启用",                    // 启用状态
					stat2:  false,                     // 启用状态
					plugin: model.Args,                // 插件名
					ttf1:   text.BoldFontFile,         // 字体1
					ttf2:   text.SakuraFontFile,       // 字体2
					info:   serviceinfo,               // 插件信息
				}
				/***********获取看板娘图片***********/
				data, err := file.GetLazyData(menu.kanban, true)
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				w, _, err := gg.GetImgWH(data, menu.kanban) // 解析图片的宽高信息
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				var back image.Image

				if w > 1024 { // 限制图片大小
					back, _, err = image.Decode(bytes.NewReader(data))
					if err != nil {
						ctx.SendChain(message.Text("ERROR:", err))
						return
					}
					back = img.Limit(back, 1024, 1920)
					goto label //跳转
				}

				back, _, err = image.Decode(bytes.NewReader(data))
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
			label:
				menu.im = back                             // 设置看板娘图片
				_, err = file.GetLazyData(menu.ttf1, true) // 获取字体
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				_, err = file.GetLazyData(menu.ttf2, true)
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				if service.EnableMarkIn(gid) {
					menu.stat = "●已启用"
					menu.stat2 = true
				}
				pic, err := dyna(&menu)
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				data, cl := writer.ToBytes(pic) // 生成图片
				ctx.SendChain(message.ImageBytes(data))
				cl()
			})

		zero.OnCommandGroup([]string{"服务列表", "service_list"}, zero.UserOrGrpAdmin).SetBlock(true).SecondPriority().
			Handle(func(ctx *zero.Ctx) {
				i := 0
				j := 0
				gid := ctx.Event.GroupID
				if gid == 0 {
					gid = -ctx.Event.UserID
				}
				managers.RLock()
				msg := []string{"--------服务列表--------\n发送\"/用法 name\"查看详情\n发送\"/响应\"启用会话"}
				managers.RUnlock()
				var enableService []string
				var disableService []string
				managers.ForEach(func(key string, manager *ctrl.Control[*zero.Ctx]) bool {
					if manager.IsEnabledIn(gid) {
						i++
						enableService = append(enableService, strconv.Itoa(i)+":"+key)
					} else {
						j++
						disableService = append(disableService, strconv.Itoa(j)+":"+key)
					}
					return true
				})
				msg = append(msg, "\n\n→以下服务已开启：\n", strings.Join(enableService, "\n"))
				msg = append(msg, "\n\n→以下服务未开启：\n", strings.Join(disableService, "\n"))
				data, err := text.RenderToBase64(strings.Join(msg, ""), text.FontFile, 400, 20)
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				if id := ctx.SendChain(message.Image("base64://" + helper.BytesToString(data))); id.ID() == 0 {
					ctx.SendChain(message.Text("ERROR: 可能被风控了"))
				}
			})

		zero.OnCommandGroup([]string{"服务详情", "service_detail"}, zero.UserOrGrpAdmin).SetBlock(true).SecondPriority().
			Handle(func(ctx *zero.Ctx) {
				i := 0
				gid := ctx.Event.GroupID
				if gid == 0 {
					gid = -ctx.Event.UserID
				}
				managers.RLock()
				msgs := make([]any, 1, len(managers.M)*7+1)
				managers.RUnlock()
				msgs[0] = "---服务详情---\n"
				managers.ForEach(func(key string, service *ctrl.Control[*zero.Ctx]) bool {
					i++
					msgs = append(msgs, i, ": ", service.EnableMarkIn(gid), key, "\n", service, "\n\n")
					return true
				})
				data, err := text.RenderToBase64(fmt.Sprint(msgs...), text.FontFile, 400, 20)
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				if id := ctx.SendChain(message.Image("base64://" + helper.BytesToString(data))); id.ID() == 0 {
					ctx.SendChain(message.Text("ERROR: 可能被风控了"))
				}
			})
	})
}

// 返回菜单图片
func dyna(mp *mpic) (image.Image, error) {
	toww, fonts := 1024.0, 50.0                    // 图片宽度和字体大小
	tow := gg.NewContextForImage(mp.im)            // 看板娘
	one := gg.NewContext(int(toww*2.5), int(toww)) // 新图像
	one.SetRGB255(255, 255, 255)
	one.Clear()
	one.DrawImage(tow.Image(), 0, 0) // 放入看板娘的位置
	if err := one.LoadFontFace(mp.ttf2, fonts*2); err != nil {
		return nil, err
	}
	one.SetRGBA255(55, 55, 55, 255) // 字体颜色
	one.DrawString(mp.plugin, toww*1.35, fonts*2)
	if err := one.LoadFontFace(mp.ttf1, fonts); err != nil { // 加载字体
		return nil, err
	}
	one.DrawRoundedRectangle(toww+50, fonts-5, fonts*4.5, fonts*1.5, 10) // 创建圆角矩形
	one.SetRGBA255(221, 221, 221, 255)                                   // 填充背景颜色
	if mp.stat2 {
		one.SetRGBA255(0, 200, 0, 200)
	}
	one.Fill()                                // 填充颜色
	one.SetRGBA255(55, 55, 55, 255)           // 字体颜色
	one.DrawString(mp.stat, toww+58, fonts*2) // 绘制文字

	lineX, lineY, hi := 5.0, 140.0, 0   // 宽高记录
	maxTwidth := toww * 1.35            // 文字边距
	for i := 0; i < len(mp.info); i++ { // 遍历文字切片
		lineTexts := make([]string, 0, len(mp.info[i]))
		lineText, dynaw, dynah, tmpw := "", 0.0, 0.0, 0.0

		for len(mp.info[i]) > 0 {
			lineText, tmpw = truncate(one, mp.info[i], maxTwidth)
			lineTexts = append(lineTexts, lineText)
			if tmpw > dynaw {
				dynaw = tmpw
			}
			if len(lineText) >= len(mp.info[i]) {
				break // 如果写入的文本大于等于本次写入的文本 则跳出
			}
			dynah += fonts * 1.3                    // 截断一次加一行高度
			mp.info[i] = mp.info[i][len(lineText):] // 丢弃已经写入的文字并重新赋值
		}
		threew, threeh := dynaw+fonts, (dynah + (fonts * 1.2)) // 圆角矩形宽度和高度
		drawx := toww + lineX + 35                             // 圆角矩形位置宽度
		if int(lineX+dynaw+fonts)-int(toww) >= 180 {           // 越界宽容度
			goto label
		}
		one.DrawRoundedRectangle(drawx, lineY, threew, threeh, 20.0) // 创建圆角矩形
		drawsc(one, fonts, drawx, lineY, lineTexts)
		lineX += float64(threew) + fonts/2 // 添加后加一次宽度
		hi = int(threeh)

		continue // 跳出循环
	label:
		lineY += float64(hi) + fonts/4              // 加一次高度
		lineX = 5                                   // 重置宽度位置
		if threeh+lineY+fonts >= float64(one.H()) { // 超出最大高度则进行加高
			itmp := gg.NewContext(one.W(), int(lineY+threeh*2)) // X2 高度
			itmp.SetRGB255(255, 255, 255)
			itmp.Clear()
			itmp.DrawImage(one.Image(), 0, 0)
			one = gg.NewContextForImage(itmp.Image())
			if err := one.LoadFontFace(mp.ttf1, fonts); err != nil { // 加载字体
				return nil, err
			}
		}
		drawx = toww + lineX + 35                                    // 圆角矩形位置宽度
		one.DrawRoundedRectangle(drawx, lineY, threew, threeh, 20.0) // 创建圆角矩形
		drawsc(one, fonts, drawx, lineY, lineTexts)
		lineX += threew + fonts/2 // 添加后加一次宽度
		hi = int(threeh)
	}
	return one.Image(), nil
}

// 填充颜色和绘制文字
func drawsc(one *gg.Context, fonts, lineX, lineY float64, lineTexts []string) {
	rand.Seed(time.Now().UnixNano())
	one.SetRGBA255(rand.Intn(250), rand.Intn(250), rand.Intn(250), 90)
	one.Fill() // 填充颜色
	one.SetRGBA255(55, 55, 55, 255)
	h := fonts + lineY
	for i := range lineTexts { // 逐行绘制文字
		one.DrawString(lineTexts[i], lineX+fonts/2, h)
		h += fonts + (fonts / 4)
	}
}

// 截断文字
func truncate(one *gg.Context, text string, maxW float64) (string, float64) {
	var tmp strings.Builder
	tmp.Grow(len(text))
	res, w := make([]rune, 0, len(text)), 0.0
	for _, r := range text {
		tmp.WriteRune(r)
		width, _ := one.MeasureString(tmp.String()) // 获取文字宽度
		if width > maxW {                           // 如果宽度大于文字边距
			break //跳出
		} else {
			w = width
			res = append(res, r) // 写入
		}
	}
	return string(res), w
}
