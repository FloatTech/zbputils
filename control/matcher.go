package control

import zero "github.com/wdvxdr1123/ZeroBot"

type Matcher interface {
	// SetBlock 设置是否阻断后面的 Matcher 触发
	SetBlock(block bool) Matcher
	// Handle 直接处理事件
	Handle(handler zero.Handler)
}

type matcherinstance struct {
	m *zero.Matcher
}

func (m matcherinstance) SetBlock(block bool) Matcher {
	_ = m.m.SetBlock(block)
	return m
}

func (m matcherinstance) Handle(handler zero.Handler) {
	_ = m.m.Handle(handler)
}
