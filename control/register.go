package control

import (
	"sync/atomic"

	ctrl "github.com/FloatTech/zbpctrl"
	zero "github.com/wdvxdr1123/ZeroBot"
)

var (
	enmap = make(map[string]*engineinstance)
	prio  uint64
)

// Register 注册插件控制器
func Register(service string, o *ctrl.Options[*zero.Ctx]) Engine {
	engine := newengine(service, int(atomic.AddUint64(&prio, 10)), o)
	enmap[service] = engine
	return engine
}

// Delete 删除插件控制器, 不会删除数据
func Delete(service string) {
	engine, ok := enmap[service]
	if ok {
		engine.Delete()
		managers.RLock()
		_, ok = managers.M[service]
		managers.RUnlock()
		if ok {
			managers.Lock()
			delete(managers.M, service)
			managers.Unlock()
		}
	}
}
