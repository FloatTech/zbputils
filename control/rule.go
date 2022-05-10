// Package control 控制插件的启用与优先级等
package control

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension"
	"github.com/wdvxdr1123/ZeroBot/extension/ttl"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"

	sql "github.com/FloatTech/sqlite"
	"github.com/FloatTech/zbputils/img/text"
	"github.com/FloatTech/zbputils/process"
)

var (
	db = &sql.Sqlite{DBPath: "data/control/plugins.db"}
	// managers 每个插件对应的管理
	managers  = map[string]*Control{}
	manmu     sync.RWMutex
	ctxbanmap = ttl.NewCache[*zero.Ctx, bool](10 * time.Second)
)

// Control is to control the plugins.
type Control struct {
	sync.RWMutex
	service string
	cache   map[int64]bool
	options Options
}

// newctrl returns Manager with settings.
func newctrl(service string, o *Options) *Control {
	m := &Control{
		service: service,
		cache:   make(map[int64]bool, 16),
		options: func() Options {
			if o == nil {
				return Options{}
			}
			return *o
		}(),
	}
	manmu.Lock()
	managers[service] = m
	manmu.Unlock()
	err := db.Create(service, &grpcfg{})
	if err != nil {
		panic(err)
	}
	err = db.Create(service+"ban", &ban{})
	if err != nil {
		panic(err)
	}
	return m
}

// Enable enables a group to pass the Manager.
// groupID == 0 (ALL) will operate on all grps.
func (m *Control) Enable(groupID int64) {
	var c grpcfg
	m.RLock()
	err := db.Find(m.service, &c, "WHERE gid="+strconv.FormatInt(groupID, 10))
	m.RUnlock()
	if err != nil {
		c.GroupID = groupID
	}
	c.Disable = int64(uint64(c.Disable) & 0xffffffff_fffffffe)
	m.Lock()
	m.cache[groupID] = true
	err = db.Insert(m.service, &c)
	m.Unlock()
	if err != nil {
		log.Errorf("[control] %v", err)
	}
}

// Disable disables a group to pass the Manager.
// groupID == 0 (ALL) will operate on all grps.
func (m *Control) Disable(groupID int64) {
	var c grpcfg
	m.RLock()
	err := db.Find(m.service, &c, "WHERE gid="+strconv.FormatInt(groupID, 10))
	m.RUnlock()
	if err != nil {
		c.GroupID = groupID
	}
	c.Disable |= 1
	m.Lock()
	m.cache[groupID] = false
	err = db.Insert(m.service, &c)
	m.Unlock()
	if err != nil {
		log.Errorf("[control] %v", err)
	}
}

// Reset resets the default config of a group.
// groupID == 0 (ALL) is not allowed.
func (m *Control) Reset(groupID int64) {
	if groupID != 0 {
		m.Lock()
		delete(m.cache, groupID)
		err := db.Del(m.service, "WHERE gid="+strconv.FormatInt(groupID, 10))
		m.Unlock()
		if err != nil {
			log.Errorf("[control] %v", err)
		}
	}
}

// IsEnabledIn 查询开启群
// 当全局未配置或与默认相同时, 状态取决于单独配置, 后备为默认配置；
// 当全局与默认不同时, 状态取决于全局配置, 单独配置失效。
func (m *Control) IsEnabledIn(gid int64) bool {
	var c grpcfg
	var err error

	m.RLock()
	yes, ok := m.cache[0]
	m.RUnlock()
	if !ok {
		m.RLock()
		err = db.Find(m.service, &c, "WHERE gid=0")
		m.RUnlock()
		if err == nil && c.GroupID == 0 {
			log.Debugf("[control] plugin %s of all : %d", m.service, c.Disable&1)
			yes = c.Disable&1 == 0
			ok = true
			m.Lock()
			m.cache[0] = yes
			m.Unlock()
			log.Debugf("[control] cache plugin %s of grp %d : %v", m.service, gid, yes)
		}
	}

	if ok && yes == m.options.DisableOnDefault { // global enable status is different from default value
		return yes
	}

	m.RLock()
	yes, ok = m.cache[gid]
	m.RUnlock()
	if ok {
		log.Debugf("[control] read cached %s of grp %d : %v", m.service, gid, yes)
	} else {
		m.RLock()
		err = db.Find(m.service, &c, "WHERE gid="+strconv.FormatInt(gid, 10))
		m.RUnlock()
		if err == nil && gid == c.GroupID {
			log.Debugf("[control] plugin %s of grp %d : %d", m.service, c.GroupID, c.Disable&1)
			yes = c.Disable&1 == 0
			ok = true
			m.Lock()
			m.cache[gid] = yes
			m.Unlock()
			log.Debugf("[control] cache plugin %s of grp %d : %v", m.service, gid, yes)
		}
	}

	if ok {
		return yes
	}
	return !m.options.DisableOnDefault
}

