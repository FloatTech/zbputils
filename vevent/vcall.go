package vevent

import (
	"unsafe"

	zero "github.com/wdvxdr1123/ZeroBot"
)

// APICallerHook is a caller middleware
type APICallerHook struct {
	caller   zero.APICaller
	callback func(rsp zero.APIResponse, err error)
}

// NewAPICallerHook hook ctx's caller
func NewAPICallerHook(ctx *zero.Ctx, callback func(rsp zero.APIResponse, err error)) (v *APICallerHook) {
	return &APICallerHook{
		caller:   (*(**Ctx)(unsafe.Pointer(&ctx))).caller,
		callback: callback,
	}
}

// CallApi call original caller and pass rsp to callback
func (v *APICallerHook) CallApi(request zero.APIRequest) (rsp zero.APIResponse, err error) {
	rsp, err = v.caller.CallApi(request)
	go v.callback(rsp, err)
	return
}
