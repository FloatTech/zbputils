package ctxext

import (
	zero "github.com/wdvxdr1123/ZeroBot"

	"github.com/FloatTech/zbputils/process"
)

// DoOnceOnSuccess 当返回 true, 之后直接通过, 否则下次触发仍会执行
func DoOnceOnSuccess(f zero.Rule) zero.Rule {
	init := process.NewOnce()
	return func(ctx *zero.Ctx) (success bool) {
		success = true
		init.Do(func() {
			success = f(ctx)
		})
		if !success {
			init.Reset()
		}
		return
	}
}
