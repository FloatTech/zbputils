package control

import (
	"image"
	"os"
	"strings"

	"github.com/Coloured-glaze/gg"
	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/img/writer"
	"github.com/FloatTech/floatbox/math"
	ctrl "github.com/FloatTech/zbpctrl"
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
	_, err = file.GetLazyData(kanbanpath+"icon.jpg", Md5File, true)
	if err != nil {
		panic(err)
	}
}

func drawservicesof(gid int64) (imgs [][]byte, err error) {
	plist := make([]*plugininfo, len(priomap))
	ForEachByPrio(func(i int, manager *ctrl.Control[*zero.Ctx]) bool {
		plist[i] = &plugininfo{
			name:   manager.Service,
			brief:  manager.Options.Brief,
			banner: manager.Options.Banner,
			status: manager.IsEnabledIn(gid),
		}
		return true
	})
	k := 0
	// 分页
	if len(plist) < 3*lnperpg {
		// 如果单页显示数量超出了总数量
		lnperpg = math.Ceil(len(plist), 3)
	}
	page := math.Ceil(len(plist), 3*lnperpg)
	imgs = make([][]byte, page)
	if imgtmp == nil {
		imgtmp, err = rendercard.Titleinfo{
			Line:          lnperpg,
			Lefttitle:     "服务列表",
			Leftsubtitle:  "service_list",
			Righttitle:    "FloatTech",
			Rightsubtitle: "ZeroBot-Plugin",
			Fontpath:      text.SakuraFontFile,
			Imgpath:       kanbanpath + "icon.jpg",
		}.Drawtitle()
		if err != nil {
			return
		}
	}
	var card image.Image
	for l := 0; l < page; l++ { // 页数
		one := gg.NewContextForImage(imgtmp)
		x, y := 30, 30+300+30
		for j := 0; j < lnperpg; j++ { // 行数
			for i := 0; i < 3; i++ { // 列数
				if k == len(plist) {
					break
				}
				banner := ""
				switch {
				case strings.HasPrefix(plist[k].banner, "http"):
					err = file.DownloadTo(plist[k].banner, bannerpath+plist[k].name+".png", true)
					if err != nil {
						return
					}
					banner = bannerpath + plist[k].name + ".png"
				case plist[k].banner != "":
					banner = plist[k].banner
				default:
					_, err = file.GetCustomLazyData(bannerurl, bannerpath+plist[k].name+".png")
					if err == nil {
						banner = bannerpath + plist[k].name + ".png"
					}
				}
				card, err = rendercard.Titleinfo{
					Lefttitle:    plist[k].name,
					Leftsubtitle: plist[k].brief,
					Imgpath:      banner,
					Fontpath:     text.SakuraFontFile,
					Fontpath2:    text.BoldFontFile,
					Status:       plist[k].status,
				}.Drawcard()
				if err != nil {
					return
				}
				one.DrawImage(card, x, y)
				k++
				x += 384 + 30
			}
			x = 30
			y += 256 + 30
		}
		data, cl := writer.ToBytes(one.Image()) // 生成图片
		imgs[l] = data
		cl()
	}
	return
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
			break // 跳出
		} else {
			w = width
			res = append(res, r) // 写入
		}
	}
	return string(res), w
}
