package control

import (
	"errors"
	"strings"
	"unicode"

	"github.com/FloatTech/zbputils/file"
	"github.com/sirupsen/logrus"
)

// GetLazyData 下载并获取本 engine 文件夹下的懒加载数据
func (e *engineinstance) GetLazyData(filename string, isDataMustEqual bool) ([]byte, error) {
	if e.datafolder == "" {
		return nil, errors.New("datafolder is empty")
	}
	if !strings.HasSuffix(e.datafolder, "/") || !strings.HasPrefix(e.datafolder, "data/") || !unicode.IsUpper(rune(e.datafolder[5])) {
		return nil, errors.New("invalid datafolder")
	}
	return file.GetLazyData(e.datafolder+filename, isDataMustEqual)
}

// InitWhenNoError 在 errfun 无误时执行 do
func (e *engineinstance) InitWhenNoError(errfun func() error, do func()) {
	err := errfun()
	if err != nil {
		logrus.Warn("[lazy] stop init plugin", e.service, "for error:", err)
		return
	}
	do()
}
