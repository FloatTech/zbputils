// Package controller 主要处理逻辑
package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/FloatTech/floatbox/binary"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/control/web/middleware"
	"github.com/FloatTech/zbputils/control/web/types"
	"github.com/RomiChan/syncx"
	"github.com/RomiChan/websocket"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	// MsgConn 向前端推送消息的ws链接
	MsgConn *websocket.Conn
	// LogConn 向前端推送日志的ws链接
	LogConn *websocket.Conn
	// 实现Write接口
	l logWriter
	// 存储请求事件，flag作为键，一个request对象作为值
	requestData syncx.Map[string, *zero.Event]
	upgrader    = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	validate = validator.New()
)

func init() {
	// 日志设置
	writer := io.MultiWriter(l, os.Stdout)
	log.SetOutput(writer)
	log.SetFormatter(&log.TextFormatter{DisableColors: false})

	zero.OnMessage().SetBlock(false).FirstPriority().Handle(func(ctx *zero.Ctx) {
		if MsgConn != nil {
			mi := types.MessageInfo{
				MessageType: ctx.Event.MessageType,
				MessageID:   ctx.Event.MessageID,
				GroupID:     ctx.Event.GroupID,
				GroupName:   ctx.GetGroupInfo(ctx.Event.GroupID, false).Name,
				UserID:      ctx.Event.UserID,
				Nickname:    ctx.GetStrangerInfo(ctx.Event.UserID, false).Get("nickname").String(),
				RawMessage:  ctx.Event.RawMessage,
			}
			err := MsgConn.WriteJSON(mi)
			if err != nil {
				log.Errorln("[gui] 推送消息发送错误:", err)
				return
			}
		}
	})
	// 直接注册一个request请求监听器，优先级设置为最高，设置不阻断事件传播
	zero.OnRequest().SetBlock(false).FirstPriority().Handle(func(ctx *zero.Ctx) {
		requestData.Store(ctx.Event.Flag, ctx.Event)
	})
}

// logWriter
type logWriter struct {
	// 日志连接
}

// GetBotList 获取机器人qq号
// @Description 获取机器人qq号
// @Router /api/getBotList [get]
func GetBotList(context *gin.Context) {
	var bots []int64
	zero.RangeBot(func(id int64, ctx *zero.Ctx) bool {
		bots = append(bots, id)
		return true
	})
	context.JSON(http.StatusOK, gin.H{
		"code":    0,
		"result":  bots,
		"message": "",
		"type":    "ok",
	})
}

// GetFriendList 获取好友列表
// @Description 获取好友列表
// @Router /api/getFriendList [get]
// @Param selfId query integer false "机器人qq号" default(123456)
func GetFriendList(context *gin.Context) {
	_, bot, err := getBot(context)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"result":  nil,
			"message": err.Error(),
			"type":    "error",
		})
		return
	}
	var resp []any
	list := bot.GetFriendList().String()
	err = json.Unmarshal(binary.StringToBytes(list), &resp)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"result":  nil,
			"message": err.Error(),
			"type":    "error",
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"code":    0,
		"result":  resp,
		"message": "",
		"type":    "ok",
	})
}

// GetGroupList 获取群列表
// @Description 获取群列表
// @Router /api/getGroupList [get]
// @Param selfId query integer false "机器人qq号" default(123456)
func GetGroupList(context *gin.Context) {
	_, bot, err := getBot(context)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"result":  nil,
			"message": err.Error(),
			"type":    "error",
		})
		return
	}
	var resp []any
	list := bot.GetGroupList().String()
	err = json.Unmarshal(binary.StringToBytes(list), &resp)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"result":  nil,
			"message": err.Error(),
			"type":    "error",
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"code":    0,
		"result":  resp,
		"message": "",
		"type":    "ok",
	})
}

