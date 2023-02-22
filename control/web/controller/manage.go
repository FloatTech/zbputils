package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/FloatTech/floatbox/binary"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/control/web/common"
	"github.com/FloatTech/zbputils/control/web/types"
	"github.com/FloatTech/zbputils/control/web/utils"
	"github.com/RomiChan/websocket"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	// 向前端推送消息的ws链接
	conn *websocket.Conn
	// 向前端推送日志的ws链接
	logConn *websocket.Conn

	l logWriter
	// 存储请求事件，flag作为键，一个request对象作为值
	requestData sync.Map
	upGrader    = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func init() {
	// 日志设置
	writer := io.MultiWriter(l, os.Stdout)
	log.SetOutput(writer)
	log.SetFormatter(&log.TextFormatter{DisableColors: true})

	zero.OnMessage().SetBlock(false).FirstPriority().Handle(func(ctx *zero.Ctx) {
		if conn != nil {
			err := conn.WriteJSON(ctx.Event)
			if err != nil {
				log.Errorln("[gui] 推送消息发送错误")
				return
			}
		}
	})
	// 直接注册一个request请求监听器，优先级设置为最高，设置不阻断事件传播
	zero.OnRequest(func(ctx *zero.Ctx) bool {
		if ctx.Event.RequestType == "friend" {
			ctx.State["type_name"] = "好友添加"
		} else {
			if ctx.Event.SubType == "add" {
				ctx.State["type_name"] = "加群请求"
			} else {
				ctx.State["type_name"] = "群邀请"
			}
		}
		return true
	}).SetBlock(false).FirstPriority().Handle(func(ctx *zero.Ctx) {
		r := &request{
			RequestType: ctx.Event.RequestType,
			SubType:     ctx.Event.SubType,
			Type:        ctx.State["type_name"].(string),
			GroupID:     ctx.Event.GroupID,
			UserID:      ctx.Event.UserID,
			Flag:        ctx.Event.Flag,
			Comment:     ctx.Event.Comment,
			SelfID:      ctx.Event.SelfID,
		}
		requestData.Store(ctx.Event.Flag, r)
	})
}

// logWriter
type logWriter struct {
}

// request
type request struct {
	RequestType string `json:"request_type"`
	SubType     string `json:"sub_type"`
	Type        string `json:"type"`
	Comment     string `json:"comment"`
	GroupID     int64  `json:"group_id"`
	UserID      int64  `json:"user_id"`
	Flag        string `json:"flag"`
	SelfID      int64  `json:"self_id"`
}

// GetBotList
// @Description 获取机器人qq号
// @Router /api/getBotList [get]
func GetBotList(context *gin.Context) {
	var bots []int64

	zero.RangeBot(func(id int64, ctx *zero.Ctx) bool {
		bots = append(bots, id)
		return true
	})
	common.OkWithData(bots, context)
}

// GetFriendList
// @Description 获取好友列表
// @Router /api/getFriendList [get]
// @Param selfId query integer true "机器人qq号" default(123456)
func GetFriendList(context *gin.Context) {
	var d types.BotParams
	err := common.Bind(&d, context)
	if err != nil {
		common.FailWithMessage(err.Error(), context)
		return
	}
	bot := zero.GetBot(d.SelfID)
	var resp []any
	list := bot.GetFriendList().String()
	err = json.Unmarshal(binary.StringToBytes(list), &resp)
	if err != nil {
		common.FailWithMessage(err.Error(), context)
		return
	}
	common.OkWithData(resp, context)
}

// GetGroupList
// @Description 获取群列表
// @Router /api/getGroupList [get]
// @Param selfId query integer true "机器人qq号" default(123456)
func GetGroupList(context *gin.Context) {
	var d types.BotParams
	err := common.Bind(&d, context)
	if err != nil {
		common.FailWithMessage(err.Error(), context)
		return
	}
	bot := zero.GetBot(d.SelfID)
	var resp []any
	list := bot.GetGroupList().String()
	err = json.Unmarshal(binary.StringToBytes(list), &resp)
	if err != nil {
		common.FailWithMessage(err.Error(), context)
		return
	}
	common.OkWithData(resp, context)
}

// GetAllPlugin
// @Description 获取所有插件的状态
// @Router /api/getAllPlugin [get]
// @Param groupId query integer false "群号" default(0)
func GetAllPlugin(context *gin.Context) {
	var d types.AllPluginParams
	err := common.Bind(&d, context)
	if err != nil {
		common.FailWithMessage(err.Error(), context)
		return
	}
	var datas []map[string]any
	control.ForEachByPrio(func(i int, manager *ctrl.Control[*zero.Ctx]) bool {
		datas = append(datas, map[string]any{"id": i, "name": manager.Service, "brief": manager.Options.Brief, "usage": manager.String(), "banner": "https://gitcode.net/u011570312/zbpbanner/-/raw/main/" + manager.Service + ".png", "status": manager.IsEnabledIn(d.GroupID)})
		return true
	})
	common.OkWithData(datas, context)
}

// GetPlugin
// @Description 获取某个插件的状态
// @Router /api/getPlugin [get]
// @Param groupId query integer false "群号" default(0)
// @Param name query string false "插件名" default(antibuse)
func GetPlugin(context *gin.Context) {
	var d types.PluginParams
	err := common.Bind(&d, context)
	if err != nil {
		common.FailWithMessage(err.Error(), context)
		return
	}
	con, b := control.Lookup(d.Name)
	if !b {
		common.FailWithMessage(d.Name+"服务不存在", context)
		return
	}
	data := map[string]any{"name": con.Service, "brief": con.Options.Brief, "usage": con.String(), "banner": "https://gitcode.net/u011570312/zbpbanner/-/raw/main/" + con.Service + ".png", "status": con.IsEnabledIn(d.GroupID)}
	common.OkWithData(data, context)
}

