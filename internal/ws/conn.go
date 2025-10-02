package ws

import (
	"context"
	"github.com/gorilla/websocket"
	"time"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	id     string
	conn   *websocket.Conn
	cmds   CommandSink
	events <-chan any
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewClient(id string, conn *websocket.Conn, sink CommandSink, sub <-chan any) *Client {
	return &Client{id: id, conn: conn, cmds: sink, events: sub}
}

func (c *Client) Listen(ctx context.Context) error {
	return nil
}

func (c *Client) readPump(ctx context.Context) {

}