// GetAllPlugin 获取所有插件的状态
// @Description 获取所有插件的状态
// @Router /api/manage/getAllPlugin [get]
// @Param groupId query integer false "群号" default(0)
func GetAllPlugin(context *gin.Context) {
	var d types.AllPluginParams
	err := bind(&d, context)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"result":  nil,
			"message": err.Error(),
			"type":    "error",
		})
		return
	}
	var pluginVoList []types.PluginVo
	control.ForEachByPrio(func(i int, manager *ctrl.Control[*zero.Ctx]) bool {
		p := types.PluginVo{
			ID:             i,
			Name:           manager.Service,
			Brief:          manager.Options.Brief,
			Usage:          manager.String(),
			Banner:         "https://gitcode.net/u011570312/zbpbanner/-/raw/main/" + manager.Service + ".png",
			PluginStatus:   manager.IsEnabledIn(d.GroupID),
			ResponseStatus: control.CanResponse(d.GroupID),
		}
		pluginVoList = append(pluginVoList, p)
		return true
	})
	context.JSON(http.StatusOK, gin.H{
		"code":    0,
		"result":  pluginVoList,
		"message": "",
		"type":    "ok",
	})
}

// GetPlugin 获取某个插件的状态
// @Description 获取某个插件的状态
// @Router /api/manage/getPlugin [get]
// @Param groupId query integer false "群号" default(0)
// @Param name query string false "插件名" default(antibuse)
func GetPlugin(context *gin.Context) {
	var d types.PluginParams
	err := bind(&d, context)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"result":  nil,
			"message": err.Error(),
			"type":    "error",
		})
		return
	}
	con, b := control.Lookup(d.Name)
	if !b {
		context.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"result":  nil,
			"message": d.Name + "服务不存在",
			"type":    "error",
		})
		return
	}
	p := types.PluginVo{
		Name:           con.Service,
		Brief:          con.Options.Brief,
		Usage:          con.String(),
		Banner:         "https://gitcode.net/u011570312/zbpbanner/-/raw/main/" + con.Service + ".png",
		PluginStatus:   con.IsEnabledIn(d.GroupID),
		ResponseStatus: control.CanResponse(d.GroupID),
	}
	context.JSON(http.StatusOK, gin.H{
		"code":    0,
		"result":  p,
		"message": "",
		"type":    "ok",
	})
}

// UpdatePluginStatus 更改某一个插件状态
// @Description 更改某一个插件状态
// @Router /api/manage/updatePluginStatus [post]
// @Param object body types.PluginStatusParams false "修改插件状态入参"
func UpdatePluginStatus(context *gin.Context) {
	var d types.PluginStatusParams
	err := bind(&d, context)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"result":  nil,
			"message": err.Error(),
			"type":    "error",
		})
		return
	}
	con, b := control.Lookup(d.Name)
	if !b {
		context.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"result":  nil,
			"message": d.Name + "服务不存在",
			"type":    "error",
		})
		return
	}
	switch d.Status {
	case 0:
		con.Disable(d.GroupID)
	case 1:
		con.Enable(d.GroupID)
	case 2:
		con.Reset(d.GroupID)
	default:
	}
	context.JSON(http.StatusOK, gin.H{
		"code":    0,
		"result":  nil,
		"message": "",
		"type":    "ok",
	})
}

// UpdateResponseStatus 更改某一个群响应
// @Description 更改某一个群响应
// @Router /api/manage/updateResponseStatus [post]
// @Param object body types.ResponseStatusParams false "修改群响应入参"
func UpdateResponseStatus(context *gin.Context) {
	var d types.ResponseStatusParams
	err := bind(&d, context)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"result":  nil,
			"message": err.Error(),
			"type":    "error",
		})
		return
	}
	if d.Status == 1 {
		err = control.Response(d.GroupID)
	} else {
		err = control.Silence(d.GroupID)
	}
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"result":  nil,
			"message": err.Error(),
			"type":    "error",
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"code":    0,
		"result":  nil,
		"message": "",
		"type":    "ok",
	})
}

