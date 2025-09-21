// Package vevent 虚拟事件
package vevent

import (
	"context"
	"unsafe"

	zero "github.com/wdvxdr1123/ZeroBot"
)

// APICallerReturnHook is a caller middleware
type APICallerReturnHook struct {
	caller   zero.APICaller
	callback func(rsp zero.APIResponse, err error)
}

// NewAPICallerReturnHook hook ctx's caller
func NewAPICallerReturnHook(ctx *zero.Ctx, callback func(rsp zero.APIResponse, err error)) (v *APICallerReturnHook) {
	return &APICallerReturnHook{
		caller:   (*(**Ctx)(unsafe.Pointer(&ctx))).caller,
		callback: callback,
	}
}

// CallAPI call original caller and pass rsp to callback
//
//nolint:revive
func (v *APICallerReturnHook) CallAPI(c context.Context, request zero.APIRequest) (rsp zero.APIResponse, err error) {
	rsp, err = v.caller.CallAPI(c, request)
	go v.callback(rsp, err)
	return
}
