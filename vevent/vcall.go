// Package vevent 虚拟事件
package vevent

import (
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
//nolint:stylecheck,revive
func (v *APICallerReturnHook) CallAPI(request zero.APIRequest) (rsp zero.APIResponse, err error) {
	rsp, err = v.caller.CallAPI(request)
	go v.callback(rsp, err)
	return
}
