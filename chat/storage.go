package chat

import (
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
)

const (
	BitmapRate = 0x0000ff
	BitmapTemp = 0x00ff00
	BitmapNagt = 0x010000
	BitmapNrec = 0x020000
	BitmapNrat = 0x040000
)

type Storage ctxext.Storage

func NewStorage(ctx *zero.Ctx, gid int64) (Storage, error) {
	s, err := ctxext.NewStorage(ctx, gid)
	return Storage(s), err
}

func (s Storage) Rate() uint8 {
	return uint8((ctxext.Storage)(s).Get(BitmapRate))
}

func (s Storage) Temp() float32 {
	temp := int8((ctxext.Storage)(s).Get(BitmapTemp))
	// 处理温度参数
	if temp <= 0 {
		temp = 70 // default setting
	}
	if temp > 100 {
		temp = 100
	}
	return float32(temp) / 100
}

func (s Storage) NoAgent() bool {
	return (ctxext.Storage)(s).GetBool(BitmapNagt)
}

func (s Storage) NoRecord() bool {
	return (ctxext.Storage)(s).GetBool(BitmapNrec)
}

func (s Storage) NoReplyAt() bool {
	return (ctxext.Storage)(s).GetBool(BitmapNrat)
}
