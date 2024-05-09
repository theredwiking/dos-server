package socket

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	Id   uint32 `json:"id"`
	Name string	`json:"name"`
	Conn *websocket.Conn
	Writer chan []byte
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
		Id:   uuid.New().ID(),
		Name: "John",
		Conn: conn,
	}
}

func (c *Client) Read(messages chan []byte) {
	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("Message received from client %d: %s\n", c.Id, msg)
		messages <- msg
	}
}

func (c *Client) Write() {
	for {
		select {
		case msg := <-c.Writer:
			err := c.Conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func (c *Client) Close() {
	c.Conn.Close()
}
