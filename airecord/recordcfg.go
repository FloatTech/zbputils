package aichat

import (
	"fmt"
	"strings"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
)

var (
	recordcfg = newrecordconfig()
)

// recordconfig 存储语音记录相关配置
type recordconfig struct {
	ModelName string // 语音模型名称
	ModelID   string // 语音模型ID
	Customgid int64  // 自定义群ID
}

// newrecordconfig 创建并返回默认语音记录配置
func newrecordconfig() recordconfig {
	return recordconfig{}
}

// getRecordConfig 返回当前语音记录配置信息
func getRecordConfig() recordconfig {
	return recordcfg
}

// setRecordModel 设置语音记录模型
func setRecordModel(modelName, modelID string) {
	recordcfg.ModelName = modelName
	recordcfg.ModelID = modelID
}

// setCustomGID 设置自定义群ID
func setCustomGID(gid int64) {
	recordcfg.Customgid = gid
}

// isvalid 检查语音记录配置是否有效
func (c *recordconfig) isvalid() bool {
	return c.ModelName != "" && c.ModelID != "" && c.Customgid != 0
}

// ensureRecordConfig 确保语音记录配置存在
func ensureRecordConfig(ctx *zero.Ctx) bool {
	c, ok := ctx.State["manager"].(*ctrl.Control[*zero.Ctx])
	if !ok {
		return false
	}
	if !recordcfg.isvalid() {
		err := c.GetExtra(&recordcfg)
		if err != nil {
			logrus.Warnln("ERROR: get extra err:", err)
		}
	}
	return true
}

// printRecordConfig 生成格式化的语音记录配置信息字符串
func printRecordConfig(recCfg recordconfig) string {
	var builder strings.Builder
	builder.WriteString("当前语音记录配置：\n")
	builder.WriteString(fmt.Sprintf("• 语音模型名称：%s\n", recCfg.ModelName))
	builder.WriteString(fmt.Sprintf("• 语音模型ID：%s\n", recCfg.ModelID))
	builder.WriteString(fmt.Sprintf("• 自定义群ID：%d\n", recCfg.Customgid))
	return builder.String()
}
