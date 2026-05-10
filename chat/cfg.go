package chat

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/fumiama/deepinfra"
	"github.com/fumiama/deepinfra/model"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// AC is the global agent configuration
var AC AgentConfig

var (
	apitypes = map[string]uint8{
		"OpenAI": 0,
		"OLLaMA": 1,
		"GenAI":  2,
	}
	apilist = [3]string{"OpenAI", "OLLaMA", "GenAI"}
)

// ModelType 支持打印 string 并生产 protocal
type ModelType int

func newModelType(typ string) (ModelType, error) {
	t, ok := apitypes[typ]
	if !ok {
		return 0, errors.New("未知类型 " + typ)
	}
	return ModelType(t), nil
}

func (mt ModelType) String() string {
	return apilist[mt]
}

// Protocol creates a protocol instance based on the model type
func (mt ModelType) Protocol(modn string, temp float32, topp float32, maxn uint, reasoning string) (mod model.Protocol, err error) {
	switch AC.Type {
	case 0:
		mod = model.NewOpenAI(
			modn, AC.Separator,
			temp, topp, maxn, reasoning,
		)
	case 1:
		mod = model.NewOLLaMA(
			modn, AC.Separator,
			temp, topp, maxn,
		)
	case 2:
		mod = model.NewGenAI(
			modn,
			temp, topp, maxn,
		)
	default:
		err = errors.New("unsupported model type " + strconv.Itoa(int(AC.Type)))
	}
	return
}

// ModelBool 支持打印成 "是/否"
type ModelBool bool

func (mb ModelBool) String() string {
	if mb {
		return "是"
	}
	return "否"
}

// ModelKey 支持隐藏密钥
type ModelKey string

func (mk ModelKey) String() string {
	if len(mk) == 0 {
		return "未设置"
	}
	if len(mk) <= 4 {
		return "****"
	}
	key := string(mk)
	return key[:2] + strings.Repeat("*", len(key)-4) + key[len(key)-2:]
}

// AgentConfig holds the configuration for the chat agent
type AgentConfig struct {
	ModelName       string
	ImageModelName  string
	AgentModelName  string
	Type            ModelType
	ImageType       ModelType
	AgentType       ModelType
	MaxN            uint
	TopP            float32
	SystemP         string
	AgentChar       string
	AgentSex        string
	API             string
	ImageAPI        string
	AgentAPI        string
	Key             ModelKey
	ImageKey        ModelKey
	AgentKey        ModelKey
	Separator       string
	ReasoningEffort string
	NoSystemP       ModelBool
}

func newconfig() AgentConfig {
	return AgentConfig{
		ModelName: model.ModelDeepDeek,
		SystemP:   SystemPrompt,
		API:       deepinfra.OpenAIDeepInfra,
	}
}

func (c *AgentConfig) String() string {
	topp, maxn := c.MParams()
	sb := strings.Builder{}
	fmt.Fprintf(&sb, "• 模型名：%s\n", c.ModelName)
	fmt.Fprintf(&sb, "• 图像模型名：%s\n", c.ImageModelName)
	fmt.Fprintf(&sb, "• Agent模型名：%s\n", c.AgentModelName)
	fmt.Fprintf(&sb, "• 接口类型：%v\n", c.Type)
	fmt.Fprintf(&sb, "• 图像接口类型：%v\n", c.ImageType)
	fmt.Fprintf(&sb, "• Agent接口类型：%v\n", c.AgentType)
	fmt.Fprintf(&sb, "• 最大长度：%d\n", maxn)
	fmt.Fprintf(&sb, "• TopP：%.1f\n", topp)
	fmt.Fprintf(&sb, "• 系统提示词：%s\n", c.SystemP)
	fmt.Fprintf(&sb, "• Agent性格：%s\n", c.AgentChar)
	fmt.Fprintf(&sb, "• Agent性别：%s\n", c.AgentSex)
	fmt.Fprintf(&sb, "• 接口地址：%s\n", c.API)
	fmt.Fprintf(&sb, "• 图像接口地址：%s\n", c.ImageAPI)
	fmt.Fprintf(&sb, "• Agent接口地址：%s\n", c.AgentAPI)
	fmt.Fprintf(&sb, "• 密钥：%v\n", c.Key)
	fmt.Fprintf(&sb, "• 图像密钥：%v\n", c.ImageKey)
	fmt.Fprintf(&sb, "• Agent密钥：%v\n", c.AgentKey)
	fmt.Fprintf(&sb, "• 分隔符：%s\n", c.Separator)
	fmt.Fprintf(&sb, "• 推理努力度：%s\n", c.ReasoningEffort)
	fmt.Fprintf(&sb, "• 支持系统提示词：%v\n", !c.NoSystemP)
	return sb.String()
}

func (c *AgentConfig) isvalid() bool {
	return c.Key != ""
}

