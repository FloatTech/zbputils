// Package common 响应
package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

var (
	validate = validator.New()
)

// Response 基础返回类
type Response struct {
	Code        int    `json:"code"`
	Result      any    `json:"result"`
	Message     string `json:"message"`
	MessageType string `json:"type"`
}

// Result 通用方法
func Result(code int, result any, message string, messageType string, ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response{
		code,
		result,
		message,
		messageType,
	})
}

// Ok 成功返回
func Ok(ctx *gin.Context) {
	Result(Success, map[string]any{}, "", SuccessMessageType, ctx)
}

// OkWithMessage 带消息成功返回
func OkWithMessage(message string, ctx *gin.Context) {
	Result(Success, map[string]any{}, message, SuccessMessageType, ctx)
}

// OkWithData 带数据成功返回
func OkWithData(result any, ctx *gin.Context) {
	Result(Success, result, "", SuccessMessageType, ctx)
}

// OkWithDetailed 详细返回
func OkWithDetailed(result any, message string, ctx *gin.Context) {
	Result(Success, result, message, SuccessMessageType, ctx)
}

// Fail 失败返回
func Fail(ctx *gin.Context) {
	Result(Error, map[string]any{}, "", ErrorMessageType, ctx)
}

// FailWithMessage 带消息失败返回
func FailWithMessage(message string, ctx *gin.Context) {
	Result(Error, map[string]any{}, message, ErrorMessageType, ctx)
}

// FailWithDetailed 详细返回
func FailWithDetailed(result any, message string, ctx *gin.Context) {
	Result(Error, result, message, ErrorMessageType, ctx)
}

// Bind 绑定结构体, 并校验
func Bind(obj any, ctx *gin.Context) (err error) {
	err = ctx.ShouldBind(obj)
	if err != nil {
		return
	}
	return validate.Struct(obj)
}

// NotLoggedIn 未登录
func NotLoggedIn(code int, result any, message string, messageType string, ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, Response{
		code,
		result,
		message,
		messageType,
	})
}