// UpdateAllPluginStatus 更改某群所有插件状态
// @Description 更改某群所有插件状态
// @Router /api/manage/updateAllPluginStatus [post]
// @Param object body types.AllPluginStatusParams false "修改插件状态入参"
func UpdateAllPluginStatus(context *gin.Context) {
	var d types.AllPluginStatusParams
	err := bind(&d, context)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"result":  nil,
			"message": err.Error(),
			"type":    "error",
		})
		return
	}
	if d.Status == 1 {
		control.ForEachByPrio(func(i int, manager *ctrl.Control[*zero.Ctx]) bool {
			manager.Enable(d.GroupID)
			return true
		})
	} else {
		control.ForEachByPrio(func(i int, manager *ctrl.Control[*zero.Ctx]) bool {
			manager.Disable(d.GroupID)
			return true
		})
	}
	context.JSON(http.StatusOK, gin.H{
		"code":    0,
		"result":  nil,
		"message": "",
		"type":    "ok",
	})
}

// HandleRequest 处理一个请求
// @Description 处理一个请求
// @Router /api/handleRequest [post]
// @Param flag formData string true "事件id" default(abc)
// @Param reason formData string false "原因" default(abc)
// @Param approve formData bool false "是否同意" default(true)
func HandleRequest(context *gin.Context) {
	var (
		d        types.HandleRequestParams
		typeName string
	)
	err := bind(&d, context)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"result":  nil,
			"message": err.Error(),
			"type":    "error",
		})
		return
	}
	r, ok := requestData.LoadAndDelete(d.Flag)
	if !ok {
		context.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"result":  nil,
			"message": "flag not found",
			"type":    "error",
		})
		return
	}
	bot := zero.GetBot(r.SelfID)
	if r.RequestType == "friend" {
		bot.SetFriendAddRequest(r.Flag, d.Approve, "")
		typeName = "好友添加"
	} else {
		bot.SetGroupAddRequest(r.Flag, r.SubType, d.Approve, d.Reason)
		if r.SubType == "add" {
			typeName = "加群请求"
		} else {
			typeName = "群邀请"
		}
	}
	log.Debugln("[gui] 已处理", r.UserID, "的"+typeName)
	context.JSON(http.StatusOK, gin.H{
		"code":    0,
		"result":  nil,
		"message": "",
		"type":    "ok",
	})
}

// GetRequestList 获取所有的请求
// @Description 获取所有的请求
// @Router /api/getRequestList [get]
// @Param selfId query integer false "机器人qq号" default(123456)
func GetRequestList(context *gin.Context) {
	d, bot, err := getBot(context)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"result":  nil,
			"message": err.Error(),
			"type":    "error",
		})
		return
	}
	var data []types.RequestVo
	requestData.Range(func(key string, value *zero.Event) bool {
		if d.SelfID != 0 && value.SelfID != d.SelfID {
			return true
		}
		data = append(data, types.RequestVo{
			Flag:        value.Flag,
			RequestType: value.RequestType,
			SubType:     value.SubType,
			Comment:     value.Comment,
			GroupID:     value.GroupID,
			GroupName:   bot.GetGroupInfo(value.GroupID, false).Name,
			UserID:      value.UserID,
			Nickname:    bot.GetStrangerInfo(value.UserID, false).Get("nickname").String(),
			SelfID:      value.SelfID,
		})
		return true
	})
	context.JSON(http.StatusOK, gin.H{
		"code":    0,
		"result":  data,
		"message": "",
		"type":    "ok",
	})
}

// GetLog 连接日志
func GetLog(context *gin.Context) {
	conn, err := upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		return
	}
	LogConn = conn
}

// Upgrade 连接ws，向前端推送message
func Upgrade(context *gin.Context) {
	conn, err := upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		return
	}
	MsgConn = conn
}

// SendMsg 前端调用发送信息
// @Description 前端调用发送信息
// @Router /api/sendMsg [post]
// @Param object body types.SendMsgParams false "发消息参数"
func SendMsg(context *gin.Context) {
	var (
		d     types.SendMsgParams
		msgID int64
	)
	err := bind(&d, context)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"result":  nil,
			"message": err.Error(),
			"type":    "error",
		})
		return
	}
	bot := zero.GetBot(d.SelfID)
	for _, gid := range d.GIDList {
		if gid > 0 {
			msgID = bot.SendGroupMessage(gid, message.UnescapeCQCodeText(d.Message))
		} else if gid < 0 {
			msgID = bot.SendPrivateMessage(-gid, message.UnescapeCQCodeText(d.Message))
		}
	}
	context.JSON(http.StatusOK, gin.H{
		"code":    0,
		"result":  msgID,
		"message": "",
		"type":    "ok",
	})
}

