package control

import (
	"strconv"

	zero "github.com/wdvxdr1123/ZeroBot"
)

func initBlock() (err error) {
	err = db.Create("__block", &block{})
	if err == nil {
		zero.OnMessage(func(ctx *zero.Ctx) bool {
			return isBlocked(ctx.Event.UserID)
		}).SetBlock(true).ThirdPriority()
	}
	return
}

var blockCache = make(map[int64]bool)

func doBlock(uid int64) error {
	manmu.Lock()
	defer manmu.Unlock()
	blockCache[uid] = true
	return db.Insert("__block", &block{UserID: uid})
}

func doUnblock(uid int64) error {
	manmu.Lock()
	defer manmu.Unlock()
	blockCache[uid] = false
	return db.Del("__block", "where uid = "+strconv.FormatInt(uid, 10))
}

func isBlocked(uid int64) bool {
	manmu.RLock()
	isbl, ok := blockCache[uid]
	manmu.RUnlock()
	if ok {
		return isbl
	}
	manmu.Lock()
	defer manmu.Unlock()
	isbl = db.CanFind("__block", "where uid = "+strconv.FormatInt(uid, 10))
	blockCache[uid] = isbl
	return isbl
}
