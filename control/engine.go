package control

import (
	"fmt"
	"os"
	"unicode"

	"github.com/FloatTech/zbputils/file"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/single"
)

type Engine interface {
	// Delete 移除该 Engine 注册的所有 Matchers
	Delete()
	// On 添加新的指定消息类型的匹配器
	On(typ string, rules ...zero.Rule) Matcher
	// OnMessage 消息触发器
	OnMessage(rules ...zero.Rule) Matcher
	// OnNotice 系统提示触发器
	OnNotice(rules ...zero.Rule) Matcher
	// OnRequest 请求消息触发器
	OnRequest(rules ...zero.Rule) Matcher
	// OnMetaEvent 元事件触发器
	OnMetaEvent(rules ...zero.Rule) Matcher
	// OnPrefix 前缀触发器
	OnPrefix(prefix string, rules ...zero.Rule) Matcher
	// OnSuffix 后缀触发器
	OnSuffix(suffix string, rules ...zero.Rule) Matcher
	// OnCommand 命令触发器
	OnCommand(commands string, rules ...zero.Rule) Matcher
	// OnRegex 正则触发器
	OnRegex(regexPattern string, rules ...zero.Rule) Matcher
	// OnKeyword 关键词触发器
	OnKeyword(keyword string, rules ...zero.Rule) Matcher
	// OnFullMatch 完全匹配触发器
	OnFullMatch(src string, rules ...zero.Rule) Matcher
	// OnFullMatchGroup 完全匹配触发器组
	OnFullMatchGroup(src []string, rules ...zero.Rule) Matcher
	// OnKeywordGroup 关键词触发器组
	OnKeywordGroup(keywords []string, rules ...zero.Rule) Matcher
	// OnCommandGroup 命令触发器组
	OnCommandGroup(commands []string, rules ...zero.Rule) Matcher
	// OnPrefixGroup 前缀触发器组
	OnPrefixGroup(prefix []string, rules ...zero.Rule) Matcher
	// OnSuffixGroup 后缀触发器组
	OnSuffixGroup(suffix []string, rules ...zero.Rule) Matcher
	// OnShell shell命令触发器
	OnShell(command string, model any, rules ...zero.Rule) Matcher
	// ApplySingle 应用反并发
	ApplySingle(*single.Single[int64]) Engine
	// DataFolder 本插件数据目录, 默认 data/zbp/
	DataFolder() string
	// IsEnabledIn 自己是否在 id (正群负个人零全局) 启用
	IsEnabledIn(id int64) bool
	// 下载并获取本 engine 文件夹下的懒加载数据
	GetLazyData(filename string, isDataMustEqual bool) ([]byte, error)
}

type engineinstance struct {
	en         *zero.Engine
	prio       int
	service    string
	datafolder string
}

var priomap = make(map[int]string)      // priomap is map[prio]service
var foldermap = make(map[string]string) // foldermap is map[folder]service

func newengine(service string, prio int, o *Options) (e *engineinstance) {
	e = new(engineinstance)
	s, ok := priomap[prio]
	if ok {
		panic(fmt.Sprint("prio", prio, "is used by", s))
	}
	priomap[prio] = service
	logrus.Debugln("[control]插件", service, "已设置优先级", prio)
	e.en = zero.New()
	e.en.UsePreHandler(
		func(ctx *zero.Ctx) bool {
			// 防止自触发
			return ctx.Event.UserID != ctx.Event.SelfID
		},
		newctrl(service, o).Handler,
	)
	e.prio = prio
	e.service = service
	if o.PublicDataFolder != "" {
		if unicode.IsLower([]rune(o.PublicDataFolder)[0]) {
			panic("public data folder " + o.PublicDataFolder + " must start with an upper case letter")
		}
		e.datafolder = "data/" + o.PublicDataFolder + "/"
	} else if o.PrivateDataFolder != "" {
		if unicode.IsUpper([]rune(o.PrivateDataFolder)[0]) {
			panic("private data folder " + o.PrivateDataFolder + " must start with an lower case letter")
		}
		e.datafolder = "data/" + o.PrivateDataFolder + "/"
	} else {
		e.datafolder = "data/zbp/"
	}
	if e.datafolder != "data/zbp/" {
		s, ok := foldermap[e.datafolder]
		if ok {
			panic("folder " + e.datafolder + " has been required by service " + s)
		}
		foldermap[e.datafolder] = service
	}
	if file.IsNotExist(e.datafolder) {
		err := os.MkdirAll(e.datafolder, 0755)
		if err != nil {
			panic(err)
		}
	}
	logrus.Debugln("[control]插件", service, "已设置数据目录", e.datafolder)
	return
}

