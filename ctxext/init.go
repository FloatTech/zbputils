// Package ctxext zb context 扩展
package ctxext

import (
	"github.com/FloatTech/zbputils/process"
)

// DoOnceOnSuccess 当返回 true, 之后直接通过, 否则下次触发仍会执行
func DoOnceOnSuccess[Ctx any](f func(Ctx) bool) func(Ctx) bool {
	init := process.NewOnce()
	return func(ctx Ctx) (success bool) {
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
