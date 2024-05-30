package pool

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
)

type NTImage nturl

func NewNTImage(u string) (nti NTImage, err error) {
	subs := ntcachere.FindStringSubmatch(u)
	if len(subs) != 3 {
		err = ErrInvalidNTURL
		return
	}
	nti = NTImage(u)
	return
}

func (nti NTImage) String() string {
	subs := ntcachere.FindStringSubmatch(string(nti))
	if len(subs) != 3 {
		panic(ErrInvalidNTURL)
	}
	fileid := subs[1]
	rkey, err := rs.rkey(time.Minute)
	if err != nil || rkey == "" {
		rkey = subs[2]
	}
	return fmt.Sprintf(ntcacheurl, fileid, rkey)
}

func init() {
	zero.OnMessage(zero.HasPicture).SetBlock(false).FirstPriority().Handle(func(ctx *zero.Ctx) {
		img, ok := ctx.State["image_url"].([]string)
		if !ok || len(img) == 0 {
			return
		}
		if !ntcachere.MatchString(img[0]) { // is not NTQQ
			return
		}
		rk, err := nturl(img[0]).rkey()
		if err != nil {
			logrus.Debugln("[imgpool] parse rkey error:", err, "image url:", img)
			return
		}
		err = rs.set(time.Minute, rk)
		if err != nil {
			logrus.Debugln("[imgpool] set rkey error:", err, "rkey:", rk)
			return
		}
		logrus.Debugln("[imgpool] set latest rkey:", rk)
	})
}
