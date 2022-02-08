package binary

import _ "unsafe" // to use linkname

// BytesToString 没有内存开销的转换
//go:linkname BytesToString github.com/wdvxdr1123/ZeroBot/utils/helper.BytesToString
func BytesToString(b []byte) string

// StringToBytes 没有内存开销的转换
//go:linkname StringToBytes github.com/wdvxdr1123/ZeroBot/utils/helper.StringToBytes
func StringToBytes(s string) (b []byte)