// Ban 禁止某人在某群使用本插件
func (m *Control) Ban(uid, gid int64) {
	var err error
	var digest [16]byte
	if gid != 0 { // 特定群
		digest = md5.Sum(helper.StringToBytes(fmt.Sprintf("%d_%d", uid, gid)))
		m.RLock()
		err = db.Insert(m.service+"ban", &ban{ID: int64(binary.LittleEndian.Uint64(digest[:8])), UserID: uid, GroupID: gid})
		m.RUnlock()
		if err == nil {
			log.Debugf("[control] plugin %s is banned in grp %d for usr %d.", m.service, gid, uid)
			return
		}
	}
	// 所有群
	digest = md5.Sum(helper.StringToBytes(fmt.Sprintf("%d_all", uid)))
	m.RLock()
	err = db.Insert(m.service+"ban", &ban{ID: int64(binary.LittleEndian.Uint64(digest[:8])), UserID: uid, GroupID: 0})
	m.RUnlock()
	if err == nil {
		log.Debugf("[control] plugin %s is banned in all grp for usr %d.", m.service, uid)
	}
}

// Permit 允许某人在某群使用本插件
func (m *Control) Permit(uid, gid int64) {
	var digest [16]byte
	if gid != 0 { // 特定群
		digest = md5.Sum(helper.StringToBytes(fmt.Sprintf("%d_%d", uid, gid)))
		m.RLock()
		_ = db.Del(m.service+"ban", "WHERE id = "+strconv.FormatInt(int64(binary.LittleEndian.Uint64(digest[:8])), 10))
		m.RUnlock()
		log.Debugf("[control] plugin %s is permitted in grp %d for usr %d.", m.service, gid, uid)
		return
	}
	// 所有群
	digest = md5.Sum(helper.StringToBytes(fmt.Sprintf("%d_all", uid)))
	m.RLock()
	_ = db.Del(m.service+"ban", "WHERE id = "+strconv.FormatInt(int64(binary.LittleEndian.Uint64(digest[:8])), 10))
	m.RUnlock()
	log.Debugf("[control] plugin %s is permitted in all grp for usr %d.", m.service, uid)
}

// IsBannedIn 某人是否在某群被 ban
func (m *Control) IsBannedIn(uid, gid int64) bool {
	var b ban
	var err error
	var digest [16]byte
	if gid != 0 {
		digest = md5.Sum(helper.StringToBytes(fmt.Sprintf("%d_%d", uid, gid)))
		m.RLock()
		err = db.Find(m.service+"ban", &b, "WHERE id = "+strconv.FormatInt(int64(binary.LittleEndian.Uint64(digest[:8])), 10))
		m.RUnlock()
		if err == nil && gid == b.GroupID && uid == b.UserID {
			log.Debugf("[control] plugin %s is banned in grp %d for usr %d.", m.service, b.GroupID, b.UserID)
			return true
		}
	}
	digest = md5.Sum(helper.StringToBytes(fmt.Sprintf("%d_all", uid)))
	m.RLock()
	err = db.Find(m.service+"ban", &b, "WHERE id = "+strconv.FormatInt(int64(binary.LittleEndian.Uint64(digest[:8])), 10))
	m.RUnlock()
	if err == nil && b.GroupID == 0 && uid == b.UserID {
		log.Debugf("[control] plugin %s is banned in all grp for usr %d.", m.service, b.UserID)
		return true
	}
	return false
}

