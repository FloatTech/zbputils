// Package webctrl 包含 webui 所需的所有内容
package webctrl

import (
	"github.com/FloatTech/zbputils/control/web/router"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// InitGui 初始化gui
func InitGui(addr string) {
	// 将日志重定向到前端hook
	// 监听后端
	go run(addr)
	// 注册消息handle
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
	// 注册主路径路由，使其跳转到主页面
	// engine.GET("/", func(context *gin.Context) {
	// 	context.Redirect(http.StatusMovedPermanently, "/dist/dist/default.html")
	// })
	log.Infoln("[gui] the webui is running on", addr)
	log.Infoln("[gui] ", "you input the `ZeroBot-Plugin.exe -g` can disable the gui")
	log.Infoln("[gui] ", "you can see api by", "http://"+addr+"/swagger/index.html")
	if err := r.Run(addr); err != nil {
		log.Debugln("[gui] ", err.Error())
	}
}
