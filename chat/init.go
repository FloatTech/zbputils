// Package chat 提供聊天记录管理和 AI 模型交互功能
package chat

import (
	_ "embed"
	"strings"

	"github.com/fumiama/deepinfra"
	"github.com/fumiama/deepinfra/chat"
	"github.com/fumiama/deepinfra/model"

	zero "github.com/wdvxdr1123/ZeroBot"
)

// SystemPrompt 将 README.md 内容嵌入为默认系统提示词
//
//go:embed README.md
var SystemPrompt string

// lst 全局聊天记录，每群/每用户独立保存最近 8 条
var lst = chat.NewLog(8, "\n\n", "自己随机开启新话题", "【", "】", ">>")

func init() {
	// 注册 ZeroBot 消息钩子，记录所有非空文本
	zero.OnMessage(func(ctx *zero.Ctx) bool {
		txt := ctx.ExtractPlainText()
		ctx.State["__zbputil_chat_txt__"] = txt
		return txt != ""
	}).FirstPriority().SetBlock(false).Handle(func(ctx *zero.Ctx) {
		// 计算群组 ID（私聊时使用负的 UserID）
		gid := ctx.Event.GroupID
		if gid == 0 {
			gid = -ctx.Event.UserID
		}
		// 将用户消息追加到对应群组的聊天记录
		lst.Add(
			gid, ctx.Event.Sender.Name(),
			ctx.State["__zbputil_chat_txt__"].(string),
			false, ctx.Event.IsToMe,
		)
	})
}

// Reply 将 AI 回复追加到指定群组的聊天记录
func Reply(grp int64, txt string) {
	lst.Add(grp, "", txt, true, false)
}

// Ask 根据聊天记录构造可执行的 deepinfra 模型请求
func Ask(p model.Protocol, grp int64, sysp string, isusersys bool) deepinfra.Model {
	return lst.Modelize(p, grp, sysp, isusersys)
}

// AskCustom 通用日志转换函数，允许自定义每条记录的映射逻辑
func AskCustom[T any](grp int64, f func(int, string) T) []T {
	return chat.Modelize(&lst, grp, f)
}

// Sanitize 清洗 AI 返回文本：
// 1. 去掉换行后内容
// 2. 去掉发言前缀（如【name】或[name]）
// 3. 去掉重复 10 次以上的子串
// 4. 去除首尾空白
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

// Reset 清空全局聊天记录，重新开始
func Reset() {
	lst.Reset()
}
