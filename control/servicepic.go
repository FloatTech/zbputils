package control

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"sync"

	"github.com/Coloured-glaze/gg"
	"github.com/FloatTech/zbputils/img"
	"github.com/FloatTech/zbputils/img/text"

	// 图片输出
	"image"

	"github.com/FloatTech/floatbox/file"
)

type mpic struct {
	kanban   string      // 看板娘图片路径
	kanbanW  int         // 看板娘图片宽度
	kanbanH  int         // 看板娘图片高度
	status   string      // 启用状态
	status2  bool        // 启用状态
	double   bool        // 双列排版
	plugin   string      // 插件名
	font1    string      // 字体1
	font2    string      // 字体2
	info     []string    // 插件信息
	info2    []string    // 插件信息2
	multiple float64     // 图片拓展倍数
	fontSize float64     // 字体大小
	im       image.Image // 图片
}

type titleColor struct {
	r, g, b int  // 颜色
	s       bool // 是否随机
}

type location struct {
	lastH            int     // 上一个高度
	drawX, maxTwidth float64 // 文字边距
	rlineX, rlineY   float64 // 宽高记录
	rtitleW          float64 // 标题位置
}

func init() {
	_, err := file.GetLazyData("data/Control/kanban.png", true)
	if err != nil {
		log.Println(err)
	}
	_, err = file.GetLazyData(text.BoldFontFile, true)
	if err != nil {
		log.Println(err)
	}
	_, err = file.GetLazyData(text.SakuraFontFile, true)
	if err != nil {
		log.Println(err)
	}
}

// 返回菜单图片
func dyna(mp *mpic, lt location) (image.Image, error) {
	var err error
	fontSize := mp.fontSize                                               // 图片宽度和字体大小
	var one = gg.NewContext(int(1024.0*2.5), mp.kanbanH+int(mp.fontSize)) // 新图像
	one.SetRGB255(255, 255, 255)
	one.Clear()
	if mp.double {
		one.DrawImage(mp.im, (one.W()/2-mp.kanbanH)/2, 0) // 放入看板娘
		lt.rtitleW = float64(one.W()) / 2
	} else {
		lt.rtitleW = 1100
	}
	if err = one.LoadFontFace(mp.font2, fontSize*2); err != nil {
		return nil, err
	}
	one.SetRGBA255(55, 55, 55, 255)                             // 字体颜色
	one.DrawString(mp.plugin, lt.rtitleW*1.35, fontSize*2)      // 插件名
	if err = one.LoadFontFace(mp.font1, fontSize); err != nil { // 加载字体
		return nil, err
	}
	one.DrawRoundedRectangle(lt.rtitleW+32, fontSize-5, fontSize*4.5, fontSize*1.5, 10) // 创建圆角矩形
	if mp.status2 {                                                                     // 如果启用
		one.SetRGBA255(15, 200, 15, 200)   // 设置背景颜色
		one.Fill()                         // 填充
		one.SetRGBA255(255, 255, 255, 255) // 设置白色
	} else {
		one.SetRGBA255(221, 221, 221, 255)
		one.Fill()
		one.SetRGBA255(55, 55, 55, 255) // 设置黑色
	}
	one.DrawString(mp.status, lt.rtitleW+40, fontSize*2) // 绘制启用状态
	return createPic(one, mp, lt)
}

