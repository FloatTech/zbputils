package vevent

import (
	"sync"
	"unsafe"

	zero "github.com/wdvxdr1123/ZeroBot"
)

type Loop struct{ bot int64 }

func NewLoop(bot int64) Loop {
	return Loop{bot: bot}
}

// Ctx represents the Context which hold the event.
// 代表上下文
type Ctx struct {
	ma     *zero.Matcher
	Event  *zero.Event
	State  zero.State
	caller zero.APICaller

	// lazy message
	once    sync.Once
	message string
}

func (lo Loop) Echo(response []byte) {
	zero.RangeBot(func(id int64, ctx *zero.Ctx) bool {
		if id == lo.bot {
			processEvent(response, (*(**Ctx)(unsafe.Pointer(&ctx))).caller)
			return false
		}
		return true
	})
}
