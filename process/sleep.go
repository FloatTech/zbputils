// Package process 流程控制相关
package process

import (
	"math/rand"
	"sync"
	"time"
)

// GlobalInitMutex 在 init 时被冻结, main 初始化完成后解冻
var GlobalInitMutex = func() (mu sync.Mutex) {
	mu.Lock()
	return
}()

// SleepAbout1sTo2s 随机阻塞等待 1 ~ 2s
func SleepAbout1sTo2s() {
	time.Sleep(time.Second + time.Millisecond*time.Duration(rand.Intn(1000)))
}