// UpdatePluginStatus
// @Description 更改某一个插件状态
// @Router /api/updatePluginStatus [post]
// @Param groupId formData integer false "群号" default(0)
// @Param name formData string true "插件名" default(aireply)
// @Param status formData boolean false "插件状态" default(true)
func UpdatePluginStatus(context *gin.Context) {
	var d types.PluginStatusParams
	err := common.Bind(&d, context)
	if err != nil {
		common.FailWithMessage(err.Error(), context)
		return
	}
	con, b := control.Lookup(d.Name)
	if !b {
		common.FailWithMessage(d.Name+"服务不存在", context)
		return
	}
	if d.Status {
		con.Enable(d.GroupID)
	} else {
		con.Disable(d.GroupID)
	}
	common.Ok(context)
}

// UpdateAllPluginStatus
// @Description 更改某群所有插件状态
// @Router /api/updateAllPluginStatus [post]
// @Param groupId formData integer false "群号" default(0)
// @Param status formData boolean false "插件状态" default(true)
func UpdateAllPluginStatus(context *gin.Context) {
	var d types.AllPluginStatusParams
	err := common.Bind(&d, context)
	if err != nil {
		common.FailWithMessage(err.Error(), context)
		return
	}
	if d.Status {
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
	common.Ok(context)
}

// HandleRequest 处理一个请求
// @Description 处理一个请求
// @Router /api/handleRequest [post]
// @Param flag formData string true "事件id" default(abc)
// @Param reason formData string false "原因" default(abc)
// @Param approve formData bool false "是否同意" default(true)
func HandleRequest(context *gin.Context) {
	var d types.HandleRequestParams
	err := common.Bind(&d, context)
	if err != nil {
		common.FailWithMessage(err.Error(), context)
		return
	}
	r, ok := requestData.LoadAndDelete(d.Flag)
	if !ok {
		common.FailWithMessage("flag not found", context)
		return
	}
	r2 := r.(*request)
	r2.handle(d.Approve, d.Reason)
	common.Ok(context)
}

// GetRequestList
// @Description 获取所有的请求
// @Router /api/getRequestList [get]
func GetRequestList(context *gin.Context) {
	var data []interface{}
	requestData.Range(func(key, value interface{}) bool {
		data = append(data, value)
		return true
	})
	common.OkWithData(data, context)
}

// GetLog 连接日志
func GetLog(context *gin.Context) {
	con1, err := upGrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		return
	}
	logConn = con1
}

// Upgrade 连接ws，向前端推送message
func Upgrade(context *gin.Context) {
	con, err := upGrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		return
	}
	conn = con
}

// SendMsg 前端调用发送信息
// @Description 前端调用发送信息
// @Router /api/sendMsg [post]
// @Param selfId formData int64 true "bot的QQ号" default(abc)
// @Param id formData int64 true "群号" default(123456)
// @Param message formData string true "消息文本" default(HelloWorld)
// @Param message_type formData string false "消息类型" default(group)
func SendMsg(context *gin.Context) {
	var d types.SendMsgParams
	err := common.Bind(&d, context)
	if err != nil {
		common.FailWithMessage(err.Error(), context)
		return
	}
	bot := zero.GetBot(d.SelfID)
	var msgID int64
	if d.MessageType == "group" {
		msgID = bot.SendGroupMessage(d.ID, message.ParseMessageFromString(d.Message))
	} else {
		msgID = bot.SendPrivateMessage(d.ID, message.ParseMessageFromString(d.Message))
	}
	common.OkWithData(msgID, context)
}

// handle 提交一个请求
func (r *request) handle(approve bool, reason string) {
	bot := zero.GetBot(r.SelfID)
	if r.RequestType == "friend" {
		bot.SetFriendAddRequest(r.Flag, approve, "")
	} else {
		bot.SetGroupAddRequest(r.Flag, r.SubType, approve, reason)
	}
	log.Debugln("[gui] ", "已处理", r.UserID, "的"+r.Type)
}

// Write 写入日志
func (l logWriter) Write(p []byte) (n int, err error) {
	if logConn != nil {
		err := logConn.WriteMessage(websocket.TextMessage, p)
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
	err := common.Bind(&d, context)
	if err != nil {
		common.FailWithMessage(err.Error(), context)
		return
	}
	user, err := control.FindUser(d.Username, d.Password)
	if err != nil {
		common.FailWithMessage(err.Error(), context)
		return
	}
	token := uuid.NewV4().String()
	utils.LoginCache.Set(token, user, cache.DefaultExpiration)
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
	common.OkWithData(r, context)
}

// GetUserInfo 获得用户信息
// @Description 获得用户信息
// @Router /api/getUserInfo [get]
func GetUserInfo(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")
	i, _ := utils.LoginCache.Get(token)
	user := i.(ctrl.User)
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
		Avatar:   "https://q1.qlogo.cn/g?b=qq&nk=1156544355&s=640",
	}
	common.OkWithData(r, context)
}

// Logout 登出
// @Description 登出
// @Router /api/logout [get]
func Logout(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")
	utils.LoginCache.Delete(token)
	common.Ok(context)
}

// GetPermCode 授权码
// @Description 授权码
// @Router /api/getPermCode [get]
func GetPermCode(context *gin.Context) {
	r := []string{"1000", "3000", "5000"}
	// 先写死接口
	common.OkWithData(r, context)
}