// 创建图片
func createPic(one *gg.Context, mp *mpic, lt location) (image.Image, error) {
	var wg sync.WaitGroup
	ch := make(chan image.Image, 1)
	errch := make(chan error, 2)
	titlec := titleColor{0, 0, 0, true}
	if mp.double {
		titlec = randColor(titleColor{0, 0, 0, false})
		ch2 := make(chan image.Image, 1)
		wg.Add(1)
		go func() {
			im, err := createPic2(mp, lt, titlec, &wg, mp.info)
			if err != nil {
				errch <- err
				ch <- nil
				return
			}
			errch <- nil
			ch <- im
		}()
		wg.Add(1)
		go func() {
			im, err := createPic2(mp, lt, titlec, &wg, mp.info2)
			if err != nil {
				errch <- err
				ch2 <- nil
				return
			}
			errch <- nil
			ch2 <- im
		}()
		wg.Wait()
		close(errch)
		close(ch)
		close(ch2)
		for err := range errch {
			if err != nil {
				return nil, err
			}
		}
		var imgs [2]image.Image
		for im := range ch {
			imgs[0] = im
		}
		for im2 := range ch2 {
			imgs[1] = im2
		}
		imRY := imgs[0].Bounds().Dy() // 右边 图像的高度
		imLY := imgs[1].Bounds().Dy() // 左边 图像的高度

		max, min, left := 0, 0, false
		tmpRY := imRY + 100
		tmpLY := imLY + mp.kanbanH + 5

		if tmpRY > tmpLY {
			max = tmpRY
			min = tmpLY // 剩余空间最大
			left = true
		} else {
			max = tmpLY
			min = tmpRY
		}
		if max > one.H() {
			itmp := gg.NewContext(one.W(), max+int(mp.fontSize)) // 高度
			itmp.SetRGB255(255, 255, 255)
			itmp.Clear()
			itmp.DrawImage(one.Image(), 0, 0)
			one = gg.NewContextForImage(itmp.Image())
		}
		min = one.H() - min //
		if min < mp.kanbanH {
			itmp := gg.NewContext(one.W(), one.H()+(mp.kanbanH-min)+int(mp.fontSize)) // 高度
			itmp.SetRGB255(255, 255, 255)
			itmp.Clear()
			itmp.DrawImage(one.Image(), 0, 0)
			one = gg.NewContextForImage(itmp.Image())
		}
		one.DrawImage(imgs[0], 1250, 0) // 最终的绘制位置
		one.DrawImage(imgs[1], 0, mp.kanbanH+5)
		if left {
			one.DrawImage(mp.im, (one.W()/2-mp.kanbanH)/2-50, tmpLY+5) // 放入看板娘
		} else {
			one.DrawImage(mp.im, (one.W()/2)+(one.W()/2-mp.kanbanH)/2, imRY+5) // 放入看板娘
		}

	} else { //==================================================================>>

		titlec = randColor(titlec)
		wg.Add(1)
		go func() {
			im, err := createPic2(mp, lt, titlec, &wg, mp.info)
			if err != nil {
				errch <- err
				ch <- nil
				return
			}
			errch <- nil
			ch <- im
		}()
		wg.Wait()
		close(errch)
		close(ch)
		for err := range errch {
			if err != nil {
				return nil, err
			}
		}
		var imgs [2]image.Image
		for im := range ch {
			imgs[0] = im
		}
		imY := imgs[0].Bounds().Dy()
		if imY+int(mp.fontSize) > one.H() {
			itmp := gg.NewContext(one.W(), imY) // 高度
			itmp.SetRGB255(255, 255, 255)
			itmp.Clear()
			itmp.DrawImage(one.Image(), 0, 0)
			one = gg.NewContextForImage(itmp.Image())
		}
		one.DrawImage(mp.im, (one.W()/2-mp.kanbanH)/2-50, 0) // 放入看板娘
		one.DrawImage(imgs[0], int(lt.rtitleW), 50)          // 最终的绘制位置
		if one.H() > mp.kanbanH*3 {                          // 超出三倍高度
			one.DrawImage(mp.im, (one.W()/2-mp.kanbanH)/2-50, imY-mp.kanbanH-int(mp.fontSize*5)) // 放入看板娘
		}
	}
	return one.Image(), nil
}

// 创建图片
func createPic2(mp *mpic, lt location, titlec titleColor, wg *sync.WaitGroup,
	info []string) (image.Image, error) {
	defer wg.Done()
	fontSize := mp.fontSize
	var one = gg.NewContext(1280, mp.kanbanH)
	if err := one.LoadFontFace(mp.font1, fontSize); err != nil { // 加载字体
		return nil, err
	}
	for i := 0; i < len(info); i++ { // 遍历文字切片
		lineTexts := make([]string, 0, len(info[i]))
		lineText, textW, textH, tmpw := "", 0.0, 0.0, 0.0

		if mp.double {
			if strings.Contains(info[i], ": ● ") || strings.Contains(info[i], ": ○ ") {
				titlec = randColor(titlec) // 随机一次颜色
			}
		}
		for len(info[i]) > 0 {
			lineText, tmpw = truncate(one, info[i], lt.maxTwidth)
			lineTexts = append(lineTexts, lineText)
			if tmpw > textW {
				textW = tmpw
			}
			if len(lineText) >= len(info[i]) {
				break // 如果写入的文本大于等于本次写入的文本 则跳出
			}
			textH += fontSize * 1.3           // 截断一次加一行高度
			info[i] = info[i][len(lineText):] // 丢弃已经写入的文字并重新赋值
		}
		threeW, threeH := textW+fontSize, (textH + (fontSize * 1.2)) // 圆角矩形宽度和高度
		lt.drawX = lt.rlineX + 30                                    // 圆角矩形位置宽度
		if int(lt.rlineX+textW)+int(fontSize*2) > one.W() {          // 越界
			goto label
		}
		one.DrawRoundedRectangle(lt.drawX, lt.rlineY, threeW, threeH, 20.0) // 创建圆角矩形
		drawsc(one, titlec, fontSize, lt.drawX, lt.rlineY, lineTexts)
		lt.rlineX += threeW + fontSize/2 // 添加后加一次宽度
		lt.lastH = int(threeH)

		continue // 跳出本次循环
	label:

		lt.rlineY += float64(lt.lastH) + fontSize/4 // 加一次高度
		lt.rlineX = 5                                      // 重置宽度位置
		if threeH+lt.rlineY+fontSize >= float64(one.H()) { // 超出最大高度则进行加高
			itmp := gg.NewContext(one.W(), int(lt.rlineY+threeH*mp.multiple)) // 高度
			itmp.DrawImage(one.Image(), 0, 0)
			one = gg.NewContextForImage(itmp.Image())
			if err := one.LoadFontFace(mp.font1, mp.fontSize); err != nil { // 加载字体
				return nil, err
			}
		}
		lt.drawX = lt.rlineX + 30                                           // 圆角矩形位置宽度
		one.DrawRoundedRectangle(lt.drawX, lt.rlineY, threeW, threeH, 20.0) // 创建圆角矩形
		drawsc(one, titlec, fontSize, lt.drawX, lt.rlineY, lineTexts)
		lt.rlineX += threeW + fontSize/2 // 添加后加一次宽度
		lt.lastH = int(threeH)
	}
	return one.Image(), nil
}

