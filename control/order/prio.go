// Package order 各个插件的优先级
package order

import "sync/atomic"

var prio uint64

func AcquirePrio() int {
	return int(atomic.AddUint64(&prio, 10))
}
