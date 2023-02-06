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
	"github.com/RomiChan/websocket"
	"github.com/gin-gonic/gin"
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
	writer := io.MultiWriter(l, os.Stdout)
	log.SetOutput(writer)

	zero.OnMessage().SetBlock(false).FirstPriority().Handle(func(ctx *zero.Ctx) {
		if conn != nil {
			err := conn.WriteJSON(ctx.Event)
			if err != nil {
				log.Debugln("[gui] " + "向发送错误")
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
// @Description:
type logWriter struct {
}

// request
// @Description: 一个请求事件的结构体
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

// BotRequest GetGroupList,GetFriendList的入参
type BotRequest struct {
	SelfID int64 `json:"self_id" form:"self_id" validate:"required"`
}

// AllPluginDto GetAllPlugin的入参
type AllPluginDto struct {
	GroupID int64 `json:"group_id" form:"group_id"`
}

// PluginDto GetPlugin的入参
type PluginDto struct {
	GroupID int64  `json:"group_id" form:"group_id"`
	Name    string `json:"name" form:"name"`
}

// PluginStatusDto UpdatePluginStatus的入参
type PluginStatusDto struct {
	GroupID int64  `json:"group_id" form:"group_id"`
	Name    string `json:"name" form:"name" validate:"required"`
	Status  bool   `json:"status" form:"status"`
}

// AllPluginStatusDto UpdateAllPluginStatus的入参
type AllPluginStatusDto struct {
	GroupID int64 `json:"group_id" form:"group_id"`
	Status  bool  `json:"status" form:"status"`
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
// @Param self_id query integer true "机器人qq号" default(123456)
func GetFriendList(context *gin.Context) {
	var d BotRequest
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
// @Param self_id query integer true "机器人qq号" default(123456)
func GetGroupList(context *gin.Context) {
	var d BotRequest
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
// @Param group_id query integer false "群号" default(0)
func GetAllPlugin(context *gin.Context) {
	var d AllPluginDto
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
// @Param group_id query integer false "群号" default(0)
// @Param name query string false "插件名" default(antibuse)
func GetPlugin(context *gin.Context) {
	var d PluginDto
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
// @Param group_id formData integer false "群号" default(0)
// @Param name formData string true "插件名" default(aireply)
// @Param status formData boolean false "插件状态" default(true)
func UpdatePluginStatus(context *gin.Context) {
	var d PluginStatusDto
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
// @Param group_id formData integer false "群号" default(0)
// @Param status formData boolean false "插件状态" default(true)
func UpdateAllPluginStatus(context *gin.Context) {
	var d AllPluginStatusDto
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

// HandelRequest 处理一个请求
func HandelRequest(context *gin.Context) {
	var data map[string]interface{}
	err := context.BindJSON(&data)
	if err != nil {
		context.JSON(404, nil)
		return
	}
	r, ok := requestData.LoadAndDelete(data["flag"].(string))
	if !ok {
		context.JSON(404, "flag not found")
	}
	r2 := r.(*request)
	r2.handle(data["approve"].(bool), data["reason"].(string))
	context.JSON(200, "操作成功")
}

// GetRequests 获取所有的请求
func GetRequests(context *gin.Context) {
	var data []interface{}
	requestData.Range(func(key, value interface{}) bool {
		data = append(data, value)
		return true
	})
	context.JSON(200, data)
}

// GetLogs 连接日志
func GetLogs(context *gin.Context) {
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
func SendMsg(context *gin.Context) {
	var data map[string]interface{}
	err := context.BindJSON(&data)
	if err != nil {
		context.JSON(404, nil)
		return
	}
	selfID := int64(data["self_id"].(float64))
	id := int64(data["id"].(float64))
	message1 := data["message"].(string)
	messageType := data["message_type"].(string)

	bot := zero.GetBot(selfID)
	var msgID int64
	if messageType == "group" {
		msgID = bot.SendGroupMessage(id, message.ParseMessageFromString(message1))
	} else {
		msgID = bot.SendPrivateMessage(id, message.ParseMessageFromString(message1))
	}
	context.JSON(200, msgID)
}

// Handle 提交一个请求
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
