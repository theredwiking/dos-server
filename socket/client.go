package socket

import (
	"log"
	"net/http"
	"strings"
	"time"

	"firebase.google.com/go/v4/auth"
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
		Id:   r.Context().Value("user").(*auth.Token).UID,
		Name: strings.Split(r.Context().Value("user").(*auth.Token).Claims["email"].(string), "@")[0],
		Conn: conn,
		send: make(chan []byte),
	}
}

func (c *Client) Read(messages chan []byte) {
	for {
		_, msg, err := c.Conn.ReadMessage()
		left := false
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
				left = true
				msg = []byte("left:" + c.Id)
			} else {
				log.Printf("Client read error: %v", err)
				return
			}
		}
		msg = []byte(c.Id + ":" + string(msg))
		messages <- msg
		if left {
			return
		}
	}
}

func (c *Client) Write() {
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case msg := <-c.send:
			if string(msg) == "leave" {
				break
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("Client write error: %v", err)
			}

		case <-ticker.C:
			if err := c.Conn.WriteMessage(websocket.PingMessage, []byte("keepalive")); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func (c *Client) Close() {
	c.Conn.Close()
}
