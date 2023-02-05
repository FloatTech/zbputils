// Package webctrl 包含 webui 所需的所有内容
package webctrl

import (
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/FloatTech/zbputils/control/web/router"
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
)

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

// InitGui 初始化gui
func InitGui(addr string) {
	// 将日志重定向到前端hook
	writer := io.MultiWriter(l, os.Stdout)
	log.SetOutput(writer)
	// 监听后端
	go run(addr)
	// 注册消息handle
	messageHandle()
}

// websocket的协议升级
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// @title zbp api
// @version 1.0
// @description zbp restful api document
// @host 127.0.0.1:3000
// @BasePath /
func run(addr string) {
	defer func() {
		err := recover()
		if err != nil {
			log.Errorln("[gui]" + "bot-manager出现不可恢复的错误")
			log.Errorln("[gui]", err)
		}
	}()

	r := gin.New()
	router.SetRouters(r)
	engine := r.Group("/api")
	// 注册主路径路由，使其跳转到主页面
	// engine.GET("/", func(context *gin.Context) {
	// 	context.Redirect(http.StatusMovedPermanently, "/dist/dist/default.html")
	// })
	// 获取所有请求
	engine.POST("/get_requests", getRequests)
	// 执行一个请求事件
	engine.POST("handle_request", handelRequest)
	// 链接日志
	engine.GET("/get_log", getLogs)
	// 获取前端标签
	engine.GET("/get_label", func(context *gin.Context) {
		context.JSON(200, "ZeroBot-Plugin")
	})
	// 发送信息
	engine.POST("/send_msg", sendMsg)
	engine.GET("/data", upgrade)
	log.Infoln("[gui] the webui is running on", addr)
	log.Infoln("[gui] ", "you input the `ZeroBot-Plugin.exe -g` can disable the gui")
	log.Infoln("[gui] ", "you can see api by", "http://"+addr+"/swagger/index.html")
	if err := r.Run(addr); err != nil {
		log.Debugln("[gui] ", err.Error())
	}
}

// handelRequest
/**
 * @Description: 处理一个请求
 * @param context
 */
func handelRequest(context *gin.Context) {
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

// getRequests
/**
 * @Description: 获取所有的请求
 * @param context
 */
func getRequests(context *gin.Context) {
	var data []interface{}
	requestData.Range(func(key, value interface{}) bool {
		data = append(data, value)
		return true
	})
	context.JSON(200, data)
}

// getLogs
/**
 * @Description: 连接日志
 * @param context
 * example
 */
func getLogs(context *gin.Context) {
	con1, err := upGrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		return
	}
	logConn = con1
}

// MessageHandle
/**
 * @Description: 定义一个向前端发送信息的handle
 * example
 */
func messageHandle() {
	defer func() {
		err := recover()
		if err != nil {
			log.Errorln("[gui]" + "bot-manager出现不可恢复的错误")
			log.Errorln("[gui] ", err)
		}
	}()

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

// upgrade
/**
 * @Description: 连接ws，向前端推送message
 * @param context
 * example
 */
func upgrade(context *gin.Context) {
	con, err := upGrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		return
	}
	conn = con
}

// sendMsg
/**
 * @Description: 前端调用发送信息
 * @param context
 * example
 */
func sendMsg(context *gin.Context) {
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

// handle
/**
 * @Description: 提交一个请求
 * @receiver r
 * @param approve 是否通过
 * @param reason 拒绝的理由
 */
func (r *request) handle(approve bool, reason string) {
	bot := zero.GetBot(r.SelfID)
	if r.RequestType == "friend" {
		bot.SetFriendAddRequest(r.Flag, approve, "")
	} else {
		bot.SetGroupAddRequest(r.Flag, r.SubType, approve, reason)
	}
	log.Debugln("[gui] ", "已处理", r.UserID, "的"+r.Type)
}

func (l logWriter) Write(p []byte) (n int, err error) {
	if logConn != nil {
		err := logConn.WriteMessage(websocket.TextMessage, p)
		if err != nil {
			return len(p), nil
		}
	}
	return len(p), nil
}
