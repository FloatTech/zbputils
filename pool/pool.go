// Package pool url 缓存池
package pool

import (
	"errors"
	"sync"
	"time"

	"github.com/fumiama/go-registry"

	"github.com/FloatTech/zbputils/process"
)

var (
	reg         = registry.NewRegReader("reilia.fumiama.top:35354", "fumiama")
	wg          sync.WaitGroup
	isconnected bool
)

type Item struct {
	name string
	u    string
}

// NewItem 唯一标识文件名 文件链接
func NewItem(name, u string) (*Item, error) {
	if len(name) > 126 {
		return nil, errors.New("name too long")
	}
	if len(u) > 126 {
		return nil, errors.New("url too long")
	}
	return &Item{name: name, u: u}, nil
}

// GetItem 唯一标识文件名
func GetItem(name string) (*Item, error) {
	if !isconnected {
		err := reg.ConnectIn(time.Second * 4)
		if err != nil {
			return nil, err
		}
		isconnected = true
		go func() {
			for c := 0; c < 10; c++ {
				process.SleepAbout1sTo2s()
				wg.Wait()
			}
			_ = reg.Close()
			isconnected = false
		}()
	}
	wg.Add(1)
	defer wg.Done()
	u, err := reg.Get(name)
	if err != nil {
		return nil, err
	}
	return &Item{name: name, u: u}, nil
}

// Push 推送 item
func (t *Item) Push(key string) (err error) {
	r := registry.NewRegedit("reilia.fumiama.top:35354", "fumiama", key)
	err = r.ConnectIn(time.Second * 4)
	if err != nil {
		return
	}
	err = r.Set(t.name, t.u)
	if err != nil {
		return
	}
	err = r.Close()
	return
}

// String item 的 url
func (t *Item) String() string {
	return t.u
}

// Name item 的 name
func (t *Item) Name() string {
	return t.name
}
