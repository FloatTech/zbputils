package web

import (
	"context"
	"net"
	"sync"
	"time"
)

var (
	// 默认dialer
	dialer = &net.Dialer{
		Timeout: 5 * time.Second,
	}

	// 定义resolver
	resolver = &net.Resolver{
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			return dialer.DialContext(ctx, "tcp", "8.8.8.8:53") // 通过tcp请求nameserver解析域名
		},
	}

	iptables = make(map[string][]string)
	iptmu    sync.RWMutex
)
