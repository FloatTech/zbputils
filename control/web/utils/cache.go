// Package utils 工具类
package utils

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	// LoginCache 登录缓存
	LoginCache = cache.New(24*time.Hour, 12*time.Hour)
)
