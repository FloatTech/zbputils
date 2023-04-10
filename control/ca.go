package control

import (
	"math"
	"math/rand"
	"time"

	"github.com/RomiChan/syncx"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	"github.com/FloatTech/floatbox/process"
	"github.com/FloatTech/ttl"
)

type ca struct {
	hasConflict *ttl.Cache[int64, bool]
	withdraw    syncx.Map[int64, uint8]
}

func (c *ca) handle(ctx *zero.Ctx) bool {
	if !zero.OnlyGroup(ctx) {
		return true
	}
	if c.hasConflict.Get(ctx.Event.GroupID) {
		return false
	}
	delaymax, ok := c.withdraw.Load(ctx.Event.GroupID)
	if !ok {
		return true
	}
	if delaymax == 0 {
		delaymax = math.MaxUint8
	}
	slptm := rand.Intn(int(delaymax))
	time.Sleep(time.Millisecond * 100 * time.Duration(slptm))
	if c.hasConflict.Get(ctx.Event.GroupID) {
		delaymax -= uint8(slptm)
		c.withdraw.Store(ctx.Event.GroupID, delaymax)
		c.hasConflict.Delete(ctx.Event.GroupID)
		return false
	}
	go func() {
		tok := genToken()
		t := message.Text("●ca" + tok)
		id := ctx.SendChain(t)
		process.SleepAbout1sTo2s()
		ctx.DeleteMessage(id)
	}()
	return true
}

var conflicts = ca{
	hasConflict: ttl.NewCache[int64, bool](time.Millisecond * 100 * math.MaxUint8),
}

func init() {
	zero.OnRegex("^(启|禁)用插件冲突避免$", zero.OnlyGroup, zero.AdminPermission, zero.OnlyToMe).SetBlock(true).SecondPriority().
		Handle(func(ctx *zero.Ctx) {
			cmd := ctx.State["regex_matched"].([]string)[1]
			switch cmd {
			case "启":
				conflicts.withdraw.Store(ctx.Event.GroupID, math.MaxUint8)
			case "禁":
				conflicts.withdraw.Delete(ctx.Event.GroupID)
			}
			ctx.SendChain(message.Text("成功", cmd, "用"))
		})

	zero.OnRegex("^●ca([\u4e00-\u8e00]{4})$", zero.OnlyGroup).SetBlock(true).SecondPriority().
		Handle(func(ctx *zero.Ctx) {
			if isValidToken(ctx.State["regex_matched"].([]string)[1]) {
				conflicts.hasConflict.Set(ctx.Event.GroupID, true)
			}
		})
}
