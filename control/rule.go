// Package control 控制插件的启用与优先级等
package control

import (
	"os"
	"strconv"
	"strings"
	"unsafe"

	"github.com/Coloured-glaze/gg"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension"
	"github.com/wdvxdr1123/ZeroBot/message"

	ctrl "github.com/FloatTech/zbpctrl"

	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/process"

	"github.com/FloatTech/rendercard"

	"github.com/FloatTech/zbputils/ctxext"
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
				// 处理插件帮助并且计算图像高
				plugininfo := strings.Split(strings.Trim(service.String(), "\n"), "\n")
				newplugininfo := make([]string, 0, len(plugininfo)*2)
				font := gg.NewContext(1, 1)
				err = font.LoadFontFace(text.BoldFontFile, 38)
				if err != nil {
					return
				}
				for i := 0; i < len(plugininfo); i++ {
					newlinetext, textw, tmpw := "", 0.0, 0.0
					for len(plugininfo[i]) > 0 {
						newlinetext, tmpw = truncate(font, plugininfo[i], 1272-50)
						newplugininfo = append(newplugininfo, newlinetext)
						if tmpw > textw {
							textw = tmpw
						}
						if len(newlinetext) >= len(plugininfo[i]) {
							break
						}
						plugininfo[i] = plugininfo[i][len(newlinetext):]
					}
				}
				var imgs []byte
				imgs, err = rendercard.Titleinfo{
					Lefttitle:     service.Service,
					Leftsubtitle:  service.Options.Brief,
					Righttitle:    "FloatTech",
					Rightsubtitle: "ZeroBot-Plugin",
					Imgpath:       kanbanpath + "icon.jpg",
					Textpath:      text.SakuraFontFile,
					Textpath2:     text.BoldFontFile,
					Status:        service.IsEnabledIn(gid),
				}.Drawtitledtext(newplugininfo)
				if err != nil {
					ctx.SendChain(message.Text("ERROR: ", err))
					return
				}
				if id := ctx.SendChain(message.ImageBytes(imgs)); id.ID() == 0 {
					ctx.SendChain(message.Text("ERROR: 可能被风控了"))
				}
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
				gid := ctx.Event.GroupID
				if gid == 0 {
					gid = -ctx.Event.UserID
				}
				var imgs [][]byte
				imgs, err = drawservicesof(gid)
				if err != nil {
					ctx.SendChain(message.Text("ERROR: ", err))
					return
				}
				msg := make(message.Message, len(imgs))
				for i := 0; i < len(imgs); i++ {
					msg[i] = ctxext.FakeSenderForwardNode(ctx, message.ImageBytes(imgs[i]))
				}
				if id := ctx.Send(msg); id.ID() == 0 {
					ctx.SendChain(message.Text("ERROR: 可能被风控了"))
				}
			})
	})
}
