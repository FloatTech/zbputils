package chat

import (
	"bytes"
	_ "embed"

	"github.com/RomiChan/syncx"
	goba "github.com/fumiama/go-onebot-agent"
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
		cfg.Sex, cfg.Char, cfg.Default,
	)
	ags.Store(id, &ag)
	return &ag
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