// DataFolder 本插件数据目录, 默认 data/zbp/
func (e *engineinstance) DataFolder() string {
	return e.datafolder
}

// IsEnabledIn 自己是否在 id (正群负个人零全局) 启用
func (e *engineinstance) IsEnabledIn(id int64) bool {
	c, ok := Lookup(e.service)
	if !ok {
		return false
	}
	return c.IsEnabledIn(id)
}

// Delete 移除该 Engine 注册的所有 Matchers
func (e *engineinstance) Delete() {
	e.en.Delete()
}

// On 添加新的指定消息类型的匹配器
func (e *engineinstance) On(typ string, rules ...zero.Rule) Matcher {
	return (*matcherinstance)(e.en.On(typ, rules...).SetPriority(e.prio))
}

// On 添加新的指定消息类型的匹配器(默认Engine)
func On(typ string, rules ...zero.Rule) Matcher {
	return (*matcherinstance)(zero.On(typ, rules...))
}

// OnMessage 消息触发器
func (e *engineinstance) OnMessage(rules ...zero.Rule) Matcher { return e.On("message", rules...) }

// OnNotice 系统提示触发器
func (e *engineinstance) OnNotice(rules ...zero.Rule) Matcher { return e.On("notice", rules...) }

// OnRequest 请求消息触发器
func (e *engineinstance) OnRequest(rules ...zero.Rule) Matcher { return On("request", rules...) }

// OnMetaEvent 元事件触发器
func (e *engineinstance) OnMetaEvent(rules ...zero.Rule) Matcher { return On("meta_event", rules...) }

// OnPrefix 前缀触发器
func (e *engineinstance) OnPrefix(prefix string, rules ...zero.Rule) Matcher {
	return (*matcherinstance)(e.en.OnPrefix(prefix, rules...).SetPriority(e.prio))
}

// OnSuffix 后缀触发器
func (e *engineinstance) OnSuffix(suffix string, rules ...zero.Rule) Matcher {
	return (*matcherinstance)(e.en.OnSuffix(suffix, rules...).SetPriority(e.prio))
}

// OnCommand 命令触发器
func (e *engineinstance) OnCommand(commands string, rules ...zero.Rule) Matcher {
	return (*matcherinstance)(e.en.OnCommand(commands, rules...).SetPriority(e.prio))
}

// OnRegex 正则触发器
func (e *engineinstance) OnRegex(regexPattern string, rules ...zero.Rule) Matcher {
	return (*matcherinstance)(e.en.OnRegex(regexPattern, rules...).SetPriority(e.prio))
}

// OnKeyword 关键词触发器
func (e *engineinstance) OnKeyword(keyword string, rules ...zero.Rule) Matcher {
	return (*matcherinstance)(e.en.OnKeyword(keyword, rules...).SetPriority(e.prio))
}

// OnFullMatch 完全匹配触发器
func (e *engineinstance) OnFullMatch(src string, rules ...zero.Rule) Matcher {
	return (*matcherinstance)(e.en.OnFullMatch(src, rules...).SetPriority(e.prio))
}

// OnFullMatchGroup 完全匹配触发器组
func (e *engineinstance) OnFullMatchGroup(src []string, rules ...zero.Rule) Matcher {
	return (*matcherinstance)(e.en.OnFullMatchGroup(src, rules...).SetPriority(e.prio))
}

// OnKeywordGroup 关键词触发器组
func (e *engineinstance) OnKeywordGroup(keywords []string, rules ...zero.Rule) Matcher {
	return (*matcherinstance)(e.en.OnKeywordGroup(keywords, rules...).SetPriority(e.prio))
}

// OnCommandGroup 命令触发器组
func (e *engineinstance) OnCommandGroup(commands []string, rules ...zero.Rule) Matcher {
	return (*matcherinstance)(e.en.OnCommandGroup(commands, rules...).SetPriority(e.prio))
}

// OnPrefixGroup 前缀触发器组
func (e *engineinstance) OnPrefixGroup(prefix []string, rules ...zero.Rule) Matcher {
	return (*matcherinstance)(e.en.OnPrefixGroup(prefix, rules...).SetPriority(e.prio))
}

// OnSuffixGroup 后缀触发器组
func (e *engineinstance) OnSuffixGroup(suffix []string, rules ...zero.Rule) Matcher {
	return (*matcherinstance)(e.en.OnSuffixGroup(suffix, rules...).SetPriority(e.prio))
}

// OnShell shell命令触发器
func (e *engineinstance) OnShell(command string, model any, rules ...zero.Rule) Matcher {
	return (*matcherinstance)(e.en.OnShell(command, model, rules...).SetPriority(e.prio))
}
