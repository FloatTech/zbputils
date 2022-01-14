package ctxext

import zero "github.com/wdvxdr1123/ZeroBot"

// UserOrGrpAdmin 允许用户单独使用或群管使用
func UserOrGrpAdmin(ctx *zero.Ctx) bool {
	if zero.OnlyGroup(ctx) {
		return zero.AdminPermission(ctx)
	}
	return zero.OnlyToMe(ctx)
}
