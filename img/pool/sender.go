package pool

import (
	"errors"
	"io"
	"os"

	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/file"
	"github.com/FloatTech/zbputils/web"
	"github.com/sirupsen/logrus"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// SendImageFromPool ...
func SendImageFromPool(fallbackfile string, generatefallback func(io.Writer) error, send ctxext.NoCtxSendMsg, get ctxext.NoCtxGetMsg) error {
	m, err := GetImage(fallbackfile)
	if err != nil {
		logrus.Debugln("[ctxext.img]", err)
		if file.IsNotExist(fallbackfile) {
			f, err := os.Create(fallbackfile)
			if err != nil {
				return err
			}
			err = generatefallback(f)
			_ = f.Close()
			if err != nil {
				return err
			}
		}
		m.SetFile(file.BOTPATH + "/" + fallbackfile)
		hassent, err := m.Push(send, get)
		if hassent {
			return nil
		}
		if err != nil {
			return err
		}
	}
	// 发送图片
	img := message.Image(m.String())
	id := send(message.Message{img})
	if id == 0 {
		id = send(message.Message{img.Add("cache", "0")})
		if id == 0 {
			data, err := web.GetData(m.String())
			if err != nil {
				return err
			}
			id = send(message.Message{message.ImageBytes(data)})
			if id == 0 {
				return errors.New("图片发送失败，可能被风控了~")
			}
		}
	}
	return nil
}
