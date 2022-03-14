package ctxext

import (
	zero "github.com/wdvxdr1123/ZeroBot"
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
