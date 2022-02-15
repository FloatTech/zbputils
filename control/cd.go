package control

import (
	"encoding/binary"
	"strings"
	"time"

	b14 "github.com/fumiama/go-base16384"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	binutils "github.com/FloatTech/zbputils/binary"
	"github.com/FloatTech/zbputils/math"
	"github.com/FloatTech/zbputils/process"
)

var startTime int64

func init() {
	// 插件冲突检测 会在本群发送一条消息并在约 1s 后撤回
	zero.OnFullMatch("插件冲突检测", zero.OnlyGroup, zero.AdminPermission, zero.OnlyToMe).SetBlock(true).FirstPriority().
		Handle(func(ctx *zero.Ctx) {
			tok, err := genToken()
			if err != nil {
				return
			}
			t := message.Text("●cd" + tok)
			startTime = time.Now().Unix()
			id := ctx.SendChain(t)
			process.SleepAbout1sTo2s()
			ctx.DeleteMessage(id)
		})

	zero.OnRegex("^●cd([\u4e00-\u8e00]{4})$", zero.OnlyGroup).SetBlock(true).FirstPriority().
		Handle(func(ctx *zero.Ctx) {
			if isValidToken(ctx.State["regex_matched"].([]string)[1]) {
				gid := ctx.Event.GroupID
				w := binutils.SelectWriter()
				ForEach(func(key string, manager *Control) bool {
					if manager.IsEnabledIn(gid) {
						w.WriteString("\xfe\xff")
						w.WriteString(key)
					}
					return true
				})
				if w.Len() > 2 {
					my, err := b14.UTF16be2utf8(b14.Encode(w.Bytes()[2:]))
					binutils.PutWriter(w)
					if err == nil {
						my, cl := binutils.OpenWriterF(func(w *binutils.Writer) {
							w.WriteString("●cd●")
							w.Write(my)
						})
						id := ctx.SendChain(message.Text(binutils.BytesToString(my)))
						cl()
						process.SleepAbout1sTo2s()
						ctx.DeleteMessage(id)
					}
				}
			}
		})

	zero.OnRegex("^●cd●(([\u4e00-\u8e00]*[\u3d01-\u3d06]?))", zero.OnlyGroup).SetBlock(true).FirstPriority().
		Handle(func(ctx *zero.Ctx) {
			if time.Now().Unix()-startTime < 10 {
				msg, err := b14.UTF82utf16be(binutils.StringToBytes(ctx.State["regex_matched"].([]string)[1]))
				if err == nil {
					gid := ctx.Event.GroupID
					for _, s := range strings.Split(b14.DecodeString(msg), "\xfe\xff") {
						manmu.RLock()
						c, ok := managers[s]
						manmu.RUnlock()
						if ok && c.IsEnabledIn(gid) {
							c.Disable(gid)
						}
					}
				}
			}
		})
}

func genToken() (tok string, err error) {
	timebytes, cl := binutils.OpenWriterF(func(w *binutils.Writer) {
		w.WriteUInt64(uint64(time.Now().Unix()))
	})
	timebytes, err = b14.UTF16be2utf8(b14.Encode(timebytes[1:]))
	cl()
	if err == nil {
		tok = binutils.BytesToString(timebytes)
	}
	return
}

func isValidToken(tok string) (yes bool) {
	s, err := b14.UTF82utf16be(binutils.StringToBytes(tok))
	if err == nil {
		timebytes, cl := binutils.OpenWriterF(func(w *binutils.Writer) {
			w.WriteByte(0)
			w.Write(b14.Decode(s))
		})
		yes = math.Abs64(time.Now().Unix()-int64(binary.BigEndian.Uint64(timebytes))) < 10
		cl()
	}
	return
}
