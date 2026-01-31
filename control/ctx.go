package control

import (
	"github.com/RomiChan/syncx"
	zero "github.com/wdvxdr1123/ZeroBot"
)

const (
	StateKeySyncxState = zero.StateKeyPrefixKeep + "_ctrl_syncx_state__"
)

func addsyncxstate(ctx *zero.Ctx) bool {
	ctx.State[StateKeySyncxState] = &syncx.Map[string, any]{}
	return true
}
