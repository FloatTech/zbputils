package job

import (
	"unsafe"

	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
)

func getmatcher(m control.Matcher) *zero.Matcher {
	return (*zero.Matcher)(*(*unsafe.Pointer)(unsafe.Add(unsafe.Pointer(&m), unsafe.Sizeof(uintptr(1)))))
}
