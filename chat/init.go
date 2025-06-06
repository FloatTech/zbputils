package chat

import (
	_ "embed"
	"strings"

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

func Ask(p model.Protocol, grp int64, sysp string, isusersys bool) deepinfra.Model {
	return lst.Modelize(p, grp, sysp, isusersys)
}

func AskCustom[T any](grp int64, f func(int, string) T) []T {
	return chat.Modelize(&lst, grp, f)
}

func Sanitize(msg string) string {
	msg, _, _ = strings.Cut(msg, "\n")
	msg = strings.TrimSpace(msg)
	i := strings.LastIndex(msg, "】")
	if i > 0 {
		if i+len("】") >= len(msg) {
			return ""
		}
		msg = msg[i+len("】"):]
	} else {
		i = strings.LastIndex(msg, "]")
		if i > 0 {
			if i+1 >= len(msg) {
				return ""
			}
			msg = msg[i:]
		}
	}
	if s, n := findRepeatedPattern(msg, 10); n > 0 {
		return s
	}
	return strings.TrimSpace(msg)
}
