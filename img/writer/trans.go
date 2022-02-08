package writer

import (
	"encoding/base64"
	"image"
	"image/jpeg"
	"io"

	"github.com/FloatTech/zbputils/binary"
)

// ToBase64 img 内容转为base64
func ToBase64(img image.Image) (base64Bytes []byte, err error) {
	buffer := binary.SelectWriter()
	encoder := base64.NewEncoder(base64.StdEncoding, buffer)
	var opt jpeg.Options
	opt.Quality = 70
	if err = jpeg.Encode(encoder, img, &opt); err != nil {
		return nil, err
	}
	encoder.Close()
	base64Bytes = buffer.Bytes()
	binary.PutWriter(buffer)
	return
}

// ToBytes img 内容转为 []byte
// 使用完 data 后必须调用 cl 放回缓冲区
func ToBytes(img image.Image) (data []byte, cl func()) {
	return binary.OpenWriterF(func(w *binary.Writer) {
		_ = jpeg.Encode(w, img, &jpeg.Options{Quality: 70})
	})
}

// WriteTo img 内容写入 Writer
func WriteTo(img image.Image, f io.Writer) (n int64, err error) {
	data, cl := ToBytes(img)
	defer cl()
	if len(data) > 0 {
		var c int
		c, err = f.Write(data)
		return int64(c), err
	}
	return
}
