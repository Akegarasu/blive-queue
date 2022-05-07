package main

import (
	"github.com/Akegarasu/blive-queue/eio"
	bliveClient "github.com/Akegarasu/blivedm-go/client"
	"github.com/Akegarasu/blivedm-go/message"
	"github.com/tidwall/gjson"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	Eio           *eio.Server
	DanmakuClient *bliveClient.Client
	Queue         *Queue
	Rule          *Rule
	RoomID        string
}

func NewServer() *Server {
	return &Server{
		Eio:           eio.NewServer(),
		DanmakuClient: nil,
		Queue:         NewQueue(),
		Rule:          DefaultRule(),
	}
}

func (s *Server) Init() {
	s.Eio.RegisterEventHandler("HEARTBEAT", func(event *eio.Event) {
		log.Debug("heartbeat")
	})

	s.Eio.RegisterEventHandler("CONNECT_DANMAKU", func(event *eio.Event) {
		s.ConnectDanmakuServer(event.Data)
	})

	s.Eio.RegisterEventHandler("APPLY_RULE", func(event *eio.Event) {
		s.Rule = NewRule(event.Data)
		log.Infof("设置了新的过滤规则 关键词: %s, 最大人数: %d, 仅舰长: %v, 最低牌子等级: %d", s.Rule.keyword, s.Rule.maxQueueLength, s.Rule.guardOnly, s.Rule.minMedalLevel)
	})

	s.Eio.RegisterEventHandler("ADD_USER", func(event *eio.Event) {
		_ = s.Eio.BoardCastEvent(*event)
		log.Info("测试用户: ", event.Data)
	})

	s.Eio.RegisterEventHandler("REMOVE_USER", func(event *eio.Event) {
		uid, err := strconv.Atoi(event.Data)
		if err != nil {
			log.Error("删除失败: uid转换出错")
		}
		ok := s.Queue.Remove(uid)
		if !ok {
			log.Error("删除失败啦")
			return
		}
		_ = s.Eio.BoardCastEvent(*event)
		log.Info("删除了uid: ", event.Data)
	})

	s.Eio.RegisterEventHandler("REMOVE_ALL", func(event *eio.Event) {
		s.Queue.Clear()
		_ = s.Eio.BoardCastEvent(*event)
		log.Info("清空了排队")
	})

	s.Eio.RegisterEventHandler("RESORT", func(event *eio.Event) {
		j := gjson.Parse(event.Data)
		oldIndex := int(j.Get("oldIndex").Int())
		newIndex := int(j.Get("newIndex").Int())
		s.Queue.Resort(oldIndex, newIndex)
		_ = s.Eio.BoardCastEventExceptSelf(*event)
		log.Infof("排序: %d -> %d", oldIndex, newIndex)
	})
}

func (s *Server) ConnectDanmakuServer(roomID string) {
	if s.DanmakuClient != nil {
		s.DanmakuClient.Stop()
	}
	c := bliveClient.NewClient(roomID)
	c.OnDanmaku(s.HandleDanmaku)
	c.UseDefaultHost()
	err := c.ConnectAndStart()
	if err != nil {
		log.Warn("连接弹幕服务器出错")
	}
	s.DanmakuClient = c
	s.RoomID = roomID
	log.Info("连接到房间: ", roomID)
}

func (s *Server) HandleDanmaku(d *message.Danmaku) {
	if d.Content == s.Rule.keyword {
		s.HandleJoinQueue(d)
		return
	}
	if d.Content == s.Rule.cancelKeyword {
		s.HandleLeaveQueue(d)
		return
	}
}

func (s *Server) HandleJoinQueue(d *message.Danmaku) {
	if !s.Rule.Filter(d, s.RoomID) {
		return
	}
	if s.Queue.Que.Len() >= s.Rule.maxQueueLength {
		log.Error("排队失败: 队列满了")
		return
	}
	ok := s.Queue.Add(d.Sender)
	if ok {
		log.Infof("添加排队成功: %s (uid: %d)", d.Sender.Uname, d.Sender.Uid)
		err := s.Eio.BoardCastEvent(eio.Event{
			EventName: "ADD_USER",
			Data:      d.Raw,
		})
		if err != nil {
			log.Error("排队失败")
		}
	}
}

func (s *Server) HandleLeaveQueue(d *message.Danmaku) {
	if s.Queue.In(d.Sender) {
		s.Queue.Remove(d.Sender.Uid)
		_ = s.Eio.BoardCastEvent(eio.Event{
			EventName: "REMOVE_USER",
			Data:      strconv.Itoa(d.Sender.Uid),
		})
		log.Infof("取消排队成功: %s (uid: %d)", d.Sender.Uname, d.Sender.Uid)
	} else {
		log.Errorf("取消排队失败: %s (uid: %d) 根本没有排队哦", d.Sender.Uname, d.Sender.Uid)
	}
}
