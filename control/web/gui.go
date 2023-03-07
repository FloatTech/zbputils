// Package webctrl 包含 webui 所需的所有内容
package webctrl

import (
	"context"
	"io/fs"
	"net/http"
	"runtime/debug"
	"strings"

	webui "github.com/FloatTech/ZeroBot-Plugin-Webui"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/control/web/controller"
	"github.com/FloatTech/zbputils/control/web/router"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// RunGui 运行webui
// @title zbp api
// @version 1.0
// @description zbp restful api document
// @host 127.0.0.1:3000
// @BasePath /
func RunGui(addr string) {
	defer func() {
		err := recover()
		if err != nil {
			log.Errorln("[gui] ZeroBot-Plugin-Webui出现不可恢复的错误")
			log.Errorln("[gui] err:", err, ",stack:", debug.Stack())
		}
	}()

	engine := gin.Default()
	router.SetRouters(engine)

	staticEngine := gin.Default()
	df, _ := fs.Sub(webui.Dist, "dist")
	staticEngine.StaticFS("/", http.FS(df))
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
	for canrun := range control.ListenCtrlChan {
		if canrun {
			if err := server.Shutdown(context.TODO()); err != nil {
				log.Errorln("[gui] server shutdown err: ", err.Error())
			}
			server = &http.Server{
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
			go func() {
				log.Infoln("[gui] the webui is running on", "http://"+addr)
				log.Infoln("[gui] you can see api by http://" + addr + "/swagger/index.html")
				if err := server.ListenAndServe(); err != nil {
					log.Errorln("[gui] server listen err: ", err.Error())
				}
			}()
		} else {
			if err := server.Shutdown(context.TODO()); err != nil {
				log.Errorln("[gui] server shutdown err: ", err.Error())
			}
			controller.MsgConn = nil
			controller.LogConn = nil
		}
	}
}
