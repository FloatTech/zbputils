package img

import (
	"bufio"
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/png"
	"net/http"
	"os"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/sirupsen/logrus"
)

//  加载图片
func Load(path string) image.Image {
	var img image.Image
	if strings.HasPrefix(path, "http") {
		res, err := http.Get(path)
		if err != nil {
			logrus.Errorln("[img] load err:", err)
			return nil
		}
		// 获得get请求响应的reader对象
		reader := bufio.NewReaderSize(res.Body, 32*1024)
		// 读取路径
		img, _, err = image.Decode(reader)
		res.Body.Close()
		if err != nil {
			logrus.Errorln("[img] decode err:", err)
			return nil
		}
	} else {
		//  加载路径
		file, err := os.Open(path)
		if err == nil {
			// 读取路径
			img, _, err = image.Decode(file)
			file.Close()
			if err != nil {
				logrus.Errorln("[img] decode err:", err)
				return nil
			}
		} else {
			logrus.Errorln("[img] open file", path, "err:", err)
			return nil
		}
	}
	return img
}

// 设置底图
func NewFactory(w, h int, fillColor color.Color) *ImgFactory {
	var dst ImgFactory
	dst.W = w
	dst.H = h
	c := color.NRGBAModel.Convert(fillColor).(color.NRGBA)
	if (c == color.NRGBA{0, 0, 0, 0}) {
		dst.Im = image.NewNRGBA(image.Rect(0, 0, w, h))
	} else {
		dst.Im = &image.NRGBA{
			Pix:    bytes.Repeat([]byte{c.R, c.G, c.B, c.A}, w*h),
			Stride: 4 * w,
			Rect:   image.Rect(0, 0, w, h),
		}
	}
	return &dst
}

// 载入图片第一帧作底图
func LoadFirstFrame(path string, w, h int) *ImgFactory {
	return Size(Load(path), w, h)
}

//  加载图片每一帧图片
func LoadAllFrames(path string, w, h int) []*image.NRGBA {
	var res *http.Response
	var err error
	var im *gif.GIF
	var ims []*image.NRGBA
	if strings.HasPrefix(path, "http") {
		res, err = http.Get(path)
		if err != nil {
			logrus.Errorln("[img] load err:", err)
			return nil
		}
		//  获得get请求响应的reader对象
		reader := bufio.NewReaderSize(res.Body, 32*1024)
		// 读取路径
		im, err = gif.DecodeAll(reader)
		res.Body.Close()
		if err != nil {
			logrus.Errorln("[img] decode err:", err)
			return nil
		}
	} else {
		//  加载路径
		file, err := os.Open(path)
		if err == nil {
			im, err = gif.DecodeAll(file)
			file.Close()
			if err != nil {
				logrus.Errorln("[img] decode err:", err)
				return nil
			}
		} else {
			logrus.Errorln("[img] open file", path, "err:", err)
			return nil
		}
	}
	im0 := Size(Load(path), w, h)
	if err != nil {
		ims = append(ims, Size(Load(path), w, h).Im)
	} else {
		for _, v := range im.Image {
			im1 := Size(v, w, h).Im
			im2 := im0.InsertUp(im1, 0, 0, 0, 0).Clone()
			ims = append(ims, im2.Im)
		}
	}
	return ims
}

// 变形
func Size(im image.Image, w, h int) *ImgFactory {
	var dc ImgFactory
	// 修改尺寸
	if w > 0 && h > 0 {
		dc.W = w
		dc.H = h
		dc.Im = imaging.Resize(im, w, h, imaging.Lanczos)
	} else if w <= 0 && h <= 0 {
		dc.W = im.Bounds().Size().X
		dc.H = im.Bounds().Size().Y
		dc.Im = image.NewNRGBA(image.Rect(0, 0, dc.W, dc.H))
		draw.Over.Draw(dc.Im, dc.Im.Bounds(), im, im.Bounds().Min)
	} else if w == 0 {
		dc.H = h
		dc.W = h * im.Bounds().Size().X / im.Bounds().Size().Y
		dc.Im = imaging.Resize(im, dc.W, h, imaging.Lanczos)
	} else if h == 0 {
		dc.W = w
		dc.H = w * im.Bounds().Size().Y / im.Bounds().Size().X
		dc.Im = imaging.Resize(im, w, dc.H, imaging.Lanczos)
	}
	return &dc
}

// 旋转
func Rotate(img image.Image, angle float64, w, h int) *ImgFactory {
	im := Size(img, w, h)
	var dc ImgFactory
	dc.Im = imaging.Rotate(im.Im, angle, color.NRGBA{0, 0, 0, 0})
	dc.W = dc.Im.Bounds().Size().X
	dc.H = dc.Im.Bounds().Size().Y
	return &dc
}

// 横向合并图片
func MergeW(im []*image.NRGBA) *ImgFactory {
	dc := make([]*ImgFactory, len(im))
	h := im[0].Bounds().Size().Y
	w := 0
	for i, value := range im {
		dc[i] = Size(value, 0, h)
		w += dc[i].W
	}
	ds := NewFactory(w, h, color.NRGBA{0, 0, 0, 0})
	x := 0
	for _, value := range dc {
		ds = ds.InsertUp(value.Im, value.W, h, x, 0)
		x += value.W
	}
	return ds
}

// 纵向合并图片
func MergeH(im []*image.NRGBA) *ImgFactory {
	dc := make([]*ImgFactory, len(im))
	w := im[0].Bounds().Size().X
	h := 0
	for i, value := range im {
		dc[i] = Size(value, 0, w)
		h += dc[i].H
	}
	ds := NewFactory(w, h, color.NRGBA{0, 0, 0, 0})
	y := 0
	for _, value := range dc {
		ds = ds.InsertUp(value.Im, w, value.H, 0, y)
		y += value.H
	}
	return ds
}

// 文本框 字体，大小，颜色 ，背景色，文本
func Text(font string, size float64, col []int, col1 []int, txt string) *ImgFactory {
	var dst ImgFactory
	dc := gg.NewContext(10, 10)
	dc.SetRGBA255(0, 0, 0, 0)
	dc.Clear()
	dc.SetRGBA255(col[0], col[1], col[2], col[3])
	dc.LoadFontFace(font, size+size/2)
	w, h := dc.MeasureString(txt)
	w = w - size*2
	dc1 := gg.NewContext(int(w), int(h))
	dc1.SetRGBA255(col1[0], col1[1], col1[2], col1[3])
	dc1.Clear()
	dc1.SetRGBA255(col[0], col[1], col[2], col[3])
	dc1.LoadFontFace(font, size)
	dc1.DrawStringAnchored(txt, w/2, h/2, 0.5, 0.5)
	dst.Im = image.NewNRGBA(image.Rect(0, 0, int(w), int(h)))
	draw.Over.Draw(dst.Im, dst.Im.Bounds(), dc1.Image(), dc1.Image().Bounds().Min)
	dst.W, dst.H = dst.Im.Bounds().Size().X, dst.Im.Bounds().Size().Y
	return &dst
}

// 保存png
func SavePng(im image.Image, path string) error {
	f, err := os.Create(path) // 创建文件
	if err == nil {
		err = png.Encode(f, im) // 写入
		f.Close()
		if err != nil {
			logrus.Errorln("[img] save img err:", err)
		}
	} else {
		logrus.Errorln("[img] create file at", path, "err:", err)
	}
	return err
}

// float64转uint8
func floatUint8(a float64) uint8 {
	b := int64(a + 0.5)
	if b > 255 {
		return 255
	}
	if b > 0 {
		return uint8(b)
	}
	return 0
}
