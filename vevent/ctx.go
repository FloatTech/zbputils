package vevent

import (
	"sync"
	"unsafe"

	zero "github.com/wdvxdr1123/ZeroBot"
)

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

func HookCtxCaller(ctx *zero.Ctx, hook zero.APICaller) {
	(*(**Ctx)(unsafe.Pointer(&ctx))).caller = hook
}
