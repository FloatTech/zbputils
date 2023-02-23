// Package types 结构体
package types

// BotParams GetGroupList,GetFriendList的入参
type BotParams struct {
	SelfID int64 `json:"selfId" form:"selfId" validate:"required"`
}

// AllPluginParams GetAllPlugin的入参
type AllPluginParams struct {
	GroupID int64 `json:"groupId" form:"groupId"`
}

// PluginParams GetPlugin的入参
type PluginParams struct {
	GroupID int64  `json:"groupId" form:"groupId"`
	Name    string `json:"name" form:"name"`
}

// PluginStatusParams UpdatePluginStatus的入参, Status,0=禁用,1=启用,2=还原
type PluginStatusParams struct {
	GroupID int64  `json:"groupId" form:"groupId"`
	Name    string `json:"name" form:"name" validate:"required"`
	Status  int    `json:"status" form:"status"`
}

// PluginStatusParams UpdatePluginStatus的入参, Status,0=沉默,1=响应
type ResponseStatusParams struct {
	GroupID int64 `json:"groupId" form:"groupId"`
	Status  int   `json:"status" form:"status"`
}

// AllPluginStatusParams UpdateAllPluginStatus的入参
type AllPluginStatusParams struct {
	GroupID int64 `json:"groupId" form:"groupId"`
	Status  int   `json:"status" form:"status"`
}

// HandleRequestParams 处理事件的入参
type HandleRequestParams struct {
	Flag    string `json:"flag" form:"flag"`
	Reason  string `json:"reason" form:"reason"`
	Approve bool   `json:"approve" form:"approve"`
}

// SendMsgParams 发送消息的入参
type SendMsgParams struct {
	SelfID  int64   `json:"selfId" form:"selfId"`
	GIDList []int64 `json:"gidList" form:"gidList"`
	Message string  `json:"message" form:"message"`
}

// LoginParams 登录参数
type LoginParams struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

// LoginResultVo 登录返回参数
type LoginResultVo struct {
	UserID   int        `json:"userId"`
	Username string     `json:"username"`
	RealName string     `json:"realName"`
	Desc     string     `json:"desc"`
	Token    string     `json:"token"`
	Roles    []RoleInfo `json:"roles"`
}

// RoleInfo 角色参数
type RoleInfo struct {
	RoleName string `json:"roleName"`
	Value    string `json:"value"`
}

// UserInfoVo 用户信息
type UserInfoVo struct {
	UserID   int        `json:"userId"`
	Username string     `json:"username"`
	RealName string     `json:"realName"`
	Desc     string     `json:"desc"`
	Token    string     `json:"token"`
	Roles    []RoleInfo `json:"roles"`
	Avatar   string     `json:"avatar"`
	HomePath string     `json:"homePath"`
	Password string     `json:"password"`
}

// MessageInfo 消息信息
type MessageInfo struct {
	MessageType string      `json:"message_type"`
	MessageID   interface{} `json:"message_id"`
	GroupID     int64       `json:"group_id"`
	GroupName   string      `json:"group_name"`
	UserID      int64       `json:"user_id"`
	Nickname    string      `json:"nickname"`
	RawMessage  string      `json:"raw_message"`
}

// PluginVo 全部插件的返回
type PluginVo struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Brief          string `json:"brief"`
	Usage          string `json:"usage"`
	Banner         string `json:"banner"`
	PluginStatus   bool   `json:"pluginStatus"`
	ResponseStatus bool   `json:"responseStatus"`
}
