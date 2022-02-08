package writer

import (
	"image"
	"image/png"
	"os"
)

// SavePNG2Path 保存 png 到 path
func SavePNG2Path(path string, im image.Image) error {
	f, err := os.Create(path) // 创建文件
	if err == nil {
		err = png.Encode(f, im) // 写入
		_ = f.Close()
	}
	return err
}
