package controller

import (
	"encoding/json"
	"strconv"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
)

// ManageController 主控制类
type ManageController struct {
	BaseController
}

// GetBotList
// @Description 获取机器人qq号
// @Router /api/getBotList [post]
func (c *ManageController) GetBotList(context *gin.Context) {
	var bots []int64

	zero.RangeBot(func(id int64, ctx *zero.Ctx) bool {
		bots = append(bots, id)
		return true
	})
	c.OkWithData(bots, context)
}

// GetGroupList
// @Description 获取群列表
// @Router /api/getGroupList [post]
// @Param self_id formData integer true "机器人qq号" default(123456)
func (c *ManageController) GetGroupList(context *gin.Context) {
	selfID, err := strconv.Atoi(context.PostForm("self_id"))
	if err != nil {
		var data map[string]interface{}
		err := context.BindJSON(&data)
		if err != nil {
			log.Errorln("[gui]" + err.Error())
			return
		}
		selfID = int(data["self_id"].(float64))
	}

	bot := zero.GetBot(int64(selfID))
	var resp []interface{}
	list := bot.GetGroupList().String()
	err = json.Unmarshal([]byte(list), &resp)
	if err != nil {
		log.Errorln("[gui]" + err.Error())
	}
	c.OkWithData(resp, context)
}

// GetPluginList
// @Description 获取所有插件的状态
// @Router /api/getPluginList [post]
// @Param group_id formData integer true "群号" default(0)
func (c *ManageController) GetPluginList(context *gin.Context) {
	groupID, err := strconv.ParseInt(context.PostForm("group_id"), 10, 64)
	if err != nil {
		var parse map[string]interface{}
		err := context.BindJSON(&parse)
		if err != nil {
			log.Errorln("[gui]" + err.Error())
			return
		}
		groupID = int64(parse["group_id"].(float64))
	}
	var datas []map[string]interface{}
	control.ForEachByPrio(func(i int, manager *ctrl.Control[*zero.Ctx]) bool {
		datas = append(datas, map[string]interface{}{"id": i, "name": manager.Service, "brief": manager.Options.Brief, "usage": manager.String(), "banner": "https://gitcode.net/u011570312/zbpbanner/-/raw/main/" + manager.Service + ".png", "status": manager.IsEnabledIn(groupID)})
		return true
	})
	c.OkWithData(datas, context)
}

// UpdatePluginStatus
// @Description 更改某一个插件状态
// @Router /api/updatePluginStatus [post]
// @Param group_id formData integer true "群号" default(0)
// @Param name formData string true "插件名" default(aireply)
// @Param status formData boolean true "插件状态" default(true)
func (c *ManageController) UpdatePluginStatus(context *gin.Context) {
	var parse map[string]interface{}
	err := context.BindJSON(&parse)
	if err != nil {
		log.Errorln("[gui] ", err)
		return
	}
	groupID := int64(parse["group_id"].(float64))
	name := parse["name"].(string)
	enable := parse["status"].(bool)
	con, b := control.Lookup(name)
	if !b {
		context.JSON(404, "服务不存在")
		return
	}
	if enable {
		con.Enable(groupID)
	} else {
		con.Disable(groupID)
	}
	c.Ok(context)
}
