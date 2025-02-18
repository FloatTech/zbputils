package chat

import (
	"github.com/fumiama/deepinfra"

	zero "github.com/wdvxdr1123/ZeroBot"
)

var lst = newlist()

func init() {
	zero.OnMessage(func(ctx *zero.Ctx) bool {
		txt := ctx.ExtractPlainText()
		ctx.State["__zbputil_chat_txt__"] = txt
		return txt != ""
	}).FirstPriority().SetBlock(false).Handle(func(ctx *zero.Ctx) {
		gid := ctx.Event.GroupID
		if gid == 0 {
			gid = -ctx.Event.UserID
		}
		lst.add(gid, ctx.Event.Sender.Name(), ctx.State["__zbputil_chat_txt__"].(string), false)
	})
}

func Reply(ctx *zero.Ctx, txt string) {
	gid := ctx.Event.GroupID
	if gid == 0 {
		gid = -ctx.Event.UserID
	}
	lst.add(gid, "", txt, true)
}

func Ask(ctx *zero.Ctx, temp float32, mn, sysp, sepstr string) deepinfra.Model {
	gid := ctx.Event.GroupID
	if gid == 0 {
		gid = -ctx.Event.UserID
	}
	return lst.modelize(temp, gid, mn, sysp, sepstr)
}
