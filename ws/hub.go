// Package ws: websocket
package ws

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/happyxhw/pkg/log"
)

// Msg info sent to Client
type Msg struct {
	Action       string `json:"action"`
	ID           string `json:"id"`
	ErrorMessage string `json:"error_message"`
	Data         string `json:"data"`
	Timestamp    int64  `json:"timestamp"`
}

// Hub websocket srv
type Hub struct {
	m       sync.RWMutex
	clients map[int64]map[string]*Client

	upgrader *websocket.Upgrader
}

// NewHub creates a new instance of client hub.
func NewHub() *Hub {
	return &Hub{
		clients: make(map[int64]map[string]*Client),
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// Upgrade conn to websocket
func (s *Hub) Upgrade(w http.ResponseWriter, r *http.Request, h http.Header) (*websocket.Conn, error) {
	return s.upgrader.Upgrade(w, r, h)
}

// Register a client
func (s *Hub) Register(c *Client) {
	s.m.Lock()
	defer s.m.Unlock()
	// close first
	if clients, ok := s.clients[c.userID]; ok {
		if cli, cOk := clients[c.id]; cOk {
			cli.close()
		}
	} else {
		s.clients[c.userID] = make(map[string]*Client)
	}
	s.clients[c.userID][c.id] = c
}

// Remove a client
func (s *Hub) Remove(c *Client) {
	s.m.Lock()
	defer s.m.Unlock()
	delete(s.clients[c.userID], c.id)
	if len(s.clients[c.userID]) == 0 {
		delete(s.clients, c.userID)
	}
}

func (s *Hub) SendToUser(userID int64, msg *Msg) {
	for _, cli := range s.clients[userID] {
		if err := cli.Send(msg); err != nil {
			log.Warn("send msg", zap.Error(err))
		}
	}
}
