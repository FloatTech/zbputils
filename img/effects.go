package img

import (
	"image/color"
	"math"

	"github.com/disintegration/imaging"
)

// AdjustBrightness 亮度(-100, 100)
func (dst *ImgFactory) AdjustBrightness(s float64) *ImgFactory {
	s = math.Min(math.Max(s, -100.0), 100.0)
	b := dst.Im.Bounds()
	for y1 := b.Min.Y; y1 <= b.Max.Y; y1++ {
		for x1 := b.Min.X; x1 <= b.Max.X; x1++ {
			a := dst.Im.At(x1, y1)
			c := color.NRGBAModel.Convert(a).(color.NRGBA)
			f := 255.0 * s / 100.0
			c.R = floatUint8(f + float64(c.R))
			c.G = floatUint8(f + float64(c.R))
			c.B = floatUint8(f + float64(c.R))
			dst.Im.Set(x1, y1, c)
		}
	}

	return dst
}

// AdjustContrast 对比度(-100, 100)
func (dst *ImgFactory) AdjustContrast(a float64) *ImgFactory {
	return &ImgFactory{
		Im: imaging.AdjustContrast(dst.Im, a),
		W:  dst.W,
		H:  dst.H,
	}
}

// AdjustSaturation 饱和度(-100, 100)
func (dst *ImgFactory) AdjustSaturation(a float64) *ImgFactory {
	return &ImgFactory{
		Im: imaging.AdjustSaturation(dst.Im, a),
		W:  dst.W,
		H:  dst.H,
	}
}

// Sharpen 锐化
func (dst *ImgFactory) Sharpen(a float64) *ImgFactory {
	return &ImgFactory{
		Im: imaging.Sharpen(dst.Im, a),
		W:  dst.W,
		H:  dst.H,
	}
}

// Blur 模糊图像 正数
func (dst *ImgFactory) Blur(a float64) *ImgFactory {
	return &ImgFactory{
		Im: imaging.Blur(dst.Im, a),
		W:  dst.W,
		H:  dst.H,
	}
}

// Grayscale 灰度
func (dst *ImgFactory) Grayscale() *ImgFactory {
	b := dst.Im.Bounds()
	for y1 := b.Min.Y; y1 <= b.Max.Y; y1++ {
		for x1 := b.Min.X; x1 <= b.Max.X; x1++ {
			a := dst.Im.At(x1, y1)
			c := color.NRGBAModel.Convert(a).(color.NRGBA)
			f := 0.299*float64(c.R) + 0.587*float64(c.G) + 0.114*float64(c.B)
			c.R = floatUint8(f)
			c.G = floatUint8(f)
			c.B = floatUint8(f)
			dst.Im.Set(x1, y1, c)
		}
	}
	return dst
}

// Invert 反色
func (dst *ImgFactory) Invert() *ImgFactory {
	b := dst.Im.Bounds()
	for y1 := b.Min.Y; y1 <= b.Max.Y; y1++ {
		for x1 := b.Min.X; x1 <= b.Max.X; x1++ {
			a := dst.Im.At(x1, y1)
			c := color.NRGBAModel.Convert(a).(color.NRGBA)
			c.R = 255 - c.R
			c.G = 255 - c.G
			c.B = 255 - c.B
			dst.Im.Set(x1, y1, c)
		}
	}
	return dst
}

// Convolve3x3 浮雕
func (dst *ImgFactory) Convolve3x3() *ImgFactory {
	return &ImgFactory{
		Im: imaging.Convolve3x3(
			dst.Im,
			[9]float64{
				-1, -1, 0,
				-1, 1, 1,
				0, 1, 1,
			},
			nil,
		),
		W: dst.W,
		H: dst.H,
	}
}
