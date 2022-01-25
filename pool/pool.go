// Package pool url 缓存池
package pool

import (
	"errors"
	"sync"
	"time"

	"github.com/fumiama/go-registry"
)

var (
	reg         = registry.NewRegReader("reilia.fumiama.top:35354", "fumiama")
	wg          sync.WaitGroup
	connmu      sync.Mutex
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
	wg.Add(1)
	defer wg.Done()
	connmu.Lock()
	if !isconnected {
		err := reg.ConnectIn(time.Second * 4)
		if err != nil {
			return nil, err
		}
		isconnected = true
		connmu.Unlock()
		go func() {
			wg.Wait()
			_ = reg.Close()
			connmu.Lock()
			isconnected = false
			connmu.Unlock()
		}()
	} else {
		connmu.Unlock()
	}
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
