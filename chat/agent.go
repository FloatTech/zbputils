package chat

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"strconv"
	"time"

	"github.com/RomiChan/syncx"
	"github.com/fumiama/deepinfra"
	"github.com/fumiama/deepinfra/model"
	goba "github.com/fumiama/go-onebot-agent"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v3"

	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/zbputils/vevent"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
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
func AgentOf(id int64, service string) *goba.Agent {
	if ag, ok := ags.Load(id); ok {
		return ag
	}
	var cfg charcfg
	err := yaml.NewDecoder(bytes.NewReader(AgentChar)).Decode(&cfg)
	if err != nil {
		panic(err)
	}
	mem, err := atomicgetmemstorage(service)
	if err != nil {
		panic(err)
	}
	ag := goba.NewAgent(
		id, 16, 8, time.Hour*24,
		zero.BotConfig.NickName[0],
		cfg.Sex, cfg.Char, cfg.Default, mem, true, false,
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
		if !errors.Is(err, io.EOF) {
			logrus.Warnln("[chat] agent err:", err, reqs)
		} else {
			logrus.Infoln("[chat] agent end action")
		}
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
			if role == goba.PermRoleUser { // check @all
				msg, ok := req.Params["message"].(string)
				if !ok {
					logrus.Warnln("[chat] invalid message type", reflect.TypeOf(req.Params["message"]))
					continue
				}
				msgs := message.ParseMessageFromArray(gjson.Parse(msg))
				for _, m := range msgs {
					if m.Type == "at" {
						qqs, ok := m.Data["qq"]
						if !ok {
							logrus.Warnln("[chat] invalid at message without qq")
							continue
						}
						qq, err := strconv.ParseInt(qqs, 10, 64)
						if err != nil {
							logrus.Warnln("[chat] invalid at qq", qqs)
							continue
						}
						if qq <= 0 {
							logrus.Warnln("[chat] invalid at qq num", qq)
							continue
						}
					}
				}
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
	vevent.HookCtxCaller(ctx, vevent.NewAPICallerReturnHook(
		ctx, func(req zero.APIRequest, rsp zero.APIResponse, err error) {
			gid := ctx.Event.GroupID
			if gid == 0 {
				gid = -ctx.Event.UserID
			}
			ag := AgentOf(ctx.Event.SelfID, "aichat")
			logrus.Infoln("[chat] agent", gid, "add requ:", &req)
			ag.AddRequest(gid, &req)
			logrus.Infoln("[chat] agent", gid, "get resp:", &rsp)
			ag.AddResponse(gid, &goba.APIResponse{
				Status:  rsp.Status,
				Data:    json.RawMessage(rsp.Data.Raw),
				Message: rsp.Message,
				Wording: rsp.Wording,
				RetCode: rsp.RetCode,
			})
		}),
	)
	gid := ctx.Event.GroupID
	if gid == 0 {
		gid = -ctx.Event.UserID
	}
	ev := togobaev(ctx.Event)
	if ev == nil {
		return
	}
	data, _ := json.Marshal(ev)
	logrus.Debugln("[chat] agent", gid, "add ev:", binary.BytesToString(data))
	AgentOf(ctx.Event.SelfID, "aichat").AddEvent(gid, ev)
}

func init() {
	zero.OnNotice().FirstPriority().SetBlock(false).Handle(logev)
	zero.OnMessage().FirstPriority().SetBlock(false).Handle(logev)
	zero.OnRequest().FirstPriority().SetBlock(false).Handle(logev)
}
