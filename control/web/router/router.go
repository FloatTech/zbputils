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
	engine.GET("/getGroupList", controller.GetGroupList)
	engine.GET("/getPluginList", controller.GetPluginList)
	engine.POST("/updatePluginStatus", controller.UpdatePluginStatus)
}
