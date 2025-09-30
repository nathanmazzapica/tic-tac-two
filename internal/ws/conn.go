package ws

import (
	"github.com/gorilla/websocket"
	"github.com/nathanmazzapica/tic-tac-two/internal/lobby"
	"time"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	id   string
	conn *websocket.Conn
	send chan<- lobby.Command
	read <-chan lobby.Event
}

func New() *Client {
	return &Client{}
}

func (c *Client) SetReadChan(read <-chan lobby.Event) {
	c.read = read
}
