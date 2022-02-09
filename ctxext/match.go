package ctxext

import (
	"time"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type ListGetter interface {
	List() []string
}

// FirstValueInList 判断正则匹配的第一个参数是否在列表中
func FirstValueInList(list ListGetter) zero.Rule {
	return func(ctx *zero.Ctx) bool {
		first := ctx.State["regex_matched"].([]string)[1]
		for _, v := range list.List() {
			if first == v {
				return true
			}
		}
		return false
	}
}

// IsPicExists 消息含有图片返回 true
func IsPicExists(ctx *zero.Ctx) bool {
	var urls = []string{}
	for _, elem := range ctx.Event.Message {
		if elem.Type == "image" {
			urls = append(urls, elem.Data["url"])
		}
	}
	if len(urls) > 0 {
		ctx.State["image_url"] = urls
		return true
	}
	return false
}

// MustProvidePicture 消息不存在图片阻塞60秒至有图片，超时返回 false
func MustProvidePicture(ctx *zero.Ctx) bool {
	if IsPicExists(ctx) {
		return true
	}
	// 没有图片就索取
	ctx.SendChain(message.Text("请发送一张图片"))
	next := zero.NewFutureEvent("message", 999, false, zero.CheckUser(ctx.Event.UserID), IsPicExists)
	recv, cancel := next.Repeat()
	select {
	case <-time.After(time.Second * 120):
		return false
	case e := <-recv:
		cancel()
		newCtx := &zero.Ctx{Event: e, State: zero.State{}}
		if IsPicExists(newCtx) {
			ctx.State["image_url"] = newCtx.State["image_url"]
			ctx.Event.MessageID = newCtx.Event.MessageID
			return true
		}
		return false
	}
}
