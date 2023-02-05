package controller

import (
	"encoding/json"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/control/web/common"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
)

type GroupDto struct {
	SelfID int64 `json:"self_id" form:"self_id" validate:"required"`
}

type PluginDto struct {
	GroupID int64 `json:"group_id" form:"group_id"`
}

type PluginStatusDto struct {
	GroupID int64  `json:"group_id" form:"group_id"`
	Name    string `json:"name" form:"name" validate:"required"`
	Status  bool   `json:"status" form:"status"`
}

// GetBotList
// @Description 获取机器人qq号
// @Router /api/getBotList [get]
func GetBotList(context *gin.Context) {
	var bots []int64

	zero.RangeBot(func(id int64, ctx *zero.Ctx) bool {
		bots = append(bots, id)
		return true
	})
	common.OkWithData(bots, context)
}

// GetGroupList
// @Description 获取群列表
// @Router /api/getGroupList [get]
// @Param self_id query integer true "机器人qq号" default(123456)
func GetGroupList(context *gin.Context) {
	var d GroupDto
	err := common.Bind(&d, context)
	if err != nil {
		common.FailWithMessage(err.Error(), context)
		return
	}
	bot := zero.GetBot(d.SelfID)
	var resp []interface{}
	list := bot.GetGroupList().String()
	err = json.Unmarshal([]byte(list), &resp)
	if err != nil {
		log.Errorln("[gui]" + err.Error())
	}
	common.OkWithData(resp, context)
}

// GetPluginList
// @Description 获取所有插件的状态
// @Router /api/getPluginList [get]
// @Param group_id query integer false "群号" default(0)
func GetPluginList(context *gin.Context) {
	var d PluginDto
	err := common.Bind(&d, context)
	if err != nil {
		common.FailWithMessage(err.Error(), context)
		return
	}
	var datas []map[string]interface{}
	control.ForEachByPrio(func(i int, manager *ctrl.Control[*zero.Ctx]) bool {
		datas = append(datas, map[string]interface{}{"id": i, "name": manager.Service, "brief": manager.Options.Brief, "usage": manager.String(), "banner": "https://gitcode.net/u011570312/zbpbanner/-/raw/main/" + manager.Service + ".png", "status": manager.IsEnabledIn(d.GroupID)})
		return true
	})
	common.OkWithData(datas, context)
}

// UpdatePluginStatus
// @Description 更改某一个插件状态
// @Router /api/updatePluginStatus [post]
// @Param group_id formData integer false "群号" default(0)
// @Param name formData string true "插件名" default(aireply)
// @Param status formData boolean false "插件状态" default(true)
func UpdatePluginStatus(context *gin.Context) {
	var d PluginStatusDto
	err := common.Bind(&d, context)
	if err != nil {
		common.FailWithMessage(err.Error(), context)
		return
	}
	con, b := control.Lookup(d.Name)
	if !b {
		common.FailWithMessage(d.Name+"服务不存在", context)
		return
	}
	if d.Status {
		con.Enable(d.GroupID)
	} else {
		con.Disable(d.GroupID)
	}
	common.Ok(context)
}
