// Package txt2img 文字转图片
package txt2img

import (
	"encoding/base64"
	"image/jpeg"
	"io"
	"os"
	"strings"

	"github.com/fogleman/gg"
	log "github.com/sirupsen/logrus"

	"github.com/FloatTech/zbputils/binary"
	"github.com/FloatTech/zbputils/file"
)

const (
	// FontPath 通用字体路径
	FontPath = "data/Font/"
	// FontFile 苹方字体
	FontFile = FontPath + "regular.ttf"
	// BoldFontFile 粗体苹方字体
	BoldFontFile = FontPath + "regular-bold.ttf"
	// SakuraFontFile ...
	SakuraFontFile = FontPath + "sakura.ttf"
	// ConsolasFontFile ...
	ConsolasFontFile = FontPath + "consolas.ttf"
)

// 加载数据库
func init() {
	_ = os.MkdirAll(FontPath, 0755)
}

type TxtCanvas struct {
	Canvas *gg.Context
}

// RenderToBase64 文字转base64
func RenderToBase64(text, font string, width, fontSize int) (base64Bytes []byte, err error) {
	txtc, err := Render(text, font, width, fontSize)
	if err != nil {
		log.Println("[txt2img]", err)
		return nil, err
	}
	base64Bytes, err = txtc.ToBase64()
	if err != nil {
		log.Println("[txt2img]", err)
		return nil, err
	}
	return
}

// Render 文字转图片 width 是图片宽度
func Render(text, font string, width, fontSize int) (txtc TxtCanvas, err error) {
	_, err = file.GetLazyData(font, false, true)
	if err != nil {
		return
	}

	txtc.Canvas = gg.NewContext(width, fontSize) // fake
	if err = txtc.Canvas.LoadFontFace(font, float64(fontSize)); err != nil {
		log.Errorln("[txt2img]", err)
		return
	}

	buff := make([]string, 0)
	for _, s := range strings.Split(text, "\n") {
		line := ""
		for _, v := range s {
			length, _ := txtc.Canvas.MeasureString(line)
			if int(length) <= width {
				line += string(v)
			} else {
				buff = append(buff, line)
				line = string(v)
			}
		}
		buff = append(buff, line)
	}

	_, h := txtc.Canvas.MeasureString("好")
	txtc.Canvas = gg.NewContext(width+int(h*2+0.5), int(float64(len(buff)*3+1)/2*h+0.5))
	txtc.Canvas.SetRGB(1, 1, 1)
	txtc.Canvas.Clear()
	txtc.Canvas.SetRGB(0, 0, 0)
	if err = txtc.Canvas.LoadFontFace(font, float64(fontSize)); err != nil {
		log.Errorln("[txt2img]", err)
		return
	}

	for i, v := range buff {
		if v != "" {
			txtc.Canvas.DrawString(v, float64(width)*0.01, float64((i+1)*3)/2*h)
		}
	}
	return
}

// ToBase64 gg内容转为base64
func (txtc TxtCanvas) ToBase64() (base64Bytes []byte, err error) {
	buffer := binary.SelectWriter()
	encoder := base64.NewEncoder(base64.StdEncoding, buffer)
	var opt jpeg.Options
	opt.Quality = 70
	if err = jpeg.Encode(encoder, txtc.Canvas.Image(), &opt); err != nil {
		return nil, err
	}
	encoder.Close()
	base64Bytes = buffer.Bytes()
	binary.PutWriter(buffer)
	return
}

// ToBytes gg内容转为 []byte
// 使用完 data 后必须调用 cl 放回缓冲区
func (txtc TxtCanvas) ToBytes() (data []byte, cl func()) {
	return binary.OpenWriterF(func(w *binary.Writer) {
		_ = jpeg.Encode(w, txtc.Canvas.Image(), &jpeg.Options{Quality: 70})
	})
}

// WriteTo gg内容写入 Writer
func (txtc TxtCanvas) WriteTo(f io.Writer) (n int64, err error) {
	data, cl := txtc.ToBytes()
	defer cl()
	if len(data) > 0 {
		var c int
		c, err = f.Write(data)
		return int64(c), err
	}
	return
}
