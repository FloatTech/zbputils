package control

import (
	"image"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/Coloured-glaze/gg"
	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/img/writer"
	"github.com/FloatTech/floatbox/math"
	"github.com/FloatTech/floatbox/process"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/disintegration/imaging"
	zero "github.com/wdvxdr1123/ZeroBot"

	"github.com/FloatTech/rendercard"

	"github.com/FloatTech/zbputils/img/text"
)

const (
	bannerpath = "data/zbpbanner/"
	kanbanpath = "data/Control/"
	bannerurl  = "https://gitcode.net/u011570312/zbpbanner/-/raw/main/"
)

type plugininfo struct {
	name   string
	brief  string
	banner string
	status bool
}

var (
	// 底图缓存
	imgtmp image.Image
	// lnperpg 每页行数
	lnperpg = 9
)

func init() {
	err := os.MkdirAll(bannerpath, 0755)
	if err != nil {
		panic(err)
	}
	_, err = file.GetLazyData(kanbanpath+"kanban.png", Md5File, true)
	if err != nil {
		panic(err)
	}
}

func drawservicesof(gid int64) (imgs [][]byte, err error) {
	limit := make(chan struct{}, runtime.NumCPU())
	pluginlen := len(priomap)
	pluginlist := make([]*plugininfo, pluginlen)
	ForEachByPrio(func(i int, manager *ctrl.Control[*zero.Ctx]) bool {
		pluginlist[i] = &plugininfo{
			name:   manager.Service,
			brief:  manager.Options.Brief,
			banner: manager.Options.Banner,
			status: manager.IsEnabledIn(gid),
		}
		return true
	})
	cardnum := lnperpg * 3
	// 分页
	if pluginlen < cardnum {
		// 如果单页显示数量超出了总数量
		lnperpg = math.Ceil(pluginlen, 3)
	}
	page := math.Ceil(pluginlen, cardnum)
	imgs = make([][]byte, page)
	if imgtmp == nil {
		imgtmp, err = (&rendercard.Title{
			Line:          lnperpg,
			LeftTitle:     "服务列表",
			LeftSubtitle:  "service_list",
			RightTitle:    "FloatTech",
			RightSubtitle: "ZeroBot-Plugin",
			TitleFont:     text.GlowSansFontFile,
			TextFont:      text.ImpactFontFile,
			ImagePath:     kanbanpath + "kanban.png",
		}).DrawTitle()
		if err != nil {
			return
		}
	}
	var wg, cwg, swg sync.WaitGroup
	cardlist := make([]image.Image, pluginlen)
	wg.Add(page)
	cwg.Add(pluginlen)
	for k, info := range pluginlist {
		go func(k int, info *plugininfo) {
			defer cwg.Done()

			limit <- struct{}{}
			defer func() { <-limit }()

			banner := ""
			switch {
			case strings.HasPrefix(info.banner, "http"):
				err = file.DownloadTo(info.banner, bannerpath+info.name+".png")
				if err != nil {
					return
				}
				process.SleepAbout1sTo2s()
				banner = bannerpath + info.name + ".png"
			case info.banner != "":
				banner = info.banner
			default:
				_, err = file.GetCustomLazyData(bannerurl, bannerpath+info.name+".png")
				if err == nil {
					banner = bannerpath + info.name + ".png"
				}
			}
			cardlist[k], err = (&rendercard.Title{
				IsEnabled:    info.status,
				LeftTitle:    info.name,
				LeftSubtitle: info.brief,
				ImagePath:    banner,
				TitleFont:    text.ImpactFontFile,
				TextFont:     text.GlowSansFontFile,
			}).DrawCard()
			if err != nil {
				return
			}
		}(k, info)
	}
	cwg.Wait()
	for l := 0; l < page; l++ { // 页数
		swg.Add(1)
		var shadowimg image.Image
		go func() {
			defer swg.Done()

			limit <- struct{}{}
			defer func() { <-limit }()

			x, y := 30+2, 30+300+30+6+4
			shadow := gg.NewContextForImage(rendercard.Transparency(imgtmp, 0))
			shadow.SetRGBA255(0, 0, 0, 192)
			for i := 0; i < math.Min(cardnum, pluginlen-cardnum*l); i++ {
				shadow.DrawRoundedRectangle(float64(x), float64(y), 384-4, 256-4, 0)
				shadow.Fill()
				x += 384 + 30
				if (i+1)%3 == 0 {
					x = 30 + 2
					y += 256 + 30
				}
			}
			shadowimg = shadow.Image()
		}()
		swg.Wait()
		go func(l int) {
			defer wg.Done()

			limit <- struct{}{}
			defer func() { <-limit }()

			one := gg.NewContextForImage(imgtmp)
			x, y := 30, 30+300+30
			one.DrawImage(imaging.Blur(shadowimg, 7), 0, 0)
			for i := 0; i < math.Min(cardnum, pluginlen-cardnum*l); i++ {
				one.DrawImage(rendercard.Fillet(cardlist[(cardnum*l)+i], 8), x, y)
				x += 384 + 30
				if (i+1)%3 == 0 {
					x = 30
					y += 256 + 30
				}
			}
			data, cl := writer.ToBytes(one.Image()) // 生成图片
			imgs[l] = data
			cl()
		}(l)
	}
	wg.Wait()
	return
}

// 获取字体和头像
func geticonandfont() (err error) {
	_, err = file.GetLazyData(text.BoldFontFile, Md5File, true)
	if err != nil {
		return
	}
	_, err = file.GetLazyData(text.SakuraFontFile, Md5File, true)
	if err != nil {
		return
	}
	_, err = file.GetLazyData(kanbanpath+"kanban.png", Md5File, true)
	if err != nil {
		return
	}
	return
}
