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

func doBlock(uid int64) error {
	return db.Insert("__block", &block{UserID: uid})
}

func doUnblock(uid int64) error {
	return db.Del("__block", "where uid = "+strconv.FormatInt(uid, 10))
}

func isBlocked(uid int64) bool {
	return db.CanFind("__block", "where uid = "+strconv.FormatInt(uid, 10))
}
