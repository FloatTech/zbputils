package chat

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"reflect"

	"github.com/RomiChan/syncx"
	"github.com/fumiama/deepinfra"
	"github.com/fumiama/deepinfra/model"
	goba "github.com/fumiama/go-onebot-agent"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	zero "github.com/wdvxdr1123/ZeroBot"
)

// AgentChar 将 char.yaml 内容嵌入为默认 agent 性格
//
//go:embed char.yaml
var AgentChar []byte

var ags = syncx.Map[int64, *goba.Agent]{}

type charcfg struct {
	Sex     string `yaml:"sex"`
	Char    string `yaml:"char"`
	Default string `yaml:"default"`
}

// AgentOf id is self_id
func AgentOf(id int64) *goba.Agent {
	if ag, ok := ags.Load(id); ok {
		return ag
	}
	var cfg charcfg
	err := yaml.NewDecoder(bytes.NewReader(AgentChar)).Decode(&cfg)
	if err != nil {
		panic(err)
	}
	ag := goba.NewAgent(
		id, 16, 8, zero.BotConfig.NickName[0],
		cfg.Sex, cfg.Char, cfg.Default, false,
	)
	ags.Store(id, &ag)
	return &ag
}

var checkgids = map[string]struct{}{
	"send_group_msg":          {},
	"set_group_kick":          {},
	"set_group_ban":           {},
	"set_group_whole_ban":     {},
	"set_group_card":          {},
	"set_group_name":          {},
	"set_group_special_title": {},
	"get_group_member_info":   {},
	"get_group_member_list":   {},
}

// CallAgent and check group API permission
func CallAgent(ag *goba.Agent, issudo bool, api deepinfra.API, p model.Protocol, grp int64, role goba.PermRole) []zero.APIRequest {
	reqs, err := ag.GetAction(api, p, grp, role, false)
	if err != nil {
		logrus.Warnln("[chat] agent err:", err, reqs)
		return nil
	}
	logrus.Infoln("[chat] agent do:", reqs)
	checkedreqs := make([]zero.APIRequest, 0, len(reqs))
	for _, req := range reqs {
		if _, ok := checkgids[req.Action]; ok {
			v, ok := req.Params["group_id"].(json.Number)
			if !ok {
				logrus.Warnln("[chat] invalid group_id type", reflect.TypeOf(req.Params["group_id"]))
				continue
			}
			gid, err := v.Int64()
			if !ok {
				logrus.Warnln("[chat] agent conv req gid err:", err)
				continue
			}
			if grp != gid && !issudo {
				logrus.Warnln("[chat] refuse to send out of grp from", grp, "to", gid)
				continue
			}
		}
		checkedreqs = append(checkedreqs, req)
	}
	return checkedreqs
}

func togobaev(ev *zero.Event) *goba.Event {
	msgid := int64(0)
	if id, ok := ev.MessageID.(int64); ok {
		msgid = id
	} else {
		return nil
	}
	return &goba.Event{
		Time:        ev.Time,
		PostType:    ev.PostType,
		MessageType: ev.MessageType,
		SubType:     ev.SubType,
		MessageID:   msgid,
		GroupID:     ev.GroupID,
		UserID:      ev.UserID,
		TargetID:    ev.TargetID,
		SelfID:      ev.SelfID,
		NoticeType:  ev.NoticeType,
		OperatorID:  ev.OperatorID,
		File:        ev.File,
		RequestType: ev.RequestType,
		Flag:        ev.Flag,
		Comment:     ev.Comment,
		Sender:      ev.Sender,
		Message:     ev.NativeMessage,
	}
}

func logev(ctx *zero.Ctx) {
	// 计算群组 ID（私聊时使用负的 UserID）
	gid := ctx.Event.GroupID
	if gid == 0 {
		gid = -ctx.Event.UserID
	}
	ev := togobaev(ctx.Event)
	if ev == nil {
		return
	}
	AgentOf(ctx.Event.SelfID).AddEvent(gid, ev)
}

func init() {
	zero.OnNotice().FirstPriority().SetBlock(false).Handle(logev)
	zero.OnMessage().FirstPriority().SetBlock(false).Handle(logev)
	zero.OnRequest().FirstPriority().SetBlock(false).Handle(logev)
}
