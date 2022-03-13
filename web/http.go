// Package web 网络处理相关
package web

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/sirupsen/logrus"
)

// NewDefaultClient ...
func NewDefaultClient() *http.Client {
	return &http.Client{}
}

// NewTLS12Client ...
func NewTLS12Client() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MaxVersion: tls.VersionTLS12,
			},
		},
	}
}

// NewPixivClient P站特殊客户端
func NewPixivClient() *http.Client {
	return &http.Client{
		// 解决中国大陆无法访问的问题
		Transport: &http.Transport{
			// 更改 dns
			Dial: func(network, addr string) (net.Conn, error) {
				host, port, err := net.SplitHostPort(addr)
				if err != nil {
					return nil, err
				}

				iptmu.RLock()
				ips, ok := iptables[host]
				iptmu.RUnlock()
				if !ok {
					ips, err = resolver.LookupHost(context.TODO(), host) // 通过自定义nameserver查询域名
					if err != nil {
						return nil, err
					}
					iptmu.Lock()
					iptables[host] = ips
					iptmu.Unlock()
				}

				for _, ip := range ips {
					// 创建链接
					conn, err := net.Dial(network, ip+":"+port)
					if err == nil {
						logrus.Debugln("[web]google DNS get host", host, ip+":"+port)
						return conn, nil
					}
				}

				return net.Dial(network, addr)
			},
			// 隐藏 sni 标志
			TLSClientConfig: &tls.Config{
				ServerName:         "-",
				InsecureSkipVerify: true,
			},
			DisableKeepAlives: true,
		},
	}
}

// GetDataWith 使用自定义请求头获取数据
func GetDataWith(client *http.Client, url string, method string, referer string, ua string) (data []byte, err error) {
	// 提交请求
	var request *http.Request
	request, err = http.NewRequest(method, url, nil)
	if err == nil {
		// 增加header选项
		request.Header.Add("Referer", referer)
		request.Header.Add("User-Agent", ua)
		var response *http.Response
		response, err = client.Do(request)
		if err == nil {
			if response.StatusCode != http.StatusOK {
				s := fmt.Sprintf("status code: %d", response.StatusCode)
				err = errors.New(s)
				return
			}
			data, err = io.ReadAll(response.Body)
			response.Body.Close()
		}
	}
	return
}

// GetData 获取数据
func GetData(url string) (data []byte, err error) {
	var response *http.Response
	response, err = http.Get(url)
	if err == nil {
		if response.StatusCode != http.StatusOK {
			s := fmt.Sprintf("status code: %d", response.StatusCode)
			err = errors.New(s)
			return
		}
		data, err = io.ReadAll(response.Body)
		response.Body.Close()
	}
	return
}

// PostData 获取数据
func PostData(url string, contentType string, body io.Reader) (data []byte, err error) {
	var response *http.Response
	response, err = http.Post(url, contentType, body)
	if err == nil {
		if response.StatusCode != http.StatusOK {
			s := fmt.Sprintf("status code: %d", response.StatusCode)
			err = errors.New(s)
			return
		}
		data, err = io.ReadAll(response.Body)
		response.Body.Close()
	}
	return
}
