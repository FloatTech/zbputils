package ctxext

import (
	"hash/crc64"
	"math/rand"
	"time"
	"unsafe"

	"github.com/FloatTech/zbputils/binary"
)

func RandSenderPerDayN(uid int64, n int) int {
	sum := crc64.New(crc64.MakeTable(crc64.ISO))
	sum.Write(binary.StringToBytes(time.Now().Format("20060102")))
	sum.Write((*[8]byte)(unsafe.Pointer(&uid))[:])
	r := rand.New(rand.NewSource(int64(sum.Sum64())))
	return r.Intn(n)
}
