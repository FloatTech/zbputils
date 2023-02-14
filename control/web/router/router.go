// Package router 路由
package router

import (
	"net/http"

	"github.com/FloatTech/zbputils/control/web/controller"
	_ "github.com/FloatTech/zbputils/control/web/docs"
	"github.com/FloatTech/zbputils/control/web/middleware"
	"github.com/gin-gonic/gin"
	webui "github.com/guohuiyuan/ZeroBot-Plugin-Webui"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetRouters 创建路由
func SetRouters(engine *gin.Engine) {
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	engine.StaticFS("/dist", http.FS(webui.Dist))
	// 支持跨域
	engine.Use(middleware.Cors())
	engine.Use(gin.Logger())

	apiRoute := engine.Group("/api")
	apiRoute.GET("/getBotList", controller.GetBotList)
	apiRoute.GET("/getFriendList", controller.GetFriendList)
	apiRoute.GET("/getGroupList", controller.GetGroupList)
	apiRoute.GET("/getPlugin", controller.GetPlugin)
	apiRoute.GET("/getAllPlugin", controller.GetAllPlugin)
	apiRoute.POST("/updatePluginStatus", controller.UpdatePluginStatus)
	apiRoute.POST("/updateAllPluginStatus", controller.UpdateAllPluginStatus)
	apiRoute.GET("/getRequests", controller.GetRequests)
	apiRoute.POST("/handleRequest", controller.HandleRequest)
	apiRoute.GET("/getLog", controller.GetLog)
	apiRoute.POST("/sendMsg", controller.SendMsg)
	apiRoute.GET("/data", controller.Upgrade)
	apiRoute.POST("/login", controller.Login)
	apiRoute.GET("/getUserInfo", controller.GetUserInfo)
	apiRoute.GET("/logout", controller.Logout)
	apiRoute.GET("/getPermCode", controller.GetPermCode)
}
