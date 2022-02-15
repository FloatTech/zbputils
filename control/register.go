package control

var enmap = make(map[string]*engineinstance)

// Register 注册插件控制器
func Register(service string, prio int, o *Options) Engine {
	engine := newengine(service, prio, o)
	enmap[service] = engine
	return engine
}

// Delete 删除插件控制器，不会删除数据
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
