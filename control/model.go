package control

import zero "github.com/wdvxdr1123/ZeroBot"

// grpcfg holds the group config for the Manager.
type grpcfg struct {
	GroupID int64 `db:"gid"`     // GroupID 群号
	Disable int64 `db:"disable"` // Disable 默认启用该插件
}

type ban struct {
	ID      int64 `db:"id"`
	UserID  int64 `db:"uid"`
	GroupID int64 `db:"gid"`
}

type block struct {
	UserID int64 `db:"uid"`
}

// Options holds the optional parameters for the Manager.
type Options struct {
	DisableOnDefault  bool
	Help              string              // 帮助文本信息
	PrivateDataFolder string              // 全部小写的数据文件夹名, 不出现在 zbpdata
	PublicDataFolder  string              // 驼峰的数据文件夹名, 出现在 zbpdata
	OnEnable          func(ctx *zero.Ctx) // 启用插件后执行的命令, 为空则打印 “已启用服务: xxx”
	OnDisable         func(ctx *zero.Ctx) // 禁用插件后执行的命令, 为空则打印 “已禁用服务: xxx”
}
