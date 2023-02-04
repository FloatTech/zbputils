// Package controller 基本返回方法
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	// Success 成功代码
	Success = 0
	// Error 错误代码
	Error = -1
	// SuccessMessageType 成功消息类型
	SuccessMessageType = "ok"
	// ErrorMessageType 错误消息类型
	ErrorMessageType = "error"
)

// BaseController 基础控制器
type BaseController struct {
}

// Response 基础返回类
type Response struct {
	Code        int         `json:"code"`
	Result      interface{} `json:"result"`
	Message     string      `json:"message"`
	MessageType string      `json:"type"`
}

// Result 通用方法
func (c *BaseController) Result(code int, result interface{}, message string, messageType string, ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response{
		code,
		result,
		message,
		messageType,
	})
}

// Ok 成功返回
func (c *BaseController) Ok(ctx *gin.Context) {
	c.Result(Success, map[string]interface{}{}, "", SuccessMessageType, ctx)
}

// OkWithMessage 带消息成功返回
func (c *BaseController) OkWithMessage(message string, ctx *gin.Context) {
	c.Result(Success, map[string]interface{}{}, message, SuccessMessageType, ctx)
}

// OkWithData 带数据成功返回
func (c *BaseController) OkWithData(result interface{}, ctx *gin.Context) {
	c.Result(Success, result, "", SuccessMessageType, ctx)
}

// OkWithDetailed 详细返回
func (c *BaseController) OkWithDetailed(result interface{}, message string, ctx *gin.Context) {
	c.Result(Success, result, message, SuccessMessageType, ctx)
}

// Fail 失败返回
func (c *BaseController) Fail(ctx *gin.Context) {
	c.Result(Error, map[string]interface{}{}, "", ErrorMessageType, ctx)
}

// FailWithMessage 带消息失败返回
func (c *BaseController) FailWithMessage(message string, ctx *gin.Context) {
	c.Result(Error, map[string]interface{}{}, message, ErrorMessageType, ctx)
}

// FailWithDetailed 详细返回
func (c *BaseController) FailWithDetailed(result interface{}, message string, ctx *gin.Context) {
	c.Result(Error, result, message, ErrorMessageType, ctx)
}
