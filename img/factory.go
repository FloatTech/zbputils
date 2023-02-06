package img

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/FloatTech/gg"
	"github.com/disintegration/imaging"
)

// Factory 处理中图像
type Factory struct {
	Im *image.NRGBA
	W  int
	H  int
}

// Clone 克隆
func (dst *Factory) Clone() *Factory {
	var src Factory
	src.Im = image.NewNRGBA(image.Rect(0, 0, dst.W, dst.H))
	draw.Over.Draw(src.Im, src.Im.Bounds(), dst.Im, dst.Im.Bounds().Min)
	src.W = dst.W
	src.H = dst.H
	return &src
}

// Reshape 变形
func (dst *Factory) Reshape(w, h int) *Factory {
	dst = Size(dst.Im, w, h)
	return dst
}

// FlipH 水平翻转
func (dst *Factory) FlipH() *Factory {
	return &Factory{
		Im: imaging.FlipH(dst.Im),
		W:  dst.W,
		H:  dst.H,
	}
}

// FlipV 垂直翻转
func (dst *Factory) FlipV() *Factory {
	return &Factory{
		Im: imaging.FlipV(dst.Im),
		W:  dst.W,
		H:  dst.H,
	}
}

// InsertUp 上部插入图片
func (dst *Factory) InsertUp(im image.Image, w, h, x, y int) *Factory {
	im1 := Size(im, w, h).Im
	// 叠加图片
	draw.Over.Draw(dst.Im, dst.Im.Bounds(), im1, im1.Bounds().Min.Sub(image.Pt(x, y)))
	return dst
}

// InsertUpC 上部插入图片 x,y是中心点
func (dst *Factory) InsertUpC(im image.Image, w, h, x, y int) *Factory {
	im1 := Size(im, w, h)
	// 叠加图片
	draw.Over.Draw(dst.Im, dst.Im.Bounds(), im1.Im, im1.Im.Bounds().Min.Sub(image.Pt(x-im1.W/2, y-im1.H/2)))
	return dst
}

// InsertBottom 底部插入图片
func (dst *Factory) InsertBottom(im image.Image, w, h, x, y int) *Factory {
	im1 := Size(im, w, h).Im
	dc := dst.Clone()
	dst = NewFactory(dst.W, dst.H, color.NRGBA{0, 0, 0, 0})
	draw.Over.Draw(dst.Im, dst.Im.Bounds(), im1, im1.Bounds().Min.Sub(image.Pt(x, y)))
	draw.Over.Draw(dst.Im, dst.Im.Bounds(), dc.Im, dc.Im.Bounds().Min)
	return dst
}

// InsertBottomC 底部插入图片 x,y是中心点
func (dst *Factory) InsertBottomC(im image.Image, w, h, x, y int) *Factory {
	im1 := Size(im, w, h)
	dc := dst.Clone()
	dst = NewFactory(dst.W, dst.H, color.NRGBA{0, 0, 0, 0})
	draw.Over.Draw(dst.Im, dst.Im.Bounds(), im1.Im, im1.Im.Bounds().Min.Sub(image.Pt(x-im1.W/2, y-im1.H/2)))
	draw.Over.Draw(dst.Im, dst.Im.Bounds(), dc.Im, dc.Im.Bounds().Min)
	return dst
}

// Circle 获取圆图
func (dst *Factory) Circle(r int) *Factory {
	if r == 0 {
		r = dst.H / 2
	}
	dst = dst.Reshape(2*r, 2*r)
	b := dst.Im.Bounds()
	for y1 := b.Min.Y; y1 < b.Max.Y; y1++ {
		for x1 := b.Min.X; x1 < b.Max.X; x1++ {
			if (x1-r)*(x1-r)+(y1-r)*(y1-r) > r*r {
				dst.Im.Set(x1, y1, color.NRGBA{0, 0, 0, 0})
			}
		}
	}
	return dst
}

// Clip 剪取方图
func (dst *Factory) Clip(w, h, x, y int) *Factory {
	dst.Im = dst.Im.SubImage(image.Rect(x, y, x+w, y+h)).(*image.NRGBA)
	dst.W = w
	dst.H = h
	return dst
}

// ClipCircleFix 裁取圆图
func (dst *Factory) ClipCircleFix(x, y, r int) *Factory {
	dst = dst.Clip(2*r, 2*r, x-r, y-r)
	b := dst.Im.Bounds()
	for y1 := b.Min.Y; y1 < b.Max.Y; y1++ {
		for x1 := b.Min.X; x1 < b.Max.X; x1++ {
			if (x1-x)*(x1-x)+(y1-y)*(y1-y) > r*r {
				dst.Im.Set(x1, y1, color.NRGBA{0, 0, 0, 0})
			}
		}
	}
	return dst
}

// ClipCircle 扣取圆
func (dst *Factory) ClipCircle(x, y, r int) *Factory {
	//  dc := dst.Clip(x-r, y-r, 2*r, 2*r)
	b := dst.Im.Bounds()
	for y1 := b.Min.Y; y1 < b.Max.Y; y1++ {
		for x1 := b.Min.X; x1 < b.Max.X; x1++ {
			if (x1-x)*(x1-x)+(y1-y)*(y1-y) <= r*r {
				dst.Im.Set(x1, y1, color.NRGBA{0, 0, 0, 0})
			}
		}
	}
	return dst
}

// InsertText 插入文本
func (dst *Factory) InsertText(font string, size float64, col []int, x, y float64, txt string) *Factory {
	dc := gg.NewContextForImage(dst.Im)
	// 字体, 大小, 颜色, 位置
	err := dc.LoadFontFace(font, size)
	if err != nil {
		return dst
	}
	dc.SetRGBA255(col[0], col[1], col[2], col[3])
	dc.DrawString(txt, x, y)
	ds := dc.Image()
	draw.Over.Draw(dst.Im, dst.Im.Bounds(), ds, ds.Bounds().Min)
	return dst
}

// InsertUpG gif 上部插入图片
func (dst *Factory) InsertUpG(im []*image.NRGBA, w, h, x, y int) []*image.NRGBA {
	if len(im) == 0 {
		return nil
	}
	ims := make([]*image.NRGBA, len(im))
	for i, v := range im {
		ims[i] = dst.Clone().InsertUp(v, w, h, x, y).Im
	}
	return ims
}
