// Package text 文字转图片
package text

import (
	"bufio"
	"image"
	"os"
	"strings"

	"github.com/FloatTech/gg"
	log "github.com/sirupsen/logrus"

	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/img/writer"
)

// 加载数据库
func init() {
	_ = os.MkdirAll(FontPath, 0755)
}

// RenderToBase64 文字转base64
func RenderToBase64(text, font string, width, fontSize int) (base64Bytes []byte, err error) {
	im, err := Render(text, font, width, fontSize)
	if err != nil {
		log.Println("[txt2img]", err)
		return nil, err
	}
	base64Bytes, err = writer.ToBase64(im)
	if err != nil {
		log.Println("[txt2img]", err)
		return nil, err
	}
	return
}

// Render 文字转图片 width 是图片宽度
func Render(text, font string, width, fontSize int) (txtPic image.Image, err error) {
	_, err = file.GetLazyData(font, "data/control/stor.spb", true)
	if err != nil {
		return
	}

	canvas := gg.NewContext(width, fontSize) // fake
	if err = canvas.LoadFontFace(font, float64(fontSize)); err != nil {
		return
	}
	buff := make([]string, 0, 32)
	s := bufio.NewScanner(strings.NewReader(text))
	line := strings.Builder{}
	for s.Scan() {
		for _, v := range s.Text() {
			length, _ := canvas.MeasureString(line.String())
			if int(length) <= width {
				line.WriteRune(v)
			} else {
				buff = append(buff, line.String())
				line.Reset()
				line.WriteRune(v)
			}
		}
		buff = append(buff, line.String())
		line.Reset()
	}
	_, h := canvas.MeasureString("好")
	canvas = gg.NewContext(width+int(h*2+0.5), int(float64(len(buff)*3+1)/2*h+0.5))
	canvas.SetRGB(1, 1, 1)
	canvas.Clear()
	canvas.SetRGB(0, 0, 0)
	if err = canvas.LoadFontFace(font, float64(fontSize)); err != nil {
		return
	}
	for i, v := range buff {
		if v != "" {
			canvas.DrawString(v, float64(width)*0.01, float64((i+1)*3)/2*h)
		}
	}
	return canvas.Image(), nil
}
