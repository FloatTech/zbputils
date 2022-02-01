package control

import (
	"strconv"

	zero "github.com/wdvxdr1123/ZeroBot"
)

func initBlock() (err error) {
	err = db.Create("block__", &block{})
	if err == nil {
		zero.OnMessage(func(ctx *zero.Ctx) bool {
			return isBlocked(ctx.Event.UserID)
		}).SetBlock(true).SecondPriority()
	}
	return
}

func doBlock(uid int64) error {
	return db.Insert("block__", &block{UserID: uid})
}

func doUnblock(uid int64) error {
	return db.Del("block__", "where uid = "+strconv.FormatInt(uid, 10))
}

func isBlocked(uid int64) bool {
	return db.CanFind("block__", "where uid = "+strconv.FormatInt(uid, 10))
}
