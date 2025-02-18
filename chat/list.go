package chat

import (
	"sync"

	"github.com/fumiama/deepinfra"
	"github.com/fumiama/deepinfra/model"
)

// listcap cannot < 2
const listcap = 8

type list struct {
	mu sync.RWMutex
	m  map[int64][]string
}

func newlist() list {
	return list{
		m: make(map[int64][]string, 64),
	}
}

func (l *list) add(grp int64, usr, txt string, isme bool) {
	if !isme {
		txt = "【" + usr + "】" + txt
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	msgs, ok := l.m[grp]
	if !ok {
		msgs = make([]string, 1, listcap)
		msgs[0] = txt
		l.m[grp] = msgs
		return
	}
	isprevusr := len(msgs)%2 != 0
	if (isprevusr && !isme) || (!isprevusr && isme) { // is same
		msgs[len(msgs)-1] += "\n\n" + txt
		return
	}
	if len(msgs) < cap(msgs) {
		msgs = append(msgs, txt)
		l.m[grp] = msgs
		return
	}
	copy(msgs, msgs[2:])
	msgs[len(msgs)-2] = txt
	l.m[grp] = msgs[:len(msgs)-1]
}

func (l *list) modelize(temp float32, grp int64, mn, sysp, sepstr string) deepinfra.Model {
	m := model.NewCustom(mn, sepstr, temp, 0.9, 1024).System(sysp)
	l.mu.RLock()
	defer l.mu.RUnlock()
	sz := len(l.m[grp])
	if sz == 0 {
		return m.User("自己随机开启新话题")
	}
	for i, msg := range l.m[grp] {
		if i%2 == 0 { // is user
			_ = m.User(msg)
		} else {
			_ = m.Assistant(msg)
		}
	}
	return m
}
