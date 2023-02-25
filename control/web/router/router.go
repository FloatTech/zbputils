// Package router 路由
package router

import (
	"github.com/FloatTech/zbputils/control/web/controller"
	_ "github.com/FloatTech/zbputils/control/web/docs" // swagger数据
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
	apiRoute.GET("/getRequestList", controller.GetRequestList)
	apiRoute.POST("/handleRequest", controller.HandleRequest)
	apiRoute.POST("/sendMsg", controller.SendMsg)
	apiRoute.GET("/getUserInfo", controller.GetUserInfo)

	manageRoute := apiRoute.Group("/manage")
	manageRoute.GET("/getPlugin", controller.GetPlugin)
	manageRoute.GET("/getAllPlugin", controller.GetAllPlugin)
	manageRoute.POST("/updatePluginStatus", controller.UpdatePluginStatus)
	manageRoute.POST("/updateResponseStatus", controller.UpdateResponseStatus)
	manageRoute.POST("/updateAllPluginStatus", controller.UpdateAllPluginStatus)

	noVerifyRoute := engine.Group("/api")
	noVerifyRoute.POST("/login", controller.Login)
	noVerifyRoute.GET("/logout", controller.Logout)
	noVerifyRoute.GET("/getPermCode", controller.GetPermCode)
	noVerifyRoute.GET("/getBotList", controller.GetBotList)
	noVerifyRoute.GET("/getLog", controller.GetLog)
	noVerifyRoute.GET("/data", controller.Upgrade)
}
