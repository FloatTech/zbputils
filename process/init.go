package process

import (
	"sync"

	zero "github.com/wdvxdr1123/ZeroBot"
)

func DoOnceOnSuccess(f func() error, onerr func(error)) zero.Rule {
	var init sync.Once
	return func(ctx *zero.Ctx) bool {
		var err error
		init.Do(func() {
			err = f()
		})
		if err != nil {
			onerr(err)
			init = sync.Once{}
			return false
		}
		return true
	}
}
