package pool

import (
	"errors"

	"github.com/wdvxdr1123/ZeroBot/message"

	"github.com/FloatTech/floatbox/file"

	"github.com/FloatTech/zbputils/ctxext"
)

// SendImageFromPool ...
func SendImageFromPool(imgname, imgpath string, genimg func() error, send ctxext.NoCtxSendMsg, get ctxext.NoCtxGetMsg) error {
	if file.IsNotExist(imgpath) {
		err := genimg()
		if err != nil {
			return err
		}
	}
	// 发送图片
	img := message.Image(file.BOTPATH + "/" + imgpath)
	id := send(message.Message{img})
	if id == 0 {
		id = send(message.Message{img.Add("cache", "0")})
		if id == 0 {
			return errors.New("图片发送失败, 可能被风控了~")
		}
	}
	return nil
}