// Write 写入日志
func (l logWriter) Write(p []byte) (n int, err error) {
	if LogConn != nil {
		err := LogConn.WriteMessage(websocket.TextMessage, p)
		if err != nil {
			return len(p), nil
		}
	}
	return len(p), nil
}

// Login 登录接口
// @Description 前端登录
// @Router /api/login [post]
// @Param username formData string true "用户名" default(xiaoguofan)
// @Param password formData string true "密码" default(123456)
func Login(context *gin.Context) {
	var d types.LoginParams
	err := bind(&d, context)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"result":  nil,
			"message": err.Error(),
			"type":    "error",
		})
		return
	}
	user, err := control.FindUser(d.Username, d.Password)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"result":  nil,
			"message": err.Error(),
			"type":    "error",
		})
		return
	}
	token := uuid.NewString()
	middleware.LoginCache.Set(token, user, cache.DefaultExpiration)
	r := types.LoginResultVo{
		Desc:     "manager",
		RealName: user.Username,
		Roles: []types.RoleInfo{
			types.RoleInfo{
				RoleName: "Super Admin",
				Value:    "super",
			},
		},
		Token:    token,
		UserID:   int(user.ID),
		Username: user.Username,
	}
	context.JSON(http.StatusOK, gin.H{
		"code":    0,
		"result":  r,
		"message": "",
		"type":    "ok",
	})
}

// GetUserInfo 获得用户信息
// @Description 获得用户信息
// @Router /api/getUserInfo [get]
func GetUserInfo(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")
	i, _ := middleware.LoginCache.Get(token)
	user := i.(control.User)
	var qq int64
	if zero.BotConfig.SuperUsers != nil && len(zero.BotConfig.SuperUsers) > 0 {
		qq = zero.BotConfig.SuperUsers[0]
	}
	r := types.UserInfoVo{
		Desc:     "manager",
		RealName: user.Username,
		Roles: []types.RoleInfo{
			types.RoleInfo{
				RoleName: "Super Admin",
				Value:    "super",
			},
		},
		UserID:   int(user.ID),
		Username: user.Username,
		Avatar:   "https://q1.qlogo.cn/g?b=qq&nk=" + strconv.FormatInt(qq, 10) + "&s=640",
	}
	context.JSON(http.StatusOK, gin.H{
		"code":    0,
		"result":  r,
		"message": "",
		"type":    "ok",
	})
}

// Logout 登出
// @Description 登出
// @Router /api/logout [get]
func Logout(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")
	middleware.LoginCache.Delete(token)
	context.JSON(http.StatusOK, gin.H{
		"code":    0,
		"result":  nil,
		"message": "",
		"type":    "ok",
	})
}

// GetPermCode 授权码
// @Description 授权码
// @Router /api/getPermCode [get]
func GetPermCode(context *gin.Context) {
	r := []string{"1000", "3000", "5000"}
	// 先写死接口
	context.JSON(http.StatusOK, gin.H{
		"code":    0,
		"result":  r,
		"message": "",
		"type":    "ok",
	})
}

func getBot(context *gin.Context) (d types.BotParams, bot *zero.Ctx, err error) {
	err = bind(&d, context)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"result":  nil,
			"message": err.Error(),
			"type":    "error",
		})
		return
	}
	if d.SelfID == 0 {
		zero.RangeBot(func(id int64, ctx *zero.Ctx) bool {
			bot = ctx
			return false
		})
	} else {
		bot = zero.GetBot(d.SelfID)
	}
	return
}

// bind 绑定结构体, 并校验
func bind(obj any, ctx *gin.Context) (err error) {
	err = ctx.ShouldBind(obj)
	if err != nil {
		return
	}
	return validate.Struct(obj)
}
