// Package pool 图片缓存池
//
//nolint:revive
package pool

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wdvxdr1123/ZeroBot/message"

	"github.com/FloatTech/zbputils/ctxext"
)

const cacheurl = "https://gchat.qpic.cn/gchatpic_new//%s/0"

var (
	ErrImgFileOutdated = errors.New("img file outdated")
	ErrNoSuchImg       = errors.New("no such img")
	ErrSendImg         = errors.New("send image error")
	ErrGetMsg          = errors.New("get msg error")
)

var (
	oldimgre = regexp.MustCompile(`^[0-9A-F-]+$`)
)

// Image 图片数据
type Image struct {
	*item
	n, f string
}

// GetImage name
func GetImage(name string) (m *Image, err error) {
	m = new(Image)
	m.n = name
	m.item, err = getItem(name)
	if err == nil && m.u != "" {
		var resp *http.Response
		resp, err = http.Head(m.String())
		if err == nil {
			_ = resp.Body.Close()
		} else {
			goto OUTDATE
		}
		if resp.StatusCode != http.StatusOK {
			goto OUTDATE
		}
		return
	OUTDATE:
		err = ErrImgFileOutdated
		logrus.Debugln("[imgpool] image", name, m, "outdated")
		return
	}
	err = ErrNoSuchImg
	logrus.Debugln("[imgpool] no such image", name)
	return
}

// NewImage context name file
func NewImage(send ctxext.NoCtxSendMsg, get ctxext.NoCtxGetMsg, name, f string) (m *Image, err error) {
	m = new(Image)
	m.n = name
	m.SetFile(f)
	m.item, err = getItem(name)
	if err == nil && m.item.u != "" {
		var resp *http.Response
		resp, err = http.Head(m.String())
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return
			}
		}
		logrus.Debugln("[imgpool] image", name, m, "outdated, updating...")
		get = nil
	}
	err = m.Push(send, get)
	return
}

// String url
func (m *Image) String() string {
	if m.item == nil {
		return m.f
	}
	if oldimgre.MatchString(m.item.u) {
		return fmt.Sprintf(cacheurl, m.item.u)
	}
	nu, err := unpack(m.item.u)
	if err != nil {
		return m.f
	}
	return string(nu)
}

// SetFile f
func (m *Image) SetFile(f string) {
	if strings.HasPrefix(f, "http://") || strings.HasPrefix(f, "https://") || strings.HasPrefix(f, "file:///") {
		m.f = f
	} else {
		m.f = "file:///" + f
	}
}

// Push context
func (m *Image) Push(send ctxext.NoCtxSendMsg, get ctxext.NoCtxGetMsg) (err error) {
	id := send(message.Message{message.Image(m.f)})
	if id == 0 {
		err = ErrSendImg
		return
	}
	if get == nil {
		return
	}
	time.Sleep(2 * time.Second) // LLOneBot get_msg delay
	msg := get(id)
	for _, e := range msg.Elements {
		if e.Type == "image" {
			u := e.Data["url"]
			if ntcachere.MatchString(u) { // is NTQQ
				raw := ""
				raw, err = nturl(u).pack()
				if err != nil {
					logrus.Errorln("[imgpool] pack nturl err:", err)
					err = nil
					return
				}
				m.item, err = newItem(m.n, raw)
				if err != nil {
					logrus.Errorln("[imgpool] get newItem err:", err)
					err = nil
					return
				}
				logrus.Debugln("[imgpool] 缓存:", m.n, "url:", u)
				err = m.item.push("minamoto")
				if err != nil {
					logrus.Errorln("[imgpool] item.push err:", err)
					err = nil
				}
				return
			}
			i := strings.LastIndex(u, "/")
			if i <= 0 {
				break
			}
			u = u[:i]
			i = strings.LastIndex(u, "-")
			if i <= 0 {
				break
			}
			u = u[i:]
			if u == "" {
				break
			}
			m.item, err = newItem(m.n, "0-0"+u)
			if err != nil {
				logrus.Errorln("[imgpool] get newItem err:", err)
				err = nil
				return
			}
			logrus.Debugln("[imgpool] 缓存:", m.n, "url:", "0-0"+u)
			err = m.item.push("minamoto")
			if err != nil {
				logrus.Errorln("[imgpool] item.push err:", err)
				err = nil
			}
			return
		}
	}
	err = ErrGetMsg
	return
}
