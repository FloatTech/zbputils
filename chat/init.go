package chat

import (
	"github.com/fumiama/deepinfra"
	"github.com/fumiama/deepinfra/model"

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

func Reply(grp int64, txt string) {
	lst.add(grp, "", txt, true)
}

func Ask(p model.Protocol, grp int64, sysp string) deepinfra.Model {
	return lst.modelize(p, grp, sysp)
}
