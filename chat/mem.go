package chat

import (
	"os"
	"path"
	"strconv"
	"strings"
	"sync"

	binutils "github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/file"
	"github.com/RomiChan/syncx"
)

var mems = syncx.Map[string, *memstorage]{}

type memstorage struct {
	mu   sync.Mutex
	root string
}

func atomicgetmemstorage(service string) (*memstorage, error) {
	d := path.Join(file.BOTPATH, "data", service, "agent")
	err := os.MkdirAll(d, 0755)
	if err != nil {
		return nil, err
	}
	x, _ := mems.LoadOrStore(service, &memstorage{root: d})
	return x, nil
}

func (ms *memstorage) Save(grp int64, text string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	p := path.Join(ms.root, strconv.FormatInt(grp, 10)+".txt")
	f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(text)
	if err != nil {
		return err
	}
	_, err = f.WriteString("\n")
	return err
}

func (ms *memstorage) Load(grp int64) []string {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	p := path.Join(ms.root, strconv.FormatInt(grp, 10)+".txt")
	data, err := os.ReadFile(p)
	if err != nil {
		return nil
	}
	return []string{strings.TrimSpace(binutils.BytesToString(data))}
}
