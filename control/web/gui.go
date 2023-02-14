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
			log.Errorln("[gui] ZeroBot-Plugin-Webui出现不可恢复的错误")
			log.Errorln("[gui]", err)
		}
	}()

	engine := gin.New()
	router.SetRouters(engine)
	log.Infoln("[gui] the webui is running on", "http://"+addr)
	log.Infoln("[gui] ", "you input the `ZeroBot-Plugin.exe -g` can disable the gui")
	log.Infoln("[gui] ", "you can see api by", "http://"+addr+"/swagger/index.html")
	if err := engine.Run(addr); err != nil {
		log.Debugln("[gui] ", err.Error())
	}
}
