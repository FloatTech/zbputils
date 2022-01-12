// Package txt2img 文字转图片
package txt2img

import (
	"encoding/base64"
	"image/jpeg"
	"io"
	"os"
	"strings"

	"github.com/fogleman/gg"
	"github.com/mattn/go-runewidth"
	log "github.com/sirupsen/logrus"

	"github.com/FloatTech/zbputils/binary"
	"github.com/FloatTech/zbputils/file"
	"github.com/FloatTech/zbputils/process"
)

const (
	whitespace = "\t\n\r\x0b\x0c"
	// FontPath 通用字体路径
	FontPath = "data/Font/"
	// FontFile 苹方字体
	FontFile = FontPath + "regular.ttf"
	// BoldFontFile 粗体苹方字体
	BoldFontFile = FontPath + "regular-bold.ttf"
)

// 加载数据库
func init() {
	go func() {
		process.SleepAbout1sTo2s()
		_ = os.MkdirAll(FontPath, 0755)
		_, _ = file.GetLazyData(FontFile, false, true)
		_, _ = file.GetLazyData(BoldFontFile, false, true)
	}()
}

type TxtCanvas struct {
	canvas *gg.Context
}

// RenderToBase64 文字转base64
func RenderToBase64(text string, width, fontSize int) (base64Bytes []byte, err error) {
	txtc, err := Render(text, width, fontSize)
	if err != nil {
		log.Println("[txt2img]:", err)
		return nil, err
	}
	base64Bytes, err = txtc.ToBase64()
	if err != nil {
		log.Println("[txt2img]:", err)
		return nil, err
	}
	return
}

// Render 文字转图片
func Render(text string, width, fontSize int) (txtc TxtCanvas, err error) {
	buff := make([]string, 0)
	line := ""
	count := 0
	for _, v := range text {
		c := string(v)
		if strings.Contains(whitespace, c) {
			buff = append(buff, strings.TrimSpace(line))
			count = 0
			line = ""
			continue
		}
		if count <= width {
			line += c
			count += runewidth.StringWidth(c)
		} else {
			buff = append(buff, line)
			line = c
			count = runewidth.StringWidth(c)
		}
	}

	txtc.canvas = gg.NewContext((fontSize+4)*width/2, (len(buff)+2)*fontSize)
	txtc.canvas.SetRGB(1, 1, 1)
	txtc.canvas.Clear()
	txtc.canvas.SetRGB(0, 0, 0)
	if err = txtc.canvas.LoadFontFace(FontFile, float64(fontSize)); err != nil {
		log.Errorln("[txt2img]:", err)
		return
	}
	for i, v := range buff {
		if v != "" {
			txtc.canvas.DrawString(v, float64(width/2), float64((i+2)*fontSize))
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
	if err = jpeg.Encode(encoder, txtc.canvas.Image(), &opt); err != nil {
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
		_ = jpeg.Encode(w, txtc.canvas.Image(), &jpeg.Options{Quality: 70})
	})
}

// WriteTo gg内容写入 Writer
// 使用完 data 后必须调用 cl 放回缓冲区
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