// 绘制文字
func drawsc(one *gg.Context, titlec titleColor, fontSize, drawX, rlineY float64, lineTexts []string) {
	if titlec.s {
		titlec = randColor(titlec)
	}
	one.SetRGBA255(titlec.r, titlec.g, titlec.b, 85)
	one.Fill() // 填充颜色
	one.SetRGBA255(55, 55, 55, 255)
	h := fontSize + rlineY - 3
	for i := range lineTexts { // 逐行绘制文字
		one.DrawString(lineTexts[i], drawX+fontSize/2, h)
		h += fontSize + (fontSize / 4)
	}
}

// 填充颜色
func randColor(titlec titleColor) titleColor {
	//	rand.Seed(time.Now().UnixNano())
	titlec.r = rand.Intn(245) // 随机颜色
	titlec.g = rand.Intn(245)
	titlec.b = rand.Intn(245)
r:
	if titlec.r < 15 && titlec.g < 15 && titlec.b < 15 {
		rand.Seed(rand.Int63n(99999999))
		titlec.r = rand.Intn(245)
		titlec.g = rand.Intn(245)
		titlec.b = rand.Intn(245)
		goto r
	} else if titlec.r > 210 && titlec.g > 210 && titlec.b > 210 {
		rand.Seed(rand.Int63n(99999999))
		titlec.r = rand.Intn(245)
		titlec.g = rand.Intn(245)
		titlec.b = rand.Intn(245)
		goto r
	}
	return titlec
}

// 截断文字
func truncate(one *gg.Context, text string, maxW float64) (string, float64) {
	var tmp strings.Builder
	tmp.Grow(len(text))
	res, w := make([]rune, 0, len(text)), 0.0
	for _, r := range text {
		tmp.WriteRune(r)
		width, _ := one.MeasureString(tmp.String()) // 获取文字宽度
		if width > maxW {                           // 如果宽度大于文字边距
			break //跳出
		} else {
			w = width
			res = append(res, r) // 写入
		}
	}
	return string(res), w
}

// 编码看板娘图片和加载字体
func loadpic(mp *mpic) error {
	if !file.IsExist(mp.font1) { // 获取字体
		return errors.New("文件 " + mp.font1 + " 不存在")
	}
	if !file.IsExist(mp.font2) { // 获取字体
		return errors.New("文件 " + mp.font2 + " 不存在")
	}
	data, err := ioutil.ReadFile(mp.kanban)
	if err != nil {
		return err
	}
	width, height, err := gg.GetImgWH(data, mp.kanban) // 解析图片的宽高信息
	if err != nil {
		return err
	}
	if width > 1024 { // 图片超出大小则进行限制
		mp.im, _, err = image.Decode(bytes.NewReader(data))
		if err != nil {
			return err
		}
		mp.im = img.Limit(mp.im, 1024, 1500)
		mp.kanbanW, mp.kanbanH = 1024, 1500
		return nil
	}
	mp.im, _, err = image.Decode(bytes.NewReader(data))
	if err != nil {
		return err
	}
	mp.kanbanW, mp.kanbanH = width, height
	return nil
}
