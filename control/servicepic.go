package control

import (
	"image"
	"os"
	"sync"

	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/math"
	"github.com/FloatTech/gg"
	ctrl "github.com/FloatTech/zbpctrl"
	zero "github.com/wdvxdr1123/ZeroBot"

	"github.com/FloatTech/zbputils/img/text"
)

const (
	bannerpath = "data/zbpbanner/"
	kanbanpath = "data/Control/"
	bannerurl  = "https://github.com/FloatTech/zbpbanner/raw/main/"
)

type plugininfo struct {
	name   string
	brief  string
	status bool
}

var (
	// 字体 GlowSans 数据
	glowsd []byte
	// 字体 Impact 数据
	impactd []byte
	// 字体 Torus 数据
	torussd []byte
	// 字体 Yuruka 数据
	yurukasd []byte
)

func init() {
	err := os.MkdirAll(bannerpath, 0755)
	if err != nil {
		panic(err)
	}
	_, err = file.GetLazyData(kanbanpath+"kanban.png", Md5File, false)
	if err != nil {
		panic(err)
	}
	_, err = file.GetLazyData(kanbanpath+"zbpuwu.png", Md5File, false)
	if err != nil {
		panic(err)
	}
	glowsd, err = file.GetLazyData(text.GlowSansFontFile, Md5File, true)
	if err != nil {
		panic(err)
	}
	impactd, err = file.GetLazyData(text.ImpactFontFile, Md5File, true)
	if err != nil {
		panic(err)
	}
	torussd, err = file.GetLazyData(text.TorusFontFile, Md5File, true)
	if err != nil {
		panic(err)
	}
	yurukasd, err = file.GetLazyData(text.YurukaFontFile, Md5File, true)
	if err != nil {
		panic(err)
	}
}

func renderservepicof(gid int64) (img image.Image, err error) {
	pluginlist := make([]plugininfo, len(priomap))
	ForEachByPrio(func(i int, manager *ctrl.Control[*zero.Ctx]) bool {
		pluginlist[i] = plugininfo{
			name:   manager.Service,
			brief:  manager.Options.Brief,
			status: manager.IsEnabledIn(gid),
		}
		return true
	})
	logo, err := gg.LoadImage(kanbanpath + "zbpuwu.png")
	if err != nil {
		return
	}
	serverlistlogo, err := renderserverlistlogo()
	if err != nil {
		return
	}
	max := len(pluginlist)
	ln := math.Ceil(max, 3)
	w := (290+24)*3 + 24
	h := serverlistlogo.Bounds().Dy() + ln*(80+16) + serverlistlogo.Bounds().Dy()/3
	canvas := gg.NewContext(w, h)

	canvas.SetRGBA255(235, 235, 235, 255)
	canvas.Clear()

	canvas.SetRGBA255(135, 144, 173, 255)
	canvas.NewSubPath()
	canvas.MoveTo(0, 0)
	canvas.LineTo(float64(canvas.W()), 140)
	canvas.LineTo(float64(canvas.W()), 0)
	canvas.ClosePath()
	canvas.Fill()

	canvas.NewSubPath()
	canvas.MoveTo(float64(canvas.W()), float64(canvas.H()))
	canvas.LineTo(0, float64(canvas.H()))
	canvas.LineTo(0, float64(canvas.H())-140)
	canvas.ClosePath()
	canvas.Fill()

	canvas.SetRGBA255(247, 177, 170, 255)
	canvas.NewSubPath()
	canvas.MoveTo(0, 0)
	canvas.LineTo(float64(canvas.W()), 0)
	canvas.LineTo(float64(canvas.W()), 70)
	canvas.LineTo(0, 270)
	canvas.ClosePath()
	canvas.Fill()

	canvas.NewSubPath()
	canvas.MoveTo(float64(canvas.W()), float64(canvas.H()))
	canvas.LineTo(0, float64(canvas.H()))
	canvas.LineTo(0, float64(canvas.H())-70)
	canvas.LineTo(float64(canvas.W()), float64(canvas.H())-270)
	canvas.ClosePath()
	canvas.Fill()

	canvas.SetRGBA255(186, 113, 132, 255)
	canvas.NewSubPath()
	canvas.MoveTo(0, 0)
	canvas.LineTo(float64(canvas.W()), 0)
	canvas.LineTo(float64(canvas.W()), 35)
	canvas.LineTo(0, 160)
	canvas.ClosePath()
	canvas.Fill()

	canvas.NewSubPath()
	canvas.MoveTo(float64(canvas.W()), float64(canvas.H()))
	canvas.LineTo(0, float64(canvas.H()))
	canvas.LineTo(0, float64(canvas.H())-35)
	canvas.LineTo(float64(canvas.W()), float64(canvas.H())-160)
	canvas.ClosePath()
	canvas.Fill()

	canvas.ScaleAbout(0.5, 0.5, float64(canvas.W()/4), 0)
	canvas.DrawImageAnchored(logo, canvas.W()/4, 0, 0.5, 0)
	canvas.Identity()

	canvas.DrawImageAnchored(serverlistlogo, canvas.W()/4*3, 0, 0.5, 0)

	cardimgs := make([]image.Image, 4)

	wg := sync.WaitGroup{}
	cardsnum := math.Ceil(ln, 4) * 3
	wg.Add(4)
	for i := 0; i < 3; i++ {
		a := i * cardsnum
		b := (i + 1) * cardsnum
		go func(i int, list []plugininfo) {
			defer wg.Done()
			cardimgs[i], err = renderinfocards(list)
			if err != nil {
				return
			}
		}(i, pluginlist[a:b])
	}
	go func() {
		defer wg.Done()
		cardimgs[3], err = renderinfocards(pluginlist[cardsnum*3:])
		if err != nil {
			return
		}
	}()
	wg.Wait()
	spacing := 0
	for i := 0; i < len(cardimgs); i++ {
		canvas.DrawImage(cardimgs[i], 0, serverlistlogo.Bounds().Dy()+spacing)
		spacing += cardimgs[i].Bounds().Dy()
	}

	img = canvas.Image()
	return
}

