// Package control 控制插件的启用与优先级等
package control

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unsafe"
	"bytes"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"

	ctrl "github.com/FloatTech/zbpctrl"

	"github.com/FloatTech/zbputils/process"


	// 图片输出
	"image"
	"github.com/FloatTech/zbputils/file"
	"github.com/FloatTech/zbputils/img"
	"github.com/FloatTech/zbputils/img/text"
	"github.com/FloatTech/zbputils/img/writer"
	"github.com/fogleman/gg"
)

var (
	// managers 每个插件对应的管理
	managers = ctrl.NewManager[*zero.Ctx]("data/control/plugins.db", 10*time.Second)
)

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
				// 绘制图片
				/***********获取看板娘图片***********/
				data, err := file.GetLazyData("data/Control/kanban.png", true)
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				back, _, err := image.Decode(bytes.NewReader(data))
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				/***********设置图片的大小和底色***********/
				backX := 902
				backY := 1056
				canvas := gg.NewContext(backX, backY)
				// 设置文字大小
				fontSize := 50.0
				_, err = file.GetLazyData(text.BoldFontFile, true)
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				_, err = file.GetLazyData(text.SakuraFontFile, true)
				if err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				if err = canvas.LoadFontFace(text.BoldFontFile, fontSize); err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				// 计算卡片大小
				backXmax := 1500
				backYmax := 1056
				serviceinfo := strings.Split(service.String(), "\n")
				for i, info := range serviceinfo {
					width, h := canvas.MeasureString(info)
					if backXmax < int(width) {
						backXmax = int(width) + 100 // 获取最大宽度
					}
					high := 300 + i*int(h+20) // 获取文本高度
					if backYmax < high {
						backYmax = high // 获取最大高度
					}
				}
				canvas = gg.NewContext(backXmax+backX+50, backYmax)
				// 设置背景色
				canvas.SetRGB(1, 1, 1)
				canvas.Clear()
				/***********放置好看的图片***********/
				back = img.Size(back, backX, backY).Im
				canvas.DrawImage(back, 0, 0)
				/***********写入插件信息***********/
				if err = canvas.LoadFontFace(text.BoldFontFile, fontSize); err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				xCoordinate := float64(backX + 50)       // 看板娘的坐标
				length, h := canvas.MeasureString("看用法") // 获取文字宽度与高度
				// 标记启动状态
				canvas.DrawRoundedRectangle(xCoordinate-length*0.1, 130-h*2.5, length*1.2, h*2, fontSize*0.2)
				enable := "未启用"
				if service.EnableMarkIn(gid) {
					canvas.SetRGB255(0, 221, 0)
					enable = "已启用"
				} else {
					canvas.SetRGB255(221, 221, 221)
				}
				canvas.Fill()
				canvas.SetRGB(0, 0, 0)
				canvas.DrawString(enable, xCoordinate, 130-h)
				// 写入插件helper内容
				canvas.DrawRoundedRectangle(xCoordinate-30, 140, float64(backXmax)-30, float64(backYmax-160), fontSize*0.3)
				canvas.SetRGB255(0, 221, 221)
				canvas.Fill()
				canvas.SetRGB(0, 0, 0)
				for i, info := range serviceinfo {
					canvas.DrawString(info, xCoordinate, 260+(h+20)*float64(i-1))
				}
				// 写入插件名称
				if err = canvas.LoadFontFace(text.SakuraFontFile, fontSize*2); err != nil {
					ctx.SendChain(message.Text("ERROR:", err))
					return
				}
				canvas.DrawString(model.Args, xCoordinate+length*1.2, 130-h)
				// 生成图片
				data, cl := writer.ToBytes(canvas.Image())
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
