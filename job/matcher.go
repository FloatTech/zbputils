package job

import (
	"unsafe"

	zero "github.com/wdvxdr1123/ZeroBot"

	"github.com/FloatTech/zbputils/control"
)

func getmatcher(m control.Matcher) *zero.Matcher {
	return (*zero.Matcher)(*(*unsafe.Pointer)(unsafe.Add(unsafe.Pointer(&m), unsafe.Sizeof(uintptr(1)))))
}
