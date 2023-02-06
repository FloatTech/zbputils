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
func SetRouters(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Use(middleware.Cors())
	r.Use(gin.Logger())
	engine := r.Group("/api")
	// 支持跨域
	engine.GET("/getBotList", controller.GetBotList)
	engine.GET("/getFriendList", controller.GetFriendList)
	engine.GET("/getGroupList", controller.GetGroupList)
	engine.GET("/getPlugin", controller.GetPlugin)
	engine.GET("/getAllPlugin", controller.GetAllPlugin)
	engine.POST("/updatePluginStatus", controller.UpdatePluginStatus)
	engine.POST("/updateAllPluginStatus", controller.UpdateAllPluginStatus)
	engine.GET("/getRequests", controller.GetRequests)
	engine.POST("/handleRequest", controller.HandleRequest)
	engine.GET("/getLog", controller.GetLog)
	engine.POST("/sendMsg", controller.SendMsg)
	engine.GET("/data", controller.Upgrade)
}
