package socket

import (
	"log"
	"math/rand"

	"github.com/gorilla/websocket"
)

type Client struct {
	id   int
	name string
	conn *websocket.Conn
	write chan []byte
}

func NewClient(conn *websocket.Conn, name string) *Client {
	return &Client{
		id:   rand.Int(),
		name: name,
		conn: conn,
	}
}

func (c *Client) Read(messages chan []byte) {
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		messages <- msg
	}
}

func (c *Client) Write() {
	for {
		select {
		case msg := <-c.write:
			err := c.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}
