package web

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
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
			return dialer.DialContext(ctx, "tcp", "8.8.8.8") // 通过tcp请求nameserver解析域名
		},
	}

	iptables = make(map[string][]string)
	iptmu    sync.RWMutex

	// PixivClient P站特殊客户端
	PixivClient = &http.Client{
		// 解决中国大陆无法访问的问题
		Transport: &http.Transport{
			// 更改 dns
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				host, port, err := net.SplitHostPort(addr)
				if err != nil {
					return nil, err
				}

				iptmu.RLock()
				ips, ok := iptables[host]
				iptmu.RUnlock()
				if !ok {
					ips, err = resolver.LookupHost(ctx, host) // 通过自定义nameserver查询域名
					if err != nil {
						return nil, err
					}
					iptmu.Lock()
					iptables[host] = ips
					iptmu.Unlock()
				}

				for _, ip := range ips {
					// 创建链接
					conn, err := dialer.DialContext(ctx, network, ip+":"+port)
					if err == nil {
						return conn, nil
					}
				}

				return dialer.DialContext(ctx, network, addr)
			},
			// 隐藏 sni 标志
			TLSClientConfig: &tls.Config{
				ServerName:         "-",
				InsecureSkipVerify: true,
			},
			DisableKeepAlives: true,
		},
	}
)
