package vevent

import (
	"unsafe"

	zero "github.com/wdvxdr1123/ZeroBot"
)

type APICallerHook struct {
	caller   zero.APICaller
	callback func(rsp zero.APIResponse, err error)
}

func NewAPICallerHook(ctx *zero.Ctx, callback func(rsp zero.APIResponse, err error)) (v *APICallerHook) {
	return &APICallerHook{
		caller:   (*(**Ctx)(unsafe.Pointer(&ctx))).caller,
		callback: callback,
	}
}

func (v *APICallerHook) CallApi(request zero.APIRequest) (rsp zero.APIResponse, err error) {
	rsp, err = v.caller.CallApi(request)
	go v.callback(rsp, err)
	return
}
