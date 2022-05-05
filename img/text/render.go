// Package text 文字转图片
package text

import (
	"image"
	"os"
	"strings"

	"github.com/fogleman/gg"
	log "github.com/sirupsen/logrus"

	"github.com/FloatTech/zbputils/file"
	"github.com/FloatTech/zbputils/img/writer"
)

// 加载数据库
func init() {
	_ = os.MkdirAll(FontPath, 0755)
}

type Text struct {
	canvas *gg.Context
}

// RenderToBase64 文字转base64
func RenderToBase64(text, font string, width, fontSize int) (base64Bytes []byte, err error) {
	txtc, err := Render(text, font, width, fontSize)
	if err != nil {
		log.Println("[txt2img]", err)
		return nil, err
	}
	base64Bytes, err = writer.ToBase64(txtc.Image())
	if err != nil {
		log.Println("[txt2img]", err)
		return nil, err
	}
	return
}

// Render 文字转图片 width 是图片宽度
func Render(text, font string, width, fontSize int) (txtc Text, err error) {
	_, err = file.GetLazyData(font, true)
	if err != nil {
		return
	}

	txtc.canvas = gg.NewContext(width, fontSize) // fake
	if err = txtc.canvas.LoadFontFace(font, float64(fontSize)); err != nil {
		log.Errorln("[txt2img]", err)
		return
	}

	buff := make([]string, 0)
	for _, s := range strings.Split(text, "\n") {
		line := ""
		for _, v := range s {
			length, _ := txtc.canvas.MeasureString(line)
			if int(length) <= width {
				line += string(v)
			} else {
				buff = append(buff, line)
				line = string(v)
			}
		}
		buff = append(buff, line)
	}

	_, h := txtc.canvas.MeasureString("好")
	txtc.canvas = gg.NewContext(width+int(h*2+0.5), int(float64(len(buff)*3+1)/2*h+0.5))
	txtc.canvas.SetRGB(1, 1, 1)
	txtc.canvas.Clear()
	txtc.canvas.SetRGB(0, 0, 0)
	if err = txtc.canvas.LoadFontFace(font, float64(fontSize)); err != nil {
		log.Errorln("[txt2img]", err)
		return
	}

	for i, v := range buff {
		if v != "" {
			txtc.canvas.DrawString(v, float64(width)*0.01, float64((i+1)*3)/2*h)
		}
	}
	return
}

// Image ...
func (txtc Text) Image() image.Image {
	return txtc.canvas.Image()
}
