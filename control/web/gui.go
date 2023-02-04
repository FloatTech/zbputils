// Package webctrl 包含 webui 所需的所有内容
package webctrl

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"

	// 依赖gin监听server
	"github.com/RomiChan/websocket"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/control/web/controller"
	_ "github.com/FloatTech/zbputils/control/web/docs"
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

	magageController := &controller.ManageController{}
	r := gin.New()
	engine := r.Group("/api")
	// 支持跨域
	engine.Use(cors())
	// 注册静态文件
	// engine.StaticFS("/dist", http.FS(manager.Dist))
	engine.POST("/getBotList", magageController.GetBotList)
	engine.POST("/getGroupList", magageController.GetGroupList)
	engine.POST("/getFriendList", getFriendList)
	// 注册主路径路由，使其跳转到主页面
	// engine.GET("/", func(context *gin.Context) {
	// 	context.Redirect(http.StatusMovedPermanently, "/dist/dist/default.html")
	// })
	// 更改某个插件状态
	engine.POST("/updatePluginStatus", magageController.UpdatePluginStatus)
	// 更改某一个插件在所有群的状态
	engine.POST("/update_plugin_all_group_status", updatePluginAllGroupStatus)
	// 更改所有插件状态
	engine.POST("/update_all_plugin_status", updateAllPluginStatus)
	// 获取所有插件
	engine.POST("/getPluginList", magageController.GetPluginList)
	// 获取一个插件
	engine.POST("/getPlugin", getPlugin)
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
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
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

// updateAllPluginStatus
/**
 * @Description: 改变所有插件的状态
 * @param context
 * example
 */
func updateAllPluginStatus(context *gin.Context) {
	enable, err := strconv.ParseBool(context.PostForm("enable"))
	if err != nil {
		var parse map[string]interface{}
		err := context.BindJSON(&parse)
		if err != nil {
			log.Errorln("[gui] " + err.Error())
			return
		}
		enable = parse["enable"].(bool)
	}
	control.ForEachByPrio(func(i int, manager *ctrl.Control[*zero.Ctx]) bool {
		if enable {
			manager.Enable(0)
		} else {
			manager.Disable(0)
		}
		return true
	})
	context.JSON(200, nil)
}

// updatePluginAllGroupStatus
/**
 * @Description: 改变插件在所有群的状态
 * @param context
 * example
 */
func updatePluginAllGroupStatus(context *gin.Context) {
	name := context.PostForm("name")
	enable, err := strconv.ParseBool(context.PostForm("enable"))
	if err != nil {
		var parse map[string]interface{}
		err := context.BindJSON(&parse)
		if err != nil {
			log.Errorln("[gui]" + err.Error())
			return
		}
		name = parse["name"].(string)
		enable = parse["enable"].(bool)
	}
	con, b := control.Lookup(name)
	if !b {
		context.JSON(404, nil)
		return
	}
	if enable {
		con.Enable(0)
	} else {
		con.Disable(0)
	}
	context.JSON(200, nil)
}

// getPlugin
/**
 * @Description: 获取一个插件的状态
 * @param context
 * example
 */
func getPlugin(context *gin.Context) {
	groupID, err := strconv.ParseInt(context.PostForm("group_id"), 10, 64)
	name := context.PostForm("name")
	if err != nil {
		var parse map[string]interface{}
		err := context.BindJSON(&parse)
		if err != nil {
			log.Errorln("[gui]" + err.Error())
			return
		}
		groupID = int64(parse["group_id"].(float64))
		name = parse["name"].(string)
	}
	con, b := control.Lookup(name)
	if !b {
		context.JSON(404, "服务不存在")
		return
	}
	context.JSON(200, gin.H{"enable": con.IsEnabledIn(groupID)})
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

// getFriendList
/**
 * @Description: 获取好友列表
 * @param context
 * example
 */
func getFriendList(context *gin.Context) {
	selfID, err := strconv.Atoi(context.PostForm("self_id"))
	if err != nil {
		log.Errorln("[gui]" + err.Error())
		var data map[string]interface{}
		err := context.BindJSON(&data)
		if err != nil {
			log.Errorln("[gui]" + err.Error())
			log.Errorln("[gui]" + "绑定错误")
			return
		}
		selfID = int(data["self_id"].(float64))
	}
	bot := zero.GetBot(int64(selfID))
	var resp []interface{}
	list := bot.GetFriendList().String()
	err = json.Unmarshal([]byte(list), &resp)
	if err != nil {
		log.Errorln("[gui]" + err.Error())
		log.Errorln("[gui]" + "解析json错误")
	}
	context.JSON(200, resp)
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

// cors
/**
 * @Description: 支持跨域访问
 * @return gin.HandlerFunc
 * example
 */
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") // 请求头部
		if origin != "" {
			// 接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			// 服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			// 允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session, Content-Type")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			// 设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			// 允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// 允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()

		c.Next()
	}
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
