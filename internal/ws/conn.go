package ws

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/nathanmazzapica/tic-tac-two/internal/dto"
	"log"
	"time"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	ID     string
	conn   *websocket.Conn
	cmds   CommandSink
	events <-chan dto.Event
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewClient(id string, conn *websocket.Conn, sink CommandSink, sub <-chan dto.Event) *Client {
	return &Client{ID: id, conn: conn, cmds: sink, events: sub}
}

func (c *Client) Listen(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case evt := <-c.events:
			log.Println("evt")
			// marshall event into json
			// send through ws
			// simple.
			evtJson, err := json.Marshal(evt)
			// todo: figure out how to handle
			if err != nil {
				log.Printf("failed to marshal event: %s", err)
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, evtJson); err != nil {
				log.Printf("failed to send event: %s", err)
			}
		default:
		}
	}
}

func (c *Client) Send(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			_, msg, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure) {
					log.Printf("websocket connection closed unexpectedly: %s", err)
				}
				// todo: handle more properly
				return err
			}

			envelope := &dto.Envelope{}
			err = json.Unmarshal(msg, envelope)

			if err != nil {
				log.Printf("failed to unmarshal envelope: %s", err)
			}

			c.cmds.Post(envelope)

		}
	}
}
