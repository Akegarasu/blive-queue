package main

import (
	"github.com/Akegarasu/blivedm-go/message"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
)

type Rule struct {
	keyword        string
	cancelKeyword  string
	guardOnly      bool
	fuzzyMatch     bool
	minMedalLevel  int
	maxQueueLength int
	admins         []string
	blockUsers     []string
}

func NewRule(s string) *Rule {
	g := gjson.Parse(s)
	admins := strings.Split(g.Get("admins").String(), "\n")
	blockUsers := strings.Split(g.Get("blockUsers").String(), "\n")
	return &Rule{
		keyword:        "排队",
		cancelKeyword:  "取消排队",
		guardOnly:      g.Get("guardOnly").Bool(),
		fuzzyMatch:     g.Get("fuzzyMatch").Bool(),
		maxQueueLength: int(g.Get("maxQueueLength").Int()),
		minMedalLevel:  int(g.Get("minMedalLevel").Int()),
		admins:         admins,
		blockUsers:     blockUsers,
	}
}

func DefaultRule() *Rule {
	return &Rule{
		keyword:        "排队",
		cancelKeyword:  "取消排队",
		guardOnly:      false,
		fuzzyMatch:     false,
		maxQueueLength: 20,
		minMedalLevel:  0,
		admins:         nil,
		blockUsers:     nil,
	}
}

func (r Rule) Filter(danmaku *message.Danmaku, roomID string) bool {
	iu := strconv.Itoa(danmaku.Sender.Uid)
	for _, i := range r.blockUsers {
		if iu == i {
			return false
		}
	}
	// 无 0 总督 1 提督 2 舰长 3
	if r.guardOnly && danmaku.Sender.GuardLevel == 0 {
		return false
	}
	if r.minMedalLevel != 0 {
		rid, err := strconv.Atoi(roomID)
		if err != nil {
			log.Error("房间号转换失败")
		}
		if danmaku.Sender.Medal.UpRoomId != rid {
			return false
		}
		if danmaku.Sender.Medal.Level < r.minMedalLevel {
			return false
		}
	}
	return true
}

// CheckIsAdmin 检查是否为管理员
func (r Rule) CheckIsAdmin(u *message.User) bool {
	uid := strconv.Itoa(u.Uid)
	if InSlice(r.admins, uid) {
		return true
	}
	return false
}