// MParams returns the global model parameters: TopP and MaxN
func (c *AgentConfig) MParams() (topp float32, maxn uint) {
	// 处理TopP参数
	topp = c.TopP
	if topp == 0 {
		topp = 0.9
	}

	// 处理最大长度参数
	maxn = c.MaxN
	if maxn == 0 {
		maxn = 4096
	}

	return topp, maxn
}

// EnsureConfig ensures the configuration is loaded and valid
func EnsureConfig(ctx *zero.Ctx) bool {
	c, ok := ctx.State["manager"].(*ctrl.Control[*zero.Ctx])
	if !ok {
		return false
	}
	cfgp := &AC
	if !cfgp.isvalid() {
		err := c.GetExtra(cfgp)
		if err != nil {
			logrus.Debugln("ERROR: get extra err:", err)
		}
		if !cfgp.isvalid() {
			AC = newconfig()
		}
	}
	if AgentCharConfig.Sex != AC.AgentSex {
		AC.AgentSex = AgentCharConfig.Sex
	}
	if AgentCharConfig.Chars != AC.AgentChar {
		AC.AgentChar = AgentCharConfig.Chars
	}
	return true
}

// NewExtraSetStr creates a handler to set a string-based extra config value
func NewExtraSetStr[T ~string](ptr *T) func(ctx *zero.Ctx) {
	return func(ctx *zero.Ctx) {
		args := strings.TrimSpace(ctx.State["args"].(string))
		c, ok := ctx.State["manager"].(*ctrl.Control[*zero.Ctx])
		if !ok {
			ctx.SendChain(message.Text("ERROR: no such plugin"))
			return
		}
		*ptr = T(args)
		err := c.SetExtra(&AC)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: set extra err: ", err))
			return
		}
		ctx.SendChain(message.Text("成功"))
	}
}

// NewExtraSetBool creates a handler to set a boolean-based extra config value
func NewExtraSetBool[T ~bool](ptr *T) func(ctx *zero.Ctx) {
	return func(ctx *zero.Ctx) {
		args := ctx.State["regex_matched"].([]string)
		isno := args[1] == "不"
		c, ok := ctx.State["manager"].(*ctrl.Control[*zero.Ctx])
		if !ok {
			ctx.SendChain(message.Text("ERROR: no such plugin"))
			return
		}
		*ptr = T(isno)
		err := c.SetExtra(&AC)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: set extra err: ", err))
			return
		}
		ctx.SendChain(message.Text("成功"))
	}
}

// NewExtraSetUint creates a handler to set a uint extra config value
func NewExtraSetUint(ptr *uint) func(ctx *zero.Ctx) {
	return func(ctx *zero.Ctx) {
		args := strings.TrimSpace(ctx.State["args"].(string))
		if args == "" {
			ctx.SendChain(message.Text("ERROR: empty args"))
			return
		}
		c, ok := ctx.State["manager"].(*ctrl.Control[*zero.Ctx])
		if !ok {
			ctx.SendChain(message.Text("ERROR: no such plugin"))
			return
		}
		n, err := strconv.ParseUint(args, 10, 64)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: parse args err: ", err))
			return
		}
		*ptr = uint(n)
		err = c.SetExtra(&AC)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: set extra err: ", err))
			return
		}
		ctx.SendChain(message.Text("成功"))
	}
}

// NewExtraSetFloat32 creates a handler to set a float32 extra config value
func NewExtraSetFloat32(ptr *float32) func(ctx *zero.Ctx) {
	return func(ctx *zero.Ctx) {
		args := strings.TrimSpace(ctx.State["args"].(string))
		if args == "" {
			ctx.SendChain(message.Text("ERROR: empty args"))
			return
		}
		c, ok := ctx.State["manager"].(*ctrl.Control[*zero.Ctx])
		if !ok {
			ctx.SendChain(message.Text("ERROR: no such plugin"))
			return
		}
		n, err := strconv.ParseFloat(args, 32)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: parse args err: ", err))
			return
		}
		*ptr = float32(n)
		err = c.SetExtra(&AC)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: set extra err: ", err))
			return
		}
		ctx.SendChain(message.Text("成功"))
	}
}

// NewExtraSetModelType creates a handler to set a ModelType extra config value
func NewExtraSetModelType(ptr *ModelType) func(ctx *zero.Ctx) {
	return func(ctx *zero.Ctx) {
		args := strings.TrimSpace(ctx.State["args"].(string))
		if args == "" {
			ctx.SendChain(message.Text("ERROR: empty args"))
			return
		}
		c, ok := ctx.State["manager"].(*ctrl.Control[*zero.Ctx])
		if !ok {
			ctx.SendChain(message.Text("ERROR: no such plugin"))
			return
		}
		typ, err := newModelType(args)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err))
			return
		}
		*ptr = typ
		err = c.SetExtra(&AC)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: set extra err: ", err))
			return
		}
		ctx.SendChain(message.Text("成功"))
	}
}
