package eio

import "github.com/tidwall/gjson"

type Event struct {
	EventName string `json:"cmd"`
	Data      string `json:"data"`
	Uuid      string `json:"uuid"`
}

//var (
//	HeartBeatEvent = Event{EventName: "HEARTBEAT"}
//	SyncEvent      = Event{EventName: "SYNC"}
//)

func ParseEvent(b []byte, uuid string) *Event {
	g := gjson.ParseBytes(b)
	return &Event{
		EventName: g.Get("cmd").String(),
		Data:      g.Get("data").String(),
		Uuid:      uuid,
	}
}

func (s *Server) RegisterEventHandler(name string, h func(*Event)) {
	s.eventHandlers[name] = h
}

func (s *Server) HandleEvent(e *Event) {
	if h, ok := s.eventHandlers[e.EventName]; ok {
		go h(e)
	}
}
