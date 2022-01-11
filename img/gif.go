package img

import (
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"os"

	"github.com/ericpauley/go-quantize/quantize"
)

// GetPaletted 将image.Image转换为 *Paletted。最多256色
func GetPaletted(im image.Image) *image.Paletted {
	q := quantize.MedianCutQuantizer{AddTransparent: true}
	p := q.Quantize(make([]color.Color, 0, 256), im)
	cp := image.NewPaletted(image.Rect(0, 0, im.Bounds().Size().X, im.Bounds().Size().Y), p)
	draw.Src.Draw(cp, cp.Bounds(), im, im.Bounds().Min)
	return cp
}

// MergeGif 合并成gif,1 dealy 10毫秒
func MergeGif(delay int, im []*image.NRGBA) *gif.GIF {
	g := &gif.GIF{}

	for _, stc := range im {
		g.Image = append(g.Image, GetPaletted(stc))             // 每帧图片
		g.Delay = append(g.Delay, delay)                        // 每帧间隔，1=10毫秒
		g.Disposal = append(g.Disposal, gif.DisposalBackground) // 透明图片需要设置
	}
	g.LoopCount = 0
	return g
}

// SaveGif 保存gif
func SaveGif(g *gif.GIF, path string) error {
	f, err := os.Create(path) // 创建文件
	if err == nil {
		gif.EncodeAll(f, g) // 写入
		f.Close()           // 关闭文件
	}
	return err
}
