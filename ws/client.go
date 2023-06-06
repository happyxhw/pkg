package ws

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/happyxhw/pkg/log"
)

var (
	ErrClientClosed = errors.New("client closed")
)

const (
	defaultWriteWait = 60 * time.Second

	defaultPongTimeout = 60 * time.Second
	defaultPingPeriod  = 30 * time.Second

	defaultReadLimit = 1024
)

// Client websocket client, held by hub
type Client struct {
	conn    *websocket.Conn // websocket conn
	stopCh  chan struct{}   // signal to stop client
	onClose func(*Client)   // do func when conn close
	stopped bool
	mu      *sync.Mutex

	hub   *Hub
	once  sync.Once
	ch    chan *Msg
	errCh chan error

	id     string // client id
	userID int64  // client user id
}

// NewClient return client instance
func NewClient(conn *websocket.Conn, hub *Hub, id string, userID int64, readFunc func([]byte)) *Client {
	cli := Client{
		conn: conn,
		hub:  hub,

		stopCh: make(chan struct{}, 1),
		ch:     make(chan *Msg),
		errCh:  make(chan error),
		mu:     &sync.Mutex{},

		id:     id,
		userID: userID,
	}
	go cli.startReading(readFunc)
	go cli.startPingHandler()
	return &cli
}

func (c *Client) SetCloseFn(onClose func(*Client)) {
	c.onClose = onClose
}

// Close client
func (c *Client) Close() {
	c.close()
}

// close conn
func (c *Client) close() {
	c.once.Do(func() {
		c.mu.Lock()
		c.stopped = true
		defer c.mu.Unlock()
		c.stopCh <- struct{}{}
		_ = c.conn.Close()
		if c.onClose != nil {
			c.onClose(c)
		}
		c.hub.Remove(c)
		log.Info("client closed", zap.String("id", c.id), zap.Int64("user_id", c.userID))
	})
}

// StartReading starts listening on the Client connection.
// As we do not need anything from the Client,
// we ignore incoming messages. Leaves the loop on errors.
func (c *Client) startReading(readFunc func([]byte)) {
	defer c.close()
	c.conn.SetReadLimit(defaultReadLimit)
	_ = c.conn.SetReadDeadline(time.Now().Add(defaultPongTimeout))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(defaultPongTimeout))
		return nil
	})
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if strings.Contains(err.Error(), "close") {
				log.Info("ws close", zap.String("id", c.id), zap.Int64("user_id", c.userID), zap.Error(err))
			} else {
				log.Error("ws read", zap.String("id", c.id), zap.Int64("user_id", c.userID), zap.Error(err))
			}
			break
		}
		readFunc(msg)
	}
}

// ping loop, quit on error or stop signal
func (c *Client) startPingHandler() {
	pingTicker := time.NewTicker(defaultPingPeriod)
	defer func() {
		c.close()
		pingTicker.Stop()
	}()
	for {
		select {
		case <-pingTicker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(defaultWriteWait))
			err := c.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				if strings.Contains(err.Error(), "close") {
					log.Info("ping", zap.String("id", c.id), zap.Int64("user_id", c.userID), zap.Error(err))
				} else {
					log.Error("ping", zap.String("id", c.id), zap.Int64("user_id", c.userID), zap.Error(err))
				}
				return
			}
		case msg, ok := <-c.ch:
			if !ok {
				return
			}
			err := c.conn.WriteJSON(msg)
			c.errCh <- err
		case <-c.stopCh:
			log.Info("stop ws client", zap.String("id", c.id), zap.Int64("user_id", c.userID))
			return
		}
	}
}

func (c *Client) Send(ctx context.Context, msg *Msg) error {
	c.mu.Lock()
	if c.stopped {
		return ErrClientClosed
	}
	defer c.mu.Unlock()
	var err error
	select {
	case c.ch <- msg:
		err = <-c.errCh
		if err != nil {
			log.Error("send", zap.String("id", c.id), zap.Int64("user_id", c.userID), zap.Error(err))
		}
	case <-ctx.Done():
		log.Error("send timeout", zap.String("id", c.id), zap.Int64("user_id", c.userID))
		close(c.ch)
		err = errors.New("timeout")
	}

	return err
}
