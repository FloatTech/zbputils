package pool

import (
	"errors"
	"time"

	"github.com/fumiama/go-registry"
)

type item struct {
	name string
	u    string
}

// newItem 唯一标识文件名 文件链接
func newItem(name, u string) (*item, error) {
	if len(name) > 126 {
		return nil, errors.New("name too long")
	}
	if len(u) > 126 {
		return nil, errors.New("url too long")
	}
	return &item{name: name, u: u}, nil
}

// getItem 唯一标识文件名
func getItem(name string) (*item, error) {
	reg := registry.NewRegReader("reilia.fumiama.top:35354", "fumiama")
	err := reg.ConnectIn(time.Second * 4)
	if err != nil {
		return nil, err
	}
	defer reg.Close()
	u, err := reg.Get(name)
	if err != nil {
		return nil, err
	}
	return &item{name: name, u: u}, nil
}

// push 推送 item
func (t *item) push(key string) (err error) {
	r := registry.NewRegedit("reilia.fumiama.top:35354", "fumiama", key)
	err = r.ConnectIn(time.Second * 4)
	if err != nil {
		return
	}
	defer r.Close()
	err = r.Set(t.name, t.u)
	if err != nil {
		return
	}
	return
}
