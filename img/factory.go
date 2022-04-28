package img

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
)

// 处理中图像
type ImgFactory struct {
	Im *image.NRGBA
	W  int
	H  int
}

// Clone 克隆
func (dst *ImgFactory) Clone() *ImgFactory {
	var src ImgFactory
	src.Im = image.NewNRGBA(image.Rect(0, 0, dst.W, dst.H))
	draw.Over.Draw(src.Im, src.Im.Bounds(), dst.Im, dst.Im.Bounds().Min)
	src.W = dst.W
	src.H = dst.H
	return &src
}

// Reshape 变形
func (dst *ImgFactory) Reshape(w, h int) *ImgFactory {
	dst = Size(dst.Im, w, h)
	return dst
}

// FlipH 水平翻转
func (dst *ImgFactory) FlipH() *ImgFactory {
	return &ImgFactory{
		Im: imaging.FlipH(dst.Im),
		W:  dst.W,
		H:  dst.H,
	}
}

// FlipV 垂直翻转
func (dst *ImgFactory) FlipV() *ImgFactory {
	return &ImgFactory{
		Im: imaging.FlipV(dst.Im),
		W:  dst.W,
		H:  dst.H,
	}
}

// InsertUp 上部插入图片
func (dst *ImgFactory) InsertUp(im image.Image, w, h, x, y int) *ImgFactory {
	im1 := Size(im, w, h).Im
	// 叠加图片
	draw.Over.Draw(dst.Im, dst.Im.Bounds(), im1, im1.Bounds().Min.Sub(image.Pt(x, y)))
	return dst
}

// InsertUpC 上部插入图片 x,y是中心点
func (dst *ImgFactory) InsertUpC(im image.Image, w, h, x, y int) *ImgFactory {
	im1 := Size(im, w, h)
	// 叠加图片
	draw.Over.Draw(dst.Im, dst.Im.Bounds(), im1.Im, im1.Im.Bounds().Min.Sub(image.Pt(x-im1.W/2, y-im1.H/2)))
	return dst
}

// InsertBottom 底部插入图片
func (dst *ImgFactory) InsertBottom(im image.Image, w, h, x, y int) *ImgFactory {
	im1 := Size(im, w, h).Im
	dc := dst.Clone()
	dst = NewFactory(dst.W, dst.H, color.NRGBA{0, 0, 0, 0})
	draw.Over.Draw(dst.Im, dst.Im.Bounds(), im1, im1.Bounds().Min.Sub(image.Pt(x, y)))
	draw.Over.Draw(dst.Im, dst.Im.Bounds(), dc.Im, dc.Im.Bounds().Min)
	return dst
}

// InsertBottomC 底部插入图片 x,y是中心点
func (dst *ImgFactory) InsertBottomC(im image.Image, w, h, x, y int) *ImgFactory {
	im1 := Size(im, w, h)
	dc := dst.Clone()
	dst = NewFactory(dst.W, dst.H, color.NRGBA{0, 0, 0, 0})
	draw.Over.Draw(dst.Im, dst.Im.Bounds(), im1.Im, im1.Im.Bounds().Min.Sub(image.Pt(x-im1.W/2, y-im1.H/2)))
	draw.Over.Draw(dst.Im, dst.Im.Bounds(), dc.Im, dc.Im.Bounds().Min)
	return dst
}

// Circle 获取圆图
func (dst *ImgFactory) Circle(r int) *ImgFactory {
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
func (dst *ImgFactory) Clip(w, h, x, y int) *ImgFactory {
	dst.Im = dst.Im.SubImage(image.Rect(x, y, x+w, y+h)).(*image.NRGBA)
	dst.W = w
	dst.H = h
	return dst
}

// ClipCircleFix 裁取圆图
func (dst *ImgFactory) ClipCircleFix(x, y, r int) *ImgFactory {
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
func (dst *ImgFactory) ClipCircle(x, y, r int) *ImgFactory {
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
func (dst *ImgFactory) InsertText(font string, size float64, col []int, x, y float64, txt string) *ImgFactory {
	dc := gg.NewContextForImage(dst.Im)
	// 字体, 大小, 颜色, 位置
	dc.LoadFontFace(font, size)
	dc.SetRGBA255(col[0], col[1], col[2], col[3])
	dc.DrawString(txt, x, y)
	ds := dc.Image()
	draw.Over.Draw(dst.Im, dst.Im.Bounds(), ds, ds.Bounds().Min)
	return dst
}

// InsertUpG gif 上部插入图片
func (dst *ImgFactory) InsertUpG(im []*image.NRGBA, w, h, x, y int) []*image.NRGBA {
	var ims []*image.NRGBA
	for _, v := range im {
		dc := dst.Clone().InsertUp(v, w, h, x, y).Im
		ims = append(ims, dc)
	}
	return ims
}
