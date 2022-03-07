package vevent

import (
	"unsafe"

	zero "github.com/wdvxdr1123/ZeroBot"
)

type Loop struct {
	caller zero.APICaller
}

func NewLoop(bot int64) (lo Loop) {
	zero.RangeBot(func(id int64, ctx *zero.Ctx) bool {
		if id == bot {
			lo.caller = (*(**Ctx)(unsafe.Pointer(&ctx))).caller
			return false
		}
		return true
	})
	return
}

func NewLoopOf(caller zero.APICaller) Loop {
	return Loop{caller: caller}
}

func (lo Loop) Echo(response []byte) {
	processEvent(response, lo.caller)
}