// GetData 获取某个群的 63 字节配置信息
func (m *Control) GetData(gid int64) int64 {
	var c grpcfg
	var err error
	m.RLock()
	err = db.Find(m.service, &c, "WHERE gid="+strconv.FormatInt(gid, 10))
	m.RUnlock()
	if err == nil && gid == c.GroupID {
		log.Debugf("[control] plugin %s of grp %d : 0x%x", m.service, c.GroupID, c.Disable>>1)
		return c.Disable >> 1
	}
	return 0
}

// SetData 为某个群设置低 63 位配置数据
func (m *Control) SetData(groupID int64, data int64) error {
	var c grpcfg
	m.RLock()
	err := db.Find(m.service, &c, "WHERE gid="+strconv.FormatInt(groupID, 10))
	m.RUnlock()
	if err != nil {
		c.GroupID = groupID
		if m.options.DisableOnDefault {
			c.Disable = 1
		}
	}
	c.Disable &= 1
	c.Disable |= data << 1
	log.Debugf("[control] set plugin %s of grp %d : 0x%x", m.service, c.GroupID, data)
	m.Lock()
	err = db.Insert(m.service, &c)
	m.Unlock()
	if err != nil {
		log.Errorf("[control] %v", err)
	}
	return err
}

// Handler 返回 预处理器
func (m *Control) Handler(ctx *zero.Ctx) bool {
	ctx.State["manager"] = m
	grp := ctx.Event.GroupID
	if grp == 0 {
		// 个人用户
		grp = -ctx.Event.UserID
	}
	ok := ctxbanmap.Get(ctx)
	if ok {
		return m.IsEnabledIn(grp)
	}
	isnotbanned := !m.IsBannedIn(ctx.Event.UserID, grp)
	if isnotbanned {
		ctxbanmap.Set(ctx, true)
		return m.IsEnabledIn(grp)
	}
	return false
}

// String 打印帮助
func (m *Control) String() string {
	return m.options.Help
}

// EnableMark 启用：●，禁用：○
type EnableMark bool

// String 打印启用状态
func (em EnableMark) String() string {
	if bool(em) {
		return "●"
	}
	return "○"
}

// EnableMarkIn 打印 ● 或 ○
func (m *Control) EnableMarkIn(grp int64) EnableMark {
	return EnableMark(m.IsEnabledIn(grp))
}

// Lookup returns a Manager by the service name, if
// not exist, it will return nil.
func Lookup(service string) (*Control, bool) {
	manmu.RLock()
	m, ok := managers[service]
	manmu.RUnlock()
	return m, ok
}

// ForEach iterates through managers.
func ForEach(iterator func(key string, manager *Control) bool) {
	manmu.RLock()
	m := copyMap(managers)
	manmu.RUnlock()
	for k, v := range m {
		if !iterator(k, v) {
			return
		}
	}
}

func copyMap(m map[string]*Control) map[string]*Control {
	ret := make(map[string]*Control, len(m))
	for k, v := range m {
		ret[k] = v
	}
	return ret
}

