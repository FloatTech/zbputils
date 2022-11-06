package control

import (
	"encoding/json"
	"image"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/Coloured-glaze/gg"
	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/img/writer"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/img"
	"github.com/FloatTech/zbputils/img/text"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const bannerpath = "zbpbanner/"

type plugininfo struct {
	PluginName   string
	PluginBrief  string
	PluginStatus bool
}
type stylecfg struct {
	style int // 服务列表样式
}

var style *stylecfg

func init() {
	_ = os.MkdirAll(bannerpath, 0755)
	if file.IsExist(bannerpath + "config.json") {
		reader, err := os.Open(bannerpath + "config.json")
		if err != nil {
			panic(err)
		}
		err = json.NewDecoder(reader).Decode(&style)
		if err != nil {
			panic(err)
		}
	} else {
		style = &stylecfg{style: 0}
		writefile, err := os.Create(bannerpath + "config.json")
		if err != nil {
			panic(err)
		}
		err = json.NewEncoder(writefile).Encode(&style)
		if err != nil {
			panic(err)
		}
	}
}

func renderimg(ctx *zero.Ctx) (err error) {
	gid := ctx.Event.GroupID
	if gid == 0 {
		gid = -ctx.Event.UserID
	}
	var plist = make([]*plugininfo, 0, len(priomap))
	ForEachByPrio(func(i int, manager *ctrl.Control[*zero.Ctx]) bool {
		plist = append(plist, &plugininfo{
			PluginName:   manager.Service,
			PluginBrief:  manager.Options.Brief,
			PluginStatus: manager.IsEnabledIn(gid),
		})
		return true
	})
	var k int
	// 分页
	page := len(plist) / 27
	if page%27 == 0 {
		page -= 1
	}
	servicelist := make(message.Message, 0, page)
	var imgtmp image.Image

	imgw := 1272.0
	// 创建图像
	canvas := gg.NewContext(int(imgw), 30+30+300+(9*(256+30)))
	canvas.SetRGBA255(240, 240, 240, 255)
	canvas.Clear()

	// 标题背景1
	canvas.DrawRectangle(0, 30, imgw, 300)
	canvas.SetRGBA255(0, 0, 0, 153)
	canvas.Fill()

	// 标题背景2
	canvas.DrawRectangle(0, 30+40, imgw, 220)
	canvas.SetRGBA255(0, 0, 0, 153)
	canvas.Fill()

	// 加载size为108的字体
	err = canvas.LoadFontFace(text.SakuraFontFile, 108)
	if err != nil {
		return
	}

	// 绘制标题
	canvas.SetRGBA255(240, 240, 240, 255)
	canvas.DrawString("服务列表", 25, 30+40+55+canvas.FontHeight()-canvas.FontHeight()/3)

	// 加载size为54的字体
	err = canvas.LoadFontFace(text.SakuraFontFile, 54)
	if err != nil {
		return
	}

	// 绘制一系列标题
	canvas.DrawString("service_list", 25+3, 30+40+165+canvas.FontHeight()/3)

	fw, _ := canvas.MeasureString("FloatTech")
	canvas.DrawString("FloatTech", imgw-25-fw-170-25, 30+40+25+15+canvas.FontHeight()+canvas.FontHeight()/4)
	fw1, _ := canvas.MeasureString("ZeroBot-Plugin")
	canvas.DrawString("ZeroBot-Plugin", imgw-25-fw1-170-25, 30+40+25+15+canvas.FontHeight()*2+canvas.FontHeight()/2)
	canvas.SetRGBA255(240, 240, 240, 255)

	// 加载icon并绘制
	var icon *img.Factory
	icon, err = img.LoadFirstFrame(kanbanPath+"icon.jpg", 170, 170)
	if err != nil {
		return
	}
	canvas.DrawImage(icon.Im, int(imgw)-25-170, 30+40+25)
	imgtmp = canvas.Image()
	for l := 0; l <= page; l++ {
		one := gg.NewContextForImage(imgtmp)
		x, y := 30.0, 30.0
		for j := 0; j < 9; j++ {
			for i := 0; i < 3; i++ {
				if k == len(plist) {
					break
				}
				err = drawplugin(one, x, y, k+1, plist[k])
				if err != nil {
					return
				}
				k++
				x += 384 + 30
			}
			x = 30.0
			y += 256 + 30
		}
		data, cl := writer.ToBytes(one.Image()) // 生成图片
		servicelist = append(servicelist, ctxext.FakeSenderForwardNode(ctx, message.ImageBytes(data)))
		cl()
	}
	if id := ctx.Send(servicelist); id.ID() == 0 {
		ctx.SendChain(message.Text("ERROR: 可能被风控了"))
	}
	return nil
}

func drawplugin(canvas *gg.Context, x, y float64, i int, list *plugininfo) (err error) {
	var impng *img.Factory
	// 绘制图片
	if strings.HasPrefix(list.PluginBrief, "http://") || strings.HasPrefix(list.PluginBrief, "https://") {
		if file.IsNotExist(bannerpath + "network/" + list.PluginName + ".png") {
			_ = os.MkdirAll(bannerpath+"network/", 0755)
			// 地址无效也继续
			err = file.DownloadTo(list.PluginBrief, bannerpath+"network/"+list.PluginName+".png", true)
			if err != nil {
				logrus.Warn("[control] ERROR: " + list.PluginName + "图片地址无效")
			}
		}
		impng, err = img.LoadFirstFrame(bannerpath+"network/"+list.PluginName+".png", 768, 512)
	} else {
		// 下载失败也不影响运行
		_, _ = file.GetLazyData(bannerpath+list.PluginName+".png", true)
		impng, err = img.LoadFirstFrame(bannerpath+list.PluginName+".png", 768, 512)
	}
	recw, rech := 384.0, 256.0
	if err == nil {
		canvas.DrawImage(img.Size(impng.Im, int(recw), int(rech)).Im, int(x), int(y)+300+30)
	} else {
		canvas.DrawRectangle(x, y+300+30, recw, rech)
		canvas.SetRGBA255(rand.Intn(195)+15, rand.Intn(195)+15, rand.Intn(195)+15, 255)
		canvas.Fill()
	}

	// 绘制遮罩
	canvas.DrawRectangle(x, y+300+30+(rech/3*2), recw, rech/3)
	canvas.SetRGBA255(0, 0, 0, 153)
	canvas.Fill()

	// 绘制排名
	canvas.DrawRectangle(x+recw/10, y+300+30, recw/10, (rech/4)-10)
	canvas.DrawRoundedRectangle(x+recw/10, y+300+30, recw/10, (rech / 4), 8)
	if list.PluginStatus {
		canvas.SetRGBA255(15, 175, 15, 255)
	} else {
		canvas.SetRGBA255(200, 15, 15, 255)
	}
	canvas.Fill()

	// 绘制插件排名
	canvas.SetRGBA255(15, 15, 15, 255)
	var fw2 float64
	if i > 99 {
		canvas.LoadFontFace(text.SakuraFontFile, 24)
		fw2, _ = canvas.MeasureString(strconv.FormatInt(int64(i), 10))
		canvas.DrawString(strconv.FormatInt(int64(i), 10), x+recw/10+((recw/10-fw2)/2), y+300+30+canvas.FontHeight()*3/8+(rech/8))
	} else {
		canvas.LoadFontFace(text.SakuraFontFile, 28)
		fw2, _ = canvas.MeasureString(strconv.FormatInt(int64(i), 10))
		canvas.DrawString(strconv.FormatInt(int64(i), 10), x+recw/10+((recw/10-fw2)/2), y+300+30+canvas.FontHeight()*3/8+(rech/8))

	}

	// 绘制插件信息
	canvas.SetRGBA255(240, 240, 240, 255)
	err = canvas.LoadFontFace(text.SakuraFontFile, 48)
	if err != nil {
		return
	}
	canvas.DrawString(list.PluginName, x+recw/32, y+300+30+(recw*0.475)+canvas.FontHeight()-canvas.FontHeight()/4)

	err = canvas.LoadFontFace(text.SakuraFontFile, 24)
	if err != nil {
		return
	}
	canvas.DrawString(list.PluginName, x+recw/32, y+300+30+(recw*0.475)+recw/6-canvas.FontHeight()/4)
	return nil
}

func renderusage(ctx *zero.Ctx, s *ctrl.Control[*zero.Ctx], gid int64) (err error) {
	// 图像宽
	imgw := 1272.0

	// 处理插件帮助并且计算图像高
	plugininfo := strings.Split(strings.Trim(s.String(), "\n"), "\n")
	newplugininfo := make([]string, 0, len(plugininfo)*2)
	font := gg.NewContext(1, 1)
	err = font.LoadFontFace(text.BoldFontFile, 38)
	if err != nil {
		return
	}
	for i := 0; i < len(plugininfo); i++ {
		newlinetext, textw, tmpw := "", 0.0, 0.0
		for len(plugininfo[i]) > 0 {
			newlinetext, tmpw = truncate(font, plugininfo[i], imgw-50)
			newplugininfo = append(newplugininfo, newlinetext)
			if tmpw > textw {
				textw = tmpw
			}
			if len(newlinetext) >= len(plugininfo[i]) {
				break
			}
			plugininfo[i] = plugininfo[i][len(newlinetext):]
		}
	}
	imgh := len(newplugininfo)*(int(font.FontHeight())+20) + 220 + 10 + 30 + 10 + 50

	// 创建图像
	canvas := gg.NewContext(int(imgw), imgh)
	canvas.SetRGBA255(15, 15, 15, 204)
	canvas.Clear()

	// 加载icon
	var icon *img.Factory
	icon, err = img.LoadFirstFrame(kanbanPath+"icon.jpg", 170, 170)
	if err != nil {
		return
	}
	canvas.DrawImage(icon.Im, int(imgw)-25-170, 25)

	// 绘制标题与内容的分割线
	canvas.DrawRectangle(0, 220, float64(imgw), 10)
	canvas.SetRGBA255(240, 240, 240, 255)
	canvas.Fill()

	// 加载size为108的字体
	err = canvas.LoadFontFace(text.SakuraFontFile, 108)
	if err != nil {
		return
	}

	// 绘制标题
	canvas.SetRGBA255(240, 240, 240, 255)
	canvas.DrawString(s.Service, 25+40+25, 55+canvas.FontHeight()-canvas.FontHeight()/3)

	// 加载size为54的字体
	err = canvas.LoadFontFace(text.SakuraFontFile, 54)
	if err != nil {
		return
	}

	// 绘制插件开启状态
	canvas.DrawRectangle(25, 25, 40, 170)
	if s.IsEnabledIn(gid) {
		canvas.SetRGBA255(15, 175, 15, 255)
	} else {
		canvas.SetRGBA255(200, 15, 15, 255)
	}
	canvas.Fill()
	canvas.SetRGBA255(240, 240, 240, 255)

	// 绘制一系列标题
	cnname := s.Options.Brief
	canvas.DrawString(cnname, 25+3+40+25, 165+canvas.FontHeight()/3)
	fw, _ := canvas.MeasureString("FloatTech")
	canvas.DrawString("FloatTech", imgw-25-fw-170-25, 25+15+canvas.FontHeight()+canvas.FontHeight()/4)
	fw1, _ := canvas.MeasureString("ZeroBot-Plugin")
	canvas.DrawString("ZeroBot-Plugin", imgw-25-fw1-170-25, 25+15+canvas.FontHeight()*2+canvas.FontHeight()/2)

	// 加载size为42的字体
	err = canvas.LoadFontFace(text.BoldFontFile, 38)
	if err != nil {
		return
	}

	x, y := 25.0, 25.0
	for i := 0; i < len(newplugininfo); i++ {
		canvas.DrawString(newplugininfo[i], x, y+220+10+canvas.FontHeight())
		y += 20 + canvas.FontHeight()
	}
	data, cl := writer.ToBytes(canvas.Image())
	ctx.SendChain(message.ImageBytes(data))
	cl()
	return
}
