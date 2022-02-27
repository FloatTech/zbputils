// Package web 网络处理相关
package web

import (
	"crypto/tls"
	"io"
	"net/http"
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
		data, err = io.ReadAll(response.Body)
		response.Body.Close()
	}
	return
}
