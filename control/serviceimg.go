package control

import (
	"strconv"
	"strings"

	"github.com/Coloured-glaze/gg"
	"github.com/FloatTech/floatbox/img/writer"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/img"
	"github.com/FloatTech/zbputils/img/text"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type plugininfo struct {
	PluginENName string
	PluginCNName string
	PluginStatus bool
}

func renderimg(ctx *zero.Ctx) (err error) {
	gid := ctx.Event.GroupID
	if gid == 0 {
		gid = -ctx.Event.UserID
	}
	var plist = make([]*plugininfo, 0, len(priomap))
	ForEachByPrio(func(i int, manager *ctrl.Control[*zero.Ctx]) bool {
		cnname := strings.Split(manager.Options.Help, "\n")[0]
		var status bool
		if manager.IsEnabledIn(gid) {
			status = true
		} else {
			status = false
		}
		plist = append(plist, &plugininfo{
			PluginENName: manager.Service,
			PluginCNName: cnname,
			PluginStatus: status,
		})
		return true
	})
	imgw := 2400.0
	// 创建图像
	canvas := gg.NewContext(int(imgw), 75+75+413+25+(8*290))
	canvas.SetRGBA255(234, 234, 234, 255)
	canvas.Clear()

	// 标题背景1
	canvas.DrawRectangle(0, 75, imgw, 413)
	canvas.SetRGBA255(0, 0, 0, 153)
	canvas.Fill()

	// 标题背景2
	canvas.DrawRectangle(0, 131, imgw, 300)
	canvas.SetRGBA255(0, 0, 0, 153)
	canvas.Fill()

	// 加载size为144的字体
	err = canvas.LoadFontFace(text.SakuraFontFile, 144)
	if err != nil {
		return
	}

	// 绘制标题
	canvas.SetRGBA255(240, 240, 240, 255)
	canvas.DrawString("服务列表", 45, 131+canvas.FontHeight()-canvas.FontHeight()/8+53)

	// 加载size为72的字体
	err = canvas.LoadFontFace(text.SakuraFontFile, 72)
	if err != nil {
		return
	}

	// 绘制一系列标题
	canvas.DrawString("service_list", 45, 75+356-canvas.FontHeight()/4-53)
	fw, _ := canvas.MeasureString("FloatTech")
	canvas.DrawString("FloatTech", imgw-45-fw-188-15, 131+canvas.FontHeight()-canvas.FontHeight()/4+75+8)
	fw1, _ := canvas.MeasureString("ZeroBot-Plugin")
	canvas.DrawString("ZeroBot-Plugin", imgw-45-fw1-188-15, 75+356-canvas.FontHeight()/4-75-8)
	canvas.SetRGBA255(240, 240, 240, 255)

	// 加载icon
	icon, err := img.LoadFirstFrame(kanbanPath+"icon.jpg", 188, 188)
	if err != nil {
		return
	}
	canvas.DrawImage(icon.Im, int(imgw)-45-188, 131+56)

	x, y := 75.0, 75.0
	var k int
	for j := 0; j < len(plist)/5; j++ {
		for i := 0; i < 5; i++ {

			err = drawplugin(canvas, x, y, k+1, plist[k])
			if err != nil {
				return
			}
			k++
			x += 460
		}
		x = 75.0
		y += 290
	}
	data, cl := writer.ToBytes(canvas.Image()) // 生成图片
	if id := ctx.SendChain(message.ImageBytes(data)); id.ID() == 0 {
		ctx.SendChain(message.Text("ERROR: 可能被风控了"))
	}
	cl()
	return nil
}

func drawplugin(canvas *gg.Context, x, y float64, i int, list *plugininfo) (err error) {
	// canvas.DrawRectangle(x, y+413+75, 390, 240)
	// 绘制图片
	imjpg, err := img.LoadFirstFrame(kanbanPath+"input.jpg", 1920, 1080)
	if err == nil {
		canvas.DrawImage(img.Size(imjpg.Im, 390, 240).Im, int(x), int(y)+413+75)
	} else {
		canvas.DrawRectangle(x, y+413+75, 390, 240)
		canvas.SetRGBA255(15, 15, 200, 255)
		canvas.Fill()
	}

	/*canvas.DrawRoundedRectangle(x, y+413+75, 390, 240, 16)

	canvas.SetRGBA255(0, 0, 240, 200)
	canvas.Fill()*/

	canvas.DrawRectangle(x, y+413+75+140, 390, 100)
	// 绘制遮罩
	// canvas.DrawRoundedRectangle(x, y+413+75+140, 390, 100, 16)
	canvas.SetRGBA255(0, 0, 0, 153)
	canvas.Fill()

	// 绘制排名
	canvas.DrawRectangle(x+30, y+413+75, 40, 50)
	canvas.DrawRoundedRectangle(x+30, y+413+75, 40, 60, 8)
	if list.PluginStatus {
		canvas.SetRGBA255(30, 195, 30, 255)
	} else {
		canvas.SetRGBA255(215, 30, 30, 255)
	}
	canvas.Fill()
	canvas.SetRGBA255(15, 15, 15, 255)

	var fw2 float64
	if i > 99 {
		err = canvas.LoadFontFace(text.SakuraFontFile, 26)
		if err != nil {
			return
		}
		fw2, _ = canvas.MeasureString(strconv.FormatInt(int64(i), 10))
		canvas.DrawString(strconv.FormatInt(int64(i), 10), x+30+((40-fw2)/2), y+413+75+30+canvas.FontHeight()*3/8)
	} else {
		err = canvas.LoadFontFace(text.SakuraFontFile, 32)
		if err != nil {
			return
		}
		fw2, _ = canvas.MeasureString(strconv.FormatInt(int64(i), 10))
		canvas.DrawString(strconv.FormatInt(int64(i), 10), x+30+((40-fw2)/2), y+413+75+30+canvas.FontHeight()*3/8)

	}

	canvas.SetRGBA255(240, 240, 240, 255)
	err = canvas.LoadFontFace(text.SakuraFontFile, 48)
	if err != nil {
		return
	}
	canvas.DrawString(list.PluginENName, x+20, y+413+75+140+canvas.FontHeight()*6/5)

	err = canvas.LoadFontFace(text.SakuraFontFile, 22)
	if err != nil {
		return
	}
	canvas.DrawString(list.PluginCNName, x+20+2, y+413+75+140+50+canvas.FontHeight()*8/5)
	return nil
}
