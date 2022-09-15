package main

import (
	"container/list"
	"encoding/json"
	"sync"

	"github.com/Akegarasu/blivedm-go/message"
)

// Queue 基于 container/list 的封装
type Queue struct {
	Que *list.List
	// Map仅作用来判断的根据，如果需要这里可以换成 struct{} 减少空间占用
	Map map[int]int
	mu  sync.RWMutex
}

type SyncMessage struct {
	Cmd  string     `json:"cmd"`
	Data []LiveUser `json:"data"`
}

type LiveUser struct {
	Uid        string `json:"uid"`
	Nickname   string `json:"nickname"`
	GuardLevel string `json:"level"`
}

func NewLiveUser(u *message.User) LiveUser {
	return LiveUser{
		Uid:        i2s(u.Uid),
		Nickname:   u.Uname,
		GuardLevel: i2s(u.GuardLevel),
	}
}

func (b LiveUser) Json() string {
	marshal, err := json.Marshal(b)
	if err != nil {
		return ""
	}
	return string(marshal)
}

func NewQueue() *Queue {
	return &Queue{
		Que: list.New(),
		Map: make(map[int]int),
		mu:  sync.RWMutex{},
	}
}

func (q *Queue) Add(u *message.User) bool {
	if q.In(u) {
		return false
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	q.Map[u.Uid] = 1
	q.Que.PushBack(u)
	return true
}

func (q *Queue) Remove(uid int) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if _, ok := q.Map[uid]; !ok {
		return false
	}
	delete(q.Map, uid)
	for p := q.Que.Front(); p != nil; p = p.Next() {
		if uid == p.Value.(*message.User).Uid {
			q.Que.Remove(p)
			return true
		}
	}
	return false
}

func (q *Queue) Clear() {
	q.Map = make(map[int]int)
	q.Que.Init()
}

func (q *Queue) Resort(oldIndex int, newIndex int) {
	if oldIndex == newIndex {
		return
	}
	oldItem := q.Get(oldIndex)
	newItem := q.Get(newIndex)
	if oldItem == nil || newItem == nil {
		return
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	if oldIndex < newIndex {
		q.Que.MoveAfter(oldItem, newItem)
	} else {
		q.Que.MoveBefore(oldItem, newItem)
	}
}

func (q *Queue) Get(pos int) *list.Element {
	q.mu.RLock()
	defer q.mu.RUnlock()
	for i, p := 0, q.Que.Front(); p != nil; p = p.Next() {
		if i == pos {
			return p
		}
		i++
	}
	return nil
}

func (q *Queue) In(u *message.User) bool {
	q.mu.RLock()
	defer q.mu.RUnlock()
	if _, ok := q.Map[u.Uid]; ok {
		return true
	}
	return false
}

func (q *Queue) Encode() *SyncMessage {
	d := make([]LiveUser, 0)
	for p := q.Que.Front(); p != nil; p = p.Next() {
		user := p.Value.(*message.User)
		d = append(d, NewLiveUser(user))
	}
	s := &SyncMessage{
		Cmd:  "SYNC",
		Data: d,
	}
	return s
}
