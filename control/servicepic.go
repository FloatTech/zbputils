package control

import (
	"errors"
	"image"
	"os"
	"strconv"
	"strings"

	"github.com/Coloured-glaze/gg"
	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/img/writer"
	ctrl "github.com/FloatTech/zbpctrl"
	zero "github.com/wdvxdr1123/ZeroBot"

	"github.com/FloatTech/rendercard"

	"github.com/FloatTech/zbputils/img/text"
)

const (
	PAGECNT    = 27
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

// 底图缓存
var imgtmp image.Image

func init() {
	err := os.MkdirAll(bannerpath, 0755)
	if err != nil {
		panic(err)
	}
	_, err = file.GetLazyData(kanbanpath+"icon.jpg", true)
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
	page := len(plist) / PAGECNT
	if len(plist)%PAGECNT != 0 {
		page++
	}
	imgs = make([][]byte, 0, page)
	if imgtmp == nil {
		imgtmp, err = rendercard.Titleinfo{
			Lefttitle:     "服务列表",
			Leftsubtitle:  "service_list",
			Righttitle:    "FloatTech",
			Rightsubtitle: "ZeroBot-Plugin",
			Textpath:      text.SakuraFontFile,
			Imgpath:       kanbanpath + "icon.jpg",
		}.Drawtitle()
		if err != nil {
			return
		}
	}
	var card image.Image
	for l := 0; l < page; l++ {
		one := gg.NewContextForImage(imgtmp)
		x, y := 30, 30
		for j := 0; j < 9; j++ {
			for i := 0; i < 3; i++ {
				if k == len(plist) {
					break
				}
				kstr := strconv.Itoa(k)
				var banner string
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
					if err != nil {
						err = errors.New("ERROR: 插件背景图下载失败或是自定义插件")
						return
					}
					banner = bannerpath + plist[k].name
				}
				card, err = rendercard.Titleinfo{
					Lefttitle:     plist[k].name,
					Leftsubtitle:  plist[k].brief,
					Rightsubtitle: kstr,
					Imgpath:       banner,
					Textpath:      text.SakuraFontFile,
					Textpath2:     text.BoldFontFile,
					Status:        plist[k].status,
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
