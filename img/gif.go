package img

import (
	"image"
	"image/color"
	"image/draw"
	"image/gif"

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
	g := &gif.GIF{
		Image:    make([]*image.Paletted, len(im)),
		Delay:    make([]int, len(im)),
		Disposal: make([]byte, len(im)),
	}
	for i, stc := range im {
		g.Image[i] = GetPaletted(stc)          // 每帧图片
		g.Delay[i] = delay                     // 每帧间隔, 1=10毫秒
		g.Disposal[i] = gif.DisposalBackground // 透明图片需要设置
	}
	return g
}
