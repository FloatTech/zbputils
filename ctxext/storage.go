package ctxext

import (
	"errors"
	"math/bits"
	"strconv"
	"strings"

	ctrl "github.com/FloatTech/zbpctrl"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	// ErrManagerNotFound unexpected
	ErrManagerNotFound = errors.New("manager not found")
)

// Storage is a wrapper of *ctrl.Control[*zero.Ctx].Get/SetData
type Storage int64

// NewStorage private is -uid
func NewStorage(ctx *zero.Ctx, gid int64) (Storage, error) {
	c, ok := ctx.State["manager"].(*ctrl.Control[*zero.Ctx])
	if !ok {
		return 0, ErrManagerNotFound
	}
	x := c.GetData(gid)
	return Storage(x), nil
}

// SaveTo ...
func (s Storage) SaveTo(ctx *zero.Ctx, gid int64) error {
	c, ok := ctx.State["manager"].(*ctrl.Control[*zero.Ctx])
	if !ok {
		return ErrManagerNotFound
	}
	return c.SetData(int64(s), gid)
}

// Get bmp is a continuous 1's sequence like 0x00ff00
func (s Storage) Get(bmp int64) int64 {
	sft := bits.TrailingZeros64(uint64(bmp))
	return (int64(s) & bmp) >> int64(sft)
}

// Set bmp is a continuous 1's sequence like 0x00ff00
func (s Storage) Set(x int64, bmp int64) Storage {
	if bmp == 0 {
		panic("cannot use bmp == 0")
	}
	sft := bits.TrailingZeros64(uint64(bmp))
	return Storage((int64(s) & (^bmp)) | ((x & (bmp >> int64(sft))) << int64(sft)))
}

// GetBool bmp must be 1-bit-long
func (s Storage) GetBool(bmp int64) bool {
	return s.Get(bmp) != 0
}

// NewStorageSaveBitmapHandler is a easy handler that cooperate with OnPrefix
// and extract int64 from args, save it to storage.
func NewStorageSaveBitmapHandler(bmp int64, minv, maxv int64) func(ctx *zero.Ctx) {
	return func(ctx *zero.Ctx) {
		args := strings.TrimSpace(ctx.State["args"].(string))
		if args == "" {
			ctx.SendChain(message.Text("ERROR: empty args"))
			return
		}
		r, err := strconv.ParseInt(args, 10, 64)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: parse int64 err: ", err))
			return
		}
		if r > maxv {
			r = maxv
		} else if r < minv {
			r = minv
		}
		gid := ctx.Event.GroupID
		if gid == 0 {
			gid = -ctx.Event.UserID
		}
		stor, err := NewStorage(ctx, gid)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		err = stor.Set(r, bmp).SaveTo(ctx, gid)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: set data err: ", err))
			return
		}
		ctx.SendChain(message.Text("成功"))
	}
}

// NewStorageSaveBoolHandler is a easy handler that cooperate with OnRegex
// and extract bool from args, save it to storage.
//
// Note:
//   - The first submatch must be (不)?
//   - If 不 is matched, true (1) will be set.
//   - Otherwise, false (0) will be set.
func NewStorageSaveBoolHandler(bmp int64) func(ctx *zero.Ctx) {
	if bits.OnesCount64(uint64(bmp)) != 1 {
		panic("bool bmp must be 1-bit-long")
	}
	return func(ctx *zero.Ctx) {
		args := ctx.State["regex_matched"].([]string)
		isone := args[1] == "不"
		gid := ctx.Event.GroupID
		if gid == 0 {
			gid = -ctx.Event.UserID
		}
		stor, err := NewStorage(ctx, gid)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		v := 0
		if isone {
			v = 1
		}
		err = stor.Set(int64(v), bmp).SaveTo(ctx, gid)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: set data err: ", err))
			return
		}
		ctx.SendChain(message.Text("成功"))
	}
}
