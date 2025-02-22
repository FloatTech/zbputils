package chat

import (
	_ "embed"

	"github.com/fumiama/deepinfra"
	"github.com/fumiama/deepinfra/chat"
	"github.com/fumiama/deepinfra/model"

	zero "github.com/wdvxdr1123/ZeroBot"
)

//go:embed README.md
var SystemPrompt string

var lst = chat.NewLog(8, "\n\n", "自己随机开启新话题", "【", "】", ">>")

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
		lst.Add(
			gid, ctx.Event.Sender.Name(),
			ctx.State["__zbputil_chat_txt__"].(string),
			false, ctx.Event.IsToMe,
		)
	})
}

func Reply(grp int64, txt string) {
	lst.Add(grp, "", txt, true, false)
}

func Ask(p model.Protocol, grp int64, sysp string) deepinfra.Model {
	return lst.Modelize(p, grp, sysp)
}
