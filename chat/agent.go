package chat

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"strconv"
	"strings"
	"sync"
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

// agentchar 将 char.yaml 内容嵌入为默认 agent 性格
//
//go:embed char.yaml
var agentchar []byte

var agentcharcfg charcfg

// AgentCharConfig 可随时更改的 agent 性格配置
var AgentCharConfig *goba.Config

var ags = syncx.Map[int64, *goba.Agent]{}

type charcfg struct {
	Sex  string `yaml:"sex"`
	Char string `yaml:"char"`
}

func init() {
	err := yaml.NewDecoder(bytes.NewReader(agentchar)).Decode(&agentcharcfg)
	if err != nil {
		panic(err)
	}
	AgentCharConfig = &goba.Config{
		Nickname: strings.Join(zero.BotConfig.NickName, "/"),
		Sex:      agentcharcfg.Sex,
		Chars:    agentcharcfg.Char,
	}
}

// AgentOf id is self_id
func AgentOf(id int64, service string) *goba.Agent {
	if ag, ok := ags.Load(id); ok {
		return ag
	}
	mem, err := atomicgetmemstorage(service)
	if err != nil {
		panic(err)
	}
	ag := goba.NewAgent(AgentCharConfig, 16, 8, time.Hour*24, "", mem, true, false)
	ags.Store(id, &ag)
	return &ag
}

// ResetAgents reset all agent log
func ResetAgents() {
	ks := make([]int64, 0, 8)
	ags.Range(func(key int64, _ *goba.Agent) bool {
		ks = append(ks, key)
		return true
	})
	for _, k := range ks {
		ags.Delete(k)
	}
}

// ResetAgentCharConfig resert agent char & sex
func ResetAgentCharConfig() {
	AC.AgentChar = agentcharcfg.Char
	AC.AgentSex = agentcharcfg.Sex
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
func CallAgent(ag *goba.Agent, issudo bool, iter int, api deepinfra.API, p model.Protocol, grp int64, role goba.PermRole) []zero.APIRequest {
	reqs, err := ag.GetAction(api, p, grp, role, iter, false)
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
				msgb, err := json.Marshal(req.Params["message"])
				if err != nil {
					logrus.Warnln("[chat] re-marshal msg err:", err)
					continue
				}
				msgs := message.ParseMessageFromArray(gjson.Parse(binary.BytesToString(msgb)))
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
	msgd := ev.NativeMessage
	if len(msgd) > 1024 {
		msg := message.ParseMessage(msgd)
		for _, m := range msg {
			for k, v := range m.Data {
				if len(v) > 512 {
					m.Data[k] = v[:200] + " ... " + v[len(v)-200:]
				}
			}
		}
		raw, err := json.Marshal(&msg)
		if err != nil {
			logrus.Warnln("[chat] togobaev marshal truncated", len(msgd), "bytes message err:", err)
			msgd = []byte(`[]`)
		} else {
			msgd = raw
		}
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
		Message:     msgd,
	}
}

func truncatecopy(params map[string]any) map[string]any {
	cpp := make(map[string]any, len(params))
	for k, v := range params {
		switch val := v.(type) {
		case string:
			if len(val) > 512 {
				cpp[k] = val[:200] + " ... " + val[len(val)-200:]
			} else {
				cpp[k] = val
			}
		case int64, int:
			cpp[k] = val
		case message.Message:
			valCopy := make(message.Message, len(val))
			copy(valCopy, val)
			for i, m := range val {
				valCopy[i] = message.Segment{
					Type: m.Type,
					Data: make(map[string]string, len(m.Data)),
				}
				for k, v := range m.Data {
					if len(v) > 512 {
						valCopy[i].Data[k] = v[:200] + " ... " + v[len(v)-200:]
					} else {
						valCopy[i].Data[k] = v
					}
				}
			}
			cpp[k] = valCopy
		case message.Segment:
			valCopy := message.Segment{
				Type: val.Type,
				Data: make(map[string]string, len(val.Data)),
			}
			for k, v := range val.Data {
				if len(v) > 512 {
					valCopy.Data[k] = v[:200] + " ... " + v[len(v)-200:]
				} else {
					valCopy.Data[k] = v
				}
			}
			cpp[k] = valCopy
		default:
			logrus.Warnln("[chat] agent unknown params typ:", reflect.TypeOf(v))
		}
	}
	return cpp
}

var logevmu sync.Mutex

func logev(ctx *zero.Ctx) {
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

	logevmu.Lock()
	if _, ok := ctx.State[zero.StateKeyPrefixKeep+"_chat_ag_hooked__"]; ok {
		logevmu.Unlock()
		return
	}
	ctx.State[zero.StateKeyPrefixKeep+"_chat_ag_hooked__"] = struct{}{}
	logevmu.Unlock()

	vevent.HookCtxCaller(ctx, vevent.NewAPICallerReturnHook(
		ctx, func(req zero.APIRequest, rsp zero.APIResponse, _ error) {
			gid := ctx.Event.GroupID
			if gid == 0 {
				gid = -ctx.Event.UserID
			}
			logevmu.Lock()
			if _, ok := ctx.State[zero.StateKeyPrefixKeep+"_chat_ag_triggered__"]; ok {
				logevmu.Unlock()
				logrus.Debugln("[chat] agent", gid, "skip agent triggered requ:", &req)
				return
			}
			logevmu.Unlock()
			if req.Action != "send_private_msg" &&
				req.Action != "send_group_msg" &&
				req.Action != "delete_msg" {
				logrus.Debugln("[chat] agent", gid, "skip non-msg other triggered action:", req.Action)
				return
			}
			req.Params = truncatecopy(req.Params)
			ag := AgentOf(ctx.Event.SelfID, "aichat")
			logrus.Debugln("[chat] agent others", gid, "add requ:", &req)
			ag.AddRequest(gid, &req)
			logrus.Debugln("[chat] agent others", gid, "add resp:", &rsp)
			ag.AddResponse(gid, &goba.APIResponse{
				Status:  rsp.Status,
				Data:    json.RawMessage(rsp.Data.Raw),
				Message: rsp.Message,
				Wording: rsp.Wording,
				RetCode: rsp.RetCode,
			})
		}),
	)
}

func init() {
	zero.OnNotice().FirstPriority().SetBlock(false).Handle(logev)
	zero.OnMessage().FirstPriority().SetBlock(false).Handle(logev)
	zero.OnRequest().FirstPriority().SetBlock(false).Handle(logev)
}
