package control

import (
	"os"

	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	yaml "gopkg.in/yaml.v2"

	ctrl "github.com/FloatTech/zbpctrl"
	zero "github.com/wdvxdr1123/ZeroBot"
)

var (
	enmap = make(map[string]*Engine)
	pm    = loadPrioMap()
)

// loadPrioMap 从文件中获取优先级配置文件的优先级，若文件不存在或文件格式错误，则返回一个空的切片类型
func loadPrioMap() []string {
	var pm []string
	data, err := os.ReadFile(priofile)
	if err != nil {
		logrus.Warnln("[control] 读取优先级配置文件失败,将使用默认的优先级配置:", err)
	} else {
		err := yaml.Unmarshal(data, &pm)
		if err != nil {
			logrus.Warnln("[control] 序列化优先级配置文件失败:", err)
		}
	}
	return pm
}

// writePrioMap 写入优先级配置文件的相关代码
func writePrioMap() {
	data, err := yaml.Marshal(&pm)

	if err != nil {
		logrus.Warnln("[control] 反序列化优先级配置文件失败:", err)
	}
	err0 := os.WriteFile(priofile, data, 0644)
	if err0 != nil {
		logrus.Warnln("[control] 写入优先级配置文件失败:", err0)
	}
}

// getPrioFromProfile 获取优先级配置文件中的优先级，若不存在，则在末尾添加并写回文件
func getPrioFromProfile(s string) int {
	i := slices.Index(pm, s)
	if i == -1 {
		pm = append(pm, s)
		i = slices.Index(pm, s)
		writePrioMap()
	}
	return (i + 1) * 10
}

// Register 注册插件控制器
func Register(service string, o *ctrl.Options[*zero.Ctx]) *Engine {
	engine := newengine(service, getPrioFromProfile(service), o)
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
