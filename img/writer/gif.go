package writer

import (
	"image/gif"
	"os"
)

// SaveGIF2Path 保存 gif 到 path
func SaveGIF2Path(path string, g *gif.GIF) error {
	f, err := os.Create(path) // 创建文件
	if err == nil {
		gif.EncodeAll(f, g) // 写入
		f.Close()           // 关闭文件
	}
	return err
}
