package vevent

import (
	_ "unsafe"

	zero "github.com/wdvxdr1123/ZeroBot"
)

// processEvent 处理事件
//go:linkname processEvent github.com/wdvxdr1123/ZeroBot.processEvent
func processEvent(response []byte, caller zero.APICaller)
