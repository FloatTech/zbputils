// Package webctrl 包含 webui 所需的所有内容
package webctrl

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/FloatTech/zbputils/control/web/router"
	"github.com/gin-gonic/gin"
	webui "github.com/guohuiyuan/ZeroBot-Plugin-Webui"
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

	engine := gin.Default()
	router.SetRouters(engine)

	staticEngine := gin.Default()
	df, _ := fs.Sub(webui.Dist, "dist")
	staticEngine.StaticFS("/", http.FS(df))

	log.Infoln("[gui] the webui is running on", "http://"+addr)
	log.Infoln("[gui] ", "you input the `ZeroBot-Plugin.exe -g` can disable the gui")
	log.Infoln("[gui] ", "you can see api by", "http://"+addr+"/swagger/index.html")
	server := &http.Server{
		Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			// 如果 URL 以 /api, /swagger 开头, 走后端路由
			if strings.HasPrefix(request.URL.Path, "/api") || strings.HasPrefix(request.URL.Path, "/swagger") {
				engine.ServeHTTP(writer, request)
				return
			}
			// 否则，走前端路由
			staticEngine.ServeHTTP(writer, request)
		}),
		Addr: addr,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Debugln("[gui] ", err.Error())
	}
}