func init() {
	process.NewCustomOnce(&manmu).Do(func() {
		err := os.MkdirAll("data/control", 0755)
		if err != nil {
			panic(err)
		}
		err = initBlock()
		if err != nil {
			panic(err)
		}
		zero.OnCommandGroup([]string{
			"启用", "enable", "禁用", "disable",
		}, zero.UserOrGrpAdmin).SetBlock(true).SecondPriority().Handle(func(ctx *zero.Ctx) {
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			service, ok := Lookup(model.Args)
			if !ok {
				ctx.SendChain(message.Text("没有找到指定服务!"))
				return
			}
			grp := ctx.Event.GroupID
			if grp == 0 {
				// 个人用户
				grp = -ctx.Event.UserID
			}
			if strings.Contains(model.Command, "启用") || strings.Contains(model.Command, "enable") {
				service.Enable(grp)
				if service.options.OnEnable != nil {
					service.options.OnEnable(ctx)
				} else {
					ctx.SendChain(message.Text("已启用服务: " + model.Args))
				}
			} else {
				service.Disable(grp)
				if service.options.OnDisable != nil {
					service.options.OnDisable(ctx)
				} else {
					ctx.SendChain(message.Text("已禁用服务: " + model.Args))
				}
			}
		})

		zero.OnCommandGroup([]string{
			"全局启用", "allenable", "全局禁用", "alldisable",
		}, zero.OnlyToMe, zero.SuperUserPermission).SetBlock(true).SecondPriority().Handle(func(ctx *zero.Ctx) {
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			service, ok := Lookup(model.Args)
			if !ok {
				ctx.SendChain(message.Text("没有找到指定服务!"))
				return
			}
			if strings.Contains(model.Command, "启用") || strings.Contains(model.Command, "enable") {
				service.Enable(0)
				ctx.SendChain(message.Text("已全局启用服务: " + model.Args))
			} else {
				service.Disable(0)
				ctx.SendChain(message.Text("已全局禁用服务: " + model.Args))
			}
		})

		zero.OnCommandGroup([]string{"还原", "reset"}, zero.UserOrGrpAdmin).SetBlock(true).SecondPriority().Handle(func(ctx *zero.Ctx) {
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			service, ok := Lookup(model.Args)
			if !ok {
				ctx.SendChain(message.Text("没有找到指定服务!"))
				return
			}
			grp := ctx.Event.GroupID
			if grp == 0 {
				// 个人用户
				grp = -ctx.Event.UserID
			}
			service.Reset(grp)
			ctx.SendChain(message.Text("已还原服务的默认启用状态: " + model.Args))
		})

		zero.OnCommandGroup([]string{
			"禁止", "ban", "允许", "permit",
		}, zero.AdminPermission).SetBlock(true).SecondPriority().Handle(func(ctx *zero.Ctx) {
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			args := strings.Split(model.Args, " ")
			if len(args) >= 2 {
				service, ok := Lookup(args[0])
				if !ok {
					ctx.SendChain(message.Text("没有找到指定服务!"))
					return
				}
				grp := ctx.Event.GroupID
				if grp == 0 {
					grp = -ctx.Event.UserID
				}
				msg := "**" + args[0] + "报告**"
				if strings.Contains(model.Command, "允许") || strings.Contains(model.Command, "permit") {
					for _, usr := range args[1:] {
						uid, err := strconv.ParseInt(usr, 10, 64)
						if err == nil {
							service.Permit(uid, grp)
							msg += "\n+ 已允许" + usr
						}
					}
				} else {
					for _, usr := range args[1:] {
						uid, err := strconv.ParseInt(usr, 10, 64)
						if err == nil {
							service.Ban(uid, grp)
							msg += "\n- 已禁止" + usr
						}
					}
				}
				ctx.SendChain(message.Text(msg))
				return
			}
			ctx.SendChain(message.Text("参数错误!"))
		})

		zero.OnCommandGroup([]string{
			"全局禁止", "allban", "全局允许", "allpermit",
		}, zero.SuperUserPermission).SetBlock(true).SecondPriority().Handle(func(ctx *zero.Ctx) {
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			args := strings.Split(model.Args, " ")
			if len(args) >= 2 {
				service, ok := Lookup(args[0])
				if !ok {
					ctx.SendChain(message.Text("没有找到指定服务!"))
					return
				}
				msg := "**" + args[0] + "全局报告**"
				if strings.Contains(model.Command, "允许") || strings.Contains(model.Command, "permit") {
					for _, usr := range args[1:] {
						uid, err := strconv.ParseInt(usr, 10, 64)
						if err == nil {
							service.Permit(uid, 0)
							msg += "\n+ 已允许" + usr
						}
					}
				} else {
					for _, usr := range args[1:] {
						uid, err := strconv.ParseInt(usr, 10, 64)
						if err == nil {
							service.Ban(uid, 0)
							msg += "\n- 已禁止" + usr
						}
					}
				}
				ctx.SendChain(message.Text(msg))
				return
			}
			ctx.SendChain(message.Text("参数错误!"))
		})

		zero.OnCommandGroup([]string{
			"封禁", "block", "解封", "unblock",
		}, zero.SuperUserPermission).SetBlock(true).SecondPriority().Handle(func(ctx *zero.Ctx) {
			model := extension.CommandModel{}
			_ = ctx.Parse(&model)
			args := strings.Split(model.Args, " ")
			if len(args) >= 1 {
				msg := "**报告**"
				if strings.Contains(model.Command, "解") || strings.Contains(model.Command, "un") {
					for _, usr := range args {
						uid, err := strconv.ParseInt(usr, 10, 64)
						if err == nil {
							if doUnblock(uid) == nil {
								msg += "\n- 已解封" + usr
							}
						}
					}
				} else {
					for _, usr := range args {
						uid, err := strconv.ParseInt(usr, 10, 64)
						if err == nil {
							if doBlock(uid) == nil {
								msg += "\n+ 已封禁" + usr
							}
						}
					}
				}
				ctx.SendChain(message.Text(msg))
				return
			}
			ctx.SendChain(message.Text("参数错误!"))
		})

		zero.OnCommandGroup([]string{"用法", "usage"}, zero.UserOrGrpAdmin).SetBlock(true).SecondPriority().
			Handle(func(ctx *zero.Ctx) {
				model := extension.CommandModel{}
				_ = ctx.Parse(&model)
				service, ok := Lookup(model.Args)
				if !ok {
					ctx.SendChain(message.Text("没有找到指定服务!"))
					return
				}
				if service.options.Help != "" {
					gid := ctx.Event.GroupID
					if gid == 0 {
						gid = -ctx.Event.UserID
					}
					ctx.SendChain(message.Text(service.EnableMarkIn(gid), " ", service))
				} else {
					ctx.SendChain(message.Text("该服务无帮助!"))
				}
			})

		zero.OnCommandGroup([]string{"服务列表", "service_list"}, zero.UserOrGrpAdmin).SetBlock(true).SecondPriority().
			Handle(func(ctx *zero.Ctx) {
				i := 0
				gid := ctx.Event.GroupID
				if gid == 0 {
					gid = -ctx.Event.UserID
				}
				manmu.RLock()
				msg := make([]any, 1, len(managers)*4+1)
				manmu.RUnlock()
				msg[0] = "--------服务列表--------\n发送\"/用法 name\"查看详情"
				ForEach(func(key string, manager *Control) bool {
					i++
					msg = append(msg, "\n", i, ": ", manager.EnableMarkIn(gid), key)
					return true
				})
				ctx.Send(message.Text(msg...))
			})

		zero.OnCommandGroup([]string{"服务详情", "service_detail"}, zero.UserOrGrpAdmin).SetBlock(true).SecondPriority().
			Handle(func(ctx *zero.Ctx) {
				i := 0
				gid := ctx.Event.GroupID
				if gid == 0 {
					gid = -ctx.Event.UserID
				}
				manmu.RLock()
				msgs := make([]any, 1, len(managers)*7+1)
				manmu.RUnlock()
				msgs[0] = "---服务详情---\n"
				ForEach(func(key string, service *Control) bool {
					i++
					msgs = append(msgs, i, ": ", service.EnableMarkIn(gid), key, "\n", service, "\n\n")
					return true
				})
				data, err := text.RenderToBase64(fmt.Sprint(msgs...), text.FontFile, 400, 20)
				if err != nil {
					log.Errorf("[control] %v", err)
				}
				if id := ctx.SendChain(message.Image("base64://" + helper.BytesToString(data))); id.ID() == 0 {
					ctx.SendChain(message.Text("ERROR: 可能被风控了"))
				}
			})
	})
}
