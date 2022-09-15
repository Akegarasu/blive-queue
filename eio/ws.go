package eio

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
	"runtime/debug"
	"time"
)

type Server struct {
	// TODO: 要加锁
	conn          []*Connection
	eventHandlers map[string]func(*Event)
}

type Connection struct {
	Ws   *websocket.Conn
	Uuid string
}

func NewConnection(ws *websocket.Conn) *Connection {
	return &Connection{
		Ws:   ws,
		Uuid: uuid.NewString(),
	}
}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewServer() *Server {
	s := &Server{
		conn:          make([]*Connection, 0),
		eventHandlers: make(map[string]func(*Event)),
	}
	go s.heartbeat()
	return s
}

func (s *Server) Warp(c *gin.Context) {
	s.Create(c.Writer, c.Request)
}

func (s *Server) Create(w http.ResponseWriter, r *http.Request) {
	c, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}
	conn := NewConnection(c)
	s.conn = append(s.conn, conn)
	log.Info("客户端连接成功~")
	s.listen(conn)
}

func (s *Server) listen(c *Connection) {
	defer func() {
		_ = c.Ws.Close()
		s.removeConn(c)
	}()
	for {
		t, msg, err := c.Ws.ReadMessage()
		if err != nil {
			log.Error(err)
			break
		}
		if t == websocket.TextMessage {
			s.HandleEvent(ParseEvent(msg, c.Uuid))
		}
	}
}

func (s *Server) removeConn(c *Connection) {
	for pos, i := range s.conn {
		if i == c {
			s.conn = append(s.conn[:pos], s.conn[pos+1:]...)
			log.Info("与一个客户端的连接断开了")
			return
		}
	}
	log.Error("断开连接失败")
}

func (s *Server) heartbeat() {
	for {
		for _, c := range s.conn {
			err := c.Ws.WriteMessage(websocket.TextMessage, []byte("{\"cmd\":\"HEARTBEAT\", \"data\":\"\"}"))
			if err != nil {
				log.Error(err)
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func (s *Server) BoardCastEvent(event Event) error {
	defer func() {
		if pan := recover(); pan != nil {
			fmt.Printf("sending error: %v\n%s", pan, debug.Stack())
		}
	}()
	for _, c := range s.conn {
		err := c.Ws.WriteJSON(event)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) BoardCastEventExceptSelf(event Event) error {
	defer func() {
		if pan := recover(); pan != nil {
			fmt.Printf("sending error: %v\n%s", pan, debug.Stack())
		}
	}()
	for _, c := range s.conn {
		if c.Uuid != event.Uuid {
			err := c.Ws.WriteJSON(event)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
