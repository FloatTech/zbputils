package web

import (
	"context"
	"crypto/tls"
	"net"
	"sync"
)

var (
	// 定义resolver
	resolver = &net.Resolver{
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			if IsSupportIPv6 {
				return tls.Dial("tcp", "[2001:4860:4860::8888]:853", nil) // 通过tls请求nameserver解析域名
			}
			return tls.Dial("tcp", "8.8.8.8:853", nil) // 通过tls请求nameserver解析域名
		},
	}

	iptables = make(map[string][]string)
	iptmu    sync.RWMutex
)
