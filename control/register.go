package control

import (
	"sync/atomic"
)

var (
	enmap = make(map[string]*engineinstance)
	prio  uint64
)

// Register 注册插件控制器
func Register(service string, o *Options) Engine {
	engine := newengine(service, int(atomic.AddUint64(&prio, 10)), o)
	enmap[service] = engine
	return engine
}

// Delete 删除插件控制器, 不会删除数据
func Delete(service string) {
	engine, ok := enmap[service]
	if ok {
		engine.Delete()
		manmu.RLock()
		_, ok = managers[service]
		manmu.RUnlock()
		if ok {
			manmu.Lock()
			delete(managers, service)
			manmu.Unlock()
		}
	}
}