func renderinfocards(plugininfos []plugininfo) (img image.Image, err error) {
	w := (290+24)*3 + 24
	cardnum := len(plugininfos)
	h := math.Ceil(cardnum, 3) * (80 + 16)
	cardw, cardh := 290.0, 80.0
	spacingw, spacingh := 24.0, 16.0
	canvas := gg.NewContext(w, h)
	beginw, beginh := 24.0, 0.0
	for i := 0; i < cardnum; i++ {
		canvas.SetRGBA255(204, 51, 51, 255)
		if plugininfos[i].status {
			canvas.SetRGBA255(136, 178, 0, 255)
		}
		canvas.DrawRoundedRectangle(beginw, beginh, cardw/2, cardh, 16)
		canvas.Fill()

		canvas.SetRGBA255(34, 26, 33, 255)
		canvas.DrawRoundedRectangle(beginw+10, beginh, cardw-10, cardh, 16)
		canvas.Fill()
		beginw += cardw + spacingw
		if (i+1)%3 == 0 {
			beginw = spacingw
			beginh += cardh + spacingh
		}
	}

	err = canvas.ParseFontFace(torussd, 36)
	if err != nil {
		return
	}
	canvas.SetRGBA255(235, 235, 235, 255)
	beginw, beginh = 24.0, 0.0
	for i := 0; i < cardnum; i++ {
		canvas.DrawStringAnchored(plugininfos[i].name, beginw+14, beginh+canvas.FontHeight()/2+4, 0, 0.5)
		beginw += cardw + spacingw
		if (i+1)%3 == 0 {
			beginw = spacingw
			beginh += cardh + spacingh
		}
	}
	err = canvas.ParseFontFace(glowsd, 16)
	if err != nil {
		return
	}
	beginw, beginh = 24.0, 0.0
	for i := 0; i < cardnum; i++ {
		canvas.DrawStringAnchored(plugininfos[i].brief, beginw+14, beginh+cardh-canvas.FontHeight()-4, 0, 0.5)
		beginw += cardw + spacingw
		if (i+1)%3 == 0 {
			beginw = spacingw
			beginh += cardh + spacingh
		}
	}
	img = canvas.Image()
	return
}

func renderserverlistlogo() (img image.Image, err error) {
	const w, h = 400, 200
	canvas := gg.NewContext(w, h)
	canvas.SetRGBA255(187, 122, 132, 255)
	err = canvas.ParseFontFace(yurukasd, 72)
	if err != nil {
		return
	}
	canvas.DrawStringAnchored("Server", 22, 45, 0, 1)
	canvas.SetRGBA255(246, 166, 171, 255)

	canvas.DrawStringAnchored("List", 219, 112, 0, 1)

	canvas.SetRGBA255(135, 145, 173, 255)
	err = canvas.ParseFontFace(yurukasd, 48)
	if err != nil {
		return
	}
	canvas.DrawStringAnchored("服務列表", 23, 120, 0, 1)

	canvas.SetRGBA255(103, 178, 240, 255)
	err = canvas.ParseFontFace(yurukasd, 36)
	if err != nil {
		return
	}

	canvas.DrawStringAnchored("サーバー", 82, 25, 0, 1)
	canvas.DrawStringAnchored("リスト", 280, 83, 0, 1)

	mask := canvas.AsMask()

	stroked := gg.NewContext(w, h)
	err = stroked.SetMask(mask)
	if err != nil {
		return
	}
	stroked.SetRGBA255(255, 255, 255, 255)
	stroked.DrawRectangle(0, 0, float64(stroked.W()), float64(stroked.H()))
	stroked.Fill()

	strokedimg := stroked.Image()
	coloredimg := canvas.Image()

	canvas = gg.NewContext(w, h)

	canvas.DrawImage(strokedimg, -3, 0)
	canvas.DrawImage(strokedimg, 3, 0)
	canvas.DrawImage(strokedimg, 0, -3)
	canvas.DrawImage(strokedimg, 0, 3)
	canvas.DrawImage(strokedimg, -3, -3)
	canvas.DrawImage(strokedimg, 3, 3)
	canvas.DrawImage(strokedimg, 3, -3)
	canvas.DrawImage(strokedimg, -3, 3)
	canvas.DrawImage(coloredimg, 0, 0)

	img = canvas.Image()
	return
}
