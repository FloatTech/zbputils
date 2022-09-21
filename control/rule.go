// Package control 控制插件的启用与优先级等
package control

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unsafe"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension"
	"github.com/wdvxdr1123/ZeroBot/message"

	ctrl "github.com/FloatTech/zbpctrl"

	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/process"

	"github.com/FloatTech/floatbox/img/writer"
	"github.com/FloatTech/zbputils/img/text"
)

var (
	// managers 每个插件对应的管理
	managers = ctrl.NewManager[*zero.Ctx]("data/control/plugins.db")
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
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			ctx.SendChain(message.Text("已改变全局默认启用状态: " + model.Args))
		})

		zero.OnRegex(`^(开启|关闭|重置|on|off|reset)看板娘$`, zero.SuperUserPermission).
			SetBlock(true).Handle(func(ctx *zero.Ctx) {
			str := ctx.State["regex_matched"].([]string)[1]
			switch {
			case str == "开启", str == "on":
				kanbanEnable = true
			case str == "关闭", str == "off":
				kanbanEnable = false
			case str == "重置", str == "reset":
				customKanban = false
				roleName = "kanban.png"
			default:
				ctx.SendChain(message.Text("输入错误了."))
				return
			}
			ctx.SendChain(message.Text("已经 ", str, " 了!"))
		})

		zero.OnRegex(`^设置看板娘头像(\[CQ:at,qq=(.\d+)\]|(.\d+))`, zero.SuperUserPermission).SetBlock(true).
			Handle(func(ctx *zero.Ctx) {
				str := ctx.State["regex_matched"].([]string)
				user := ""
				id, err := strconv.ParseInt(str[2], 10, 64)
				if err != nil {
					id, err = strconv.ParseInt(str[1], 10, 64)
					if err != nil {
						ctx.SendChain(message.Text("ERROR: ", err))
						return
					}
				}
				user = strconv.FormatInt(id, 10)
				if id > 100 {
					err = file.DownloadTo("http://q4.qlogo.cn/g?b=qq&nk="+user+"&s=640",
						kanbanPath+"img/"+user+".jpg", false)
					if err != nil {
						ctx.SendChain(message.Text("ERROR: ", err))
						return
					}
					customKanban = true
					roleName = "img/" + user + ".jpg"
					ctx.SendChain(message.Text("设置成功。"))
				}
			})

		zero.OnKeyword("看板娘图片", zero.MustProvidePicture, zero.SuperUserPermission).
			SetBlock(true).Handle(func(ctx *zero.Ctx) {
			id := fmt.Sprint(ctx.Event.UserID)
			url := ctx.State["image_url"].([]string)
			err := file.DownloadTo(url[0], kanbanPath+"img/"+id+".jpg", false)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			customKanban = true
			roleName = "img/" + id + ".jpg"
			ctx.SendChain(message.Text("设置成功。"))
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
				_, err := file.GetLazyData(text.BoldFontFile, true)
				if err != nil {
					ctx.SendChain(message.Text("ERROR: ", err))
					return
				}
				_, err = file.GetLazyData(text.SakuraFontFile, true)
				if err != nil {
					ctx.SendChain(message.Text("ERROR: ", err))
					return
				}
				gid := ctx.Event.GroupID
				if gid == 0 {
					gid = -ctx.Event.UserID
				}
				/***********获取看板娘图片***********/
				serviceinfo := strings.Split(strings.Trim(service.String(), "\n"), "\n")
				var menu = mpic{
					path:       kanbanPath + roleName, // 看板娘图片
					isDisplay:  kanbanEnable,          // 显示看板娘
					isCustom:   customKanban,
					statusText: "●已启用",              // 启用状态
					enableText: true,                // 启用状态
					pluginName: model.Args,          // 插件名
					font1:      text.BoldFontFile,   // 字体1
					font2:      text.SakuraFontFile, // 字体2
					info:       serviceinfo,         // 插件信息
					multiple:   2.5,                 // 倍数
					fontSize:   50,                  // 字体大小
				}
				if !service.EnableMarkIn(gid) {
					menu.statusText = "○未启用"
					menu.enableText = false
				}
				err = menu.loadpic()
				if err != nil {
					ctx.SendChain(message.Text("ERROR: ", err))
					return
				}
				pic, err := menu.dyna(&location{
					lastH:     0, //
					drawX:     0.0,
					maxTwidth: 1200.0, // 文字边距
					rlineX:    0.0,    // 宽高记录
					rlineY:    140.0,
				})
				if err != nil {
					ctx.SendChain(message.Text("ERROR: ", err))
					return
				}
				data, cl := writer.ToBytes(pic) // 生成图片
				if id := ctx.SendChain(message.ImageBytes(data)); id.ID() == 0 {
					ctx.SendChain(message.Text("ERROR: 可能被风控了"))
				}
				cl()
			})

		zero.OnCommandGroup([]string{"服务列表", "service_list"}, zero.UserOrGrpAdmin).SetBlock(true).SecondPriority().
			Handle(func(ctx *zero.Ctx) {
				_, err := file.GetLazyData(text.BoldFontFile, true)
				if err != nil {
					ctx.SendChain(message.Text("ERROR: ", err))
					return
				}
				_, err = file.GetLazyData(text.SakuraFontFile, true)
				if err != nil {
					ctx.SendChain(message.Text("ERROR: ", err))
					return
				}
				i := 0
				j := 0
				gid := ctx.Event.GroupID
				if gid == 0 {
					gid = -ctx.Event.UserID
				}
				managers.RLock()
				var tmp strings.Builder
				var enable strings.Builder
				var disable strings.Builder
				tmp.Grow(512)
				enable.Grow(256)
				tmp.WriteString("\t\t\t <----------服务列表---------> \t\t\t\n" +
					"\t\t\t\t   ◇发送\"/用法 name\"查看详情   \t\t\t\t\n" +
					"\t\t\t\t   ◇发送\"/响应\"启用会话        \t\t\t\t\t\n",
				)
				managers.RUnlock()
				enable.WriteString("\t\t\t\t\t     ↓ 以下服务已开启↓         \t\t\t\t\n")
				disable.WriteString("\t\t\t\t\t     ↓ 以下服务未开启↓         \t\t\t\t\n")

				managers.ForEach(func(key string, manager *ctrl.Control[*zero.Ctx]) bool {
					if manager.IsEnabledIn(gid) {
						i++
						enable.WriteString(strconv.Itoa(i) + ": " + key + "\n")
					} else {
						j++
						disable.WriteString(strconv.Itoa(j) + ": " + key + "\n")
					}
					return true
				})
				tmp.WriteString(enable.String())
				tmp.WriteString(disable.String())
				msg := strings.Split(strings.Trim(tmp.String(), "\n"), "\n")

				var menu = mpic{
					path:       kanbanPath + roleName, // 看板娘图片
					isDisplay:  kanbanEnable,          // 显示看板娘
					isCustom:   customKanban,
					statusText: "●Plugin",           // 启用状态
					enableText: true,                // 启用状态
					pluginName: "ZeroBot-Plugin",    // 插件名
					font1:      text.BoldFontFile,   // 字体1
					font2:      text.SakuraFontFile, // 字体2
					info:       msg,                 // 插件信息
					multiple:   2.5,                 // 倍数
					fontSize:   50,                  // 字体大小
				}
				err = menu.loadpic()
				if err != nil {
					ctx.SendChain(message.Text("ERROR: ", err))
					return
				}
				pic, err := menu.dyna(&location{
					lastH:     0,
					drawX:     0.0,
					maxTwidth: 1200.0, // 文字边距
					rlineX:    0.0,    // 宽高记录
					rlineY:    140.0,
				})
				if err != nil {
					ctx.SendChain(message.Text("ERROR: ", err))
					return
				}
				data, cl := writer.ToBytes(pic) // 生成图片
				if id := ctx.SendChain(message.ImageBytes(data)); id.ID() == 0 {
					ctx.SendChain(message.Text("ERROR: 可能被风控了"))
				}
				cl()
			})

		zero.OnCommandGroup([]string{"服务详情", "service_detail"}, zero.UserOrGrpAdmin).SetBlock(true).SecondPriority().
			Handle(func(ctx *zero.Ctx) {
				_, err := file.GetLazyData(text.BoldFontFile, true)
				if err != nil {
					ctx.SendChain(message.Text("ERROR: ", err))
					return
				}
				_, err = file.GetLazyData(text.SakuraFontFile, true)
				if err != nil {
					ctx.SendChain(message.Text("ERROR: ", err))
					return
				}
				i, j := 0, 0
				double, fontSize, multiple := true, 40.0, 5.0
				gid := ctx.Event.GroupID
				if gid == 0 {
					gid = -ctx.Event.UserID
				}
				managers.RLock()
				lenmap := len(managers.M)
				managers.RUnlock()
				if lenmap == 0 {
					ctx.SendChain(message.Text("服务数量为: ", lenmap))
					return
				}
				var tmp strings.Builder
				var tmp2 strings.Builder
				tmp.Grow(lenmap * 100)
				tmp2.Grow(lenmap * 100)

				end := lenmap / 2
				if lenmap <= 5 {
					double = false // 单列模式
					fontSize = 40
					multiple = 3
				}
				tab := "\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\n"
				tmp.WriteString("\t\t\t\t  <---------服务详情--------->   \t\t\t\t")
				managers.ForEach(func(key string, service *ctrl.Control[*zero.Ctx]) bool {
					i++
					if i > end+1 && lenmap > 5 {
						goto label
					}
					tmp.WriteString(fmt.Sprint("\n", i, ": ", service.EnableMarkIn(gid), " ", key,
						tab, strings.Trim(fmt.Sprint(service), "\n")))
					return true

				label:
					if j > 0 {
						tmp2.WriteString(fmt.Sprint("\n", i, ": ", service.EnableMarkIn(gid), " ", key,
							tab, strings.Trim(fmt.Sprint(service), "\n")))
					} else {
						tmp2.WriteString(fmt.Sprint(i, ": ", service.EnableMarkIn(gid), " ", key,
							tab, strings.Trim(fmt.Sprint(service), "\n")))
					}
					j++
					return true
				})

				msg := strings.Split(tmp.String(), "\n")
				msg2 := strings.Split(tmp2.String(), "\n")
				var menu = mpic{
					path:       kanbanPath + roleName, // 看板娘图片
					isDisplay:  kanbanEnable,          // 显示看板娘
					isCustom:   customKanban,
					statusText: "○ INFO",            // 启用状态
					enableText: false,               // 启用状态
					isDouble:   double,              // 双列排版
					pluginName: "ZeroBot-Plugin",    // 插件名
					font1:      text.BoldFontFile,   // 字体1
					font2:      text.SakuraFontFile, // 字体2
					info:       msg,                 // 插件信息
					info2:      msg2,                // 插件信息
					multiple:   multiple,            // 倍数
					fontSize:   fontSize,            // 字体大小
				}
				err = menu.loadpic()
				if err != nil {
					ctx.SendChain(message.Text("ERROR: ", err))
					return
				}
				pic, err := menu.dyna(&location{
					lastH:     0,
					drawX:     0.0,
					maxTwidth: 1200.0, // 文字边距
					rlineX:    0.0,    // 宽高记录
					rlineY:    140.0,
				})
				if err != nil {
					ctx.SendChain(message.Text("ERROR: ", err))
					return
				}
				data, cl := writer.ToBytes(pic) // 生成图片
				if id := ctx.SendChain(message.ImageBytes(data)); id.ID() == 0 {
					ctx.SendChain(message.Text("ERROR: 可能被风控了"))
				}
				cl()
			})
	})
}
