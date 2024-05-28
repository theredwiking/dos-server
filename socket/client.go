package socket

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	Id   string `json:"id"`
	Name string	`json:"name"`
	Conn *websocket.Conn
	send chan []byte
}

var upgrader = websocket.Upgrader{}

func NewClient(w http.ResponseWriter, r *http.Request) *Client {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error upgrading connection", http.StatusInternalServerError)
		return nil
	}

	return &Client{
		Id:   "",
		Name: "John" + strconv.FormatUint(uint64(uuid.New().ID()), 10),
		Conn: conn,
		send: make(chan []byte),
	}
}

func (c *Client) Read(messages chan []byte) {
	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		msg = []byte(c.Id + ": " + string(msg))
		messages <- msg
	}
}

func (c *Client) Write() {
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case msg := <-c.send:
			err := c.Conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println(err)
			}

		case <-ticker.C:
			log.Println("Sending ping to client", c.Id)
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func (c *Client) Close() {
	c.Conn.Close()
}
