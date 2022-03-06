package vevent

import (
	zero "github.com/wdvxdr1123/ZeroBot"
)

type Loop struct {
	drv zero.APICaller
}

func NewLoop(drv zero.APICaller) Loop {
	return Loop{drv: drv}
}

func (lo Loop) Echo(response []byte) {
	processEvent(response, lo.drv)
}
