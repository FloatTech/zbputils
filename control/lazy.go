package control

import (
	"errors"
	"strings"
	"unicode"

	"github.com/FloatTech/zbputils/file"
)

// 下载并获取本 engine 文件夹下的懒加载数据
func (e *engineinstance) GetLazyData(filename string, isDataMustEqual bool) ([]byte, error) {
	if e.datafolder == "" {
		return nil, errors.New("datafolder is empty")
	}
	if !strings.HasSuffix(e.datafolder, "/") || !strings.HasPrefix(e.datafolder, "data/") || !unicode.IsUpper(rune(e.datafolder[5])) {
		return nil, errors.New("invalid datafolder")
	}
	return file.GetLazyData(e.datafolder+filename, isDataMustEqual)
}
