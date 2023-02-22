// Package router 路由
package router

import (
	"github.com/FloatTech/zbputils/control/web/controller"
	_ "github.com/FloatTech/zbputils/control/web/docs"
	"github.com/FloatTech/zbputils/control/web/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetRouters 创建路由
func SetRouters(engine *gin.Engine) {
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 支持跨域
	engine.Use(middleware.Cors(), gin.Logger())

	apiRoute := engine.Group("/api")
	apiRoute.Use(middleware.TokenMiddle())
	apiRoute.GET("/getFriendList", controller.GetFriendList)
	apiRoute.GET("/getGroupList", controller.GetGroupList)
	apiRoute.GET("/getPlugin", controller.GetPlugin)
	apiRoute.GET("/getAllPlugin", controller.GetAllPlugin)
	apiRoute.POST("/updatePluginStatus", controller.UpdatePluginStatus)
	apiRoute.POST("/updateAllPluginStatus", controller.UpdateAllPluginStatus)
	apiRoute.GET("/getRequestList", controller.GetRequestList)
	apiRoute.POST("/handleRequest", controller.HandleRequest)
	apiRoute.POST("/sendMsg", controller.SendMsg)

	noverifyRoute := engine.Group("/api")
	noverifyRoute.POST("/login", controller.Login)
	noverifyRoute.GET("/getUserInfo", controller.GetUserInfo)
	noverifyRoute.GET("/logout", controller.Logout)
	noverifyRoute.GET("/getPermCode", controller.GetPermCode)
	noverifyRoute.GET("/getBotList", controller.GetBotList)
	noverifyRoute.GET("/getLog", controller.GetLog)
	noverifyRoute.GET("/data", controller.Upgrade)
}
