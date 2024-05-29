package socket

import (
	"log"
)

type GameInfo struct {
	Code string `json:"code"`
	Owner string `json:"owner"`
}

type Game struct {
	Info GameInfo `json:"info"`
	clients []Client
	Connections uint32 `json:"connections"`
	messages chan []byte
}

func NewGame(game GameInfo) *Game {
	return &Game{
		Info: game,
		clients: []Client{},
		Connections: 0,
		messages: make(chan []byte),
	}
}

func (g *Game) AddClient(client Client) {
	g.clients = append(g.clients, client)
	go client.Write()
	go client.Read(g.messages)
	g.Connections++
	g.clientList(client)
	g.Broadcast([]byte("joined:" + client.Name))
}

func (g *Game) clientList(client Client) {
	list := "list:"
	for _, c := range g.clients {
		if c.Id != client.Id {
			list += c.Name + ","
		}
	}
	client.send <- []byte("clients:" + list)
}

func (g *Game) Start() {
	g.Broadcast([]byte("game:started"))
}

func (g *Game) End() {
	g.Broadcast([]byte("game:ended"))
}

func (g *Game) Close() {
	for _, client := range g.clients {
		client.Close()
	}
}

func splitMessage(message []byte) (string, []byte, []byte) {
	id := ""
	for i, byte := range message {
		if byte == ':' {
			id := string(message[:i])
			for j, byte := range message[i+1:] {
				if byte == ':' {
					return id, message[i+1:i+j+1], message[i+j+2:]
				}
			}
		}
	}
	return id, []byte{}, message
}

func (g *Game) ReadMessages() {
	for {
		select {
		case msg := <-g.messages:
			id, action, message := splitMessage(msg)
			if id == "" {
				log.Println("Invalid message received")
				continue
			}
			if id == g.Info.Owner && string(action) == "game" {
				switch string(message) {
				case "start":
					g.Start()
				case "end":
					g.End()
				}
			}
			log.Printf("Message received from client %s: %s\n", id, message)
		}
	}
}

func (g *Game) IsFull() bool {
	return len(g.clients) >= 10
}

func (g *Game) Broadcast(message []byte) {
	for _, client := range g.clients {
		log.Printf("Broadcasting message to client %s: %s\n", client.Id, message)
		client.send <- message
	}
}

func (g *Game) RemoveClient(client Client) {
	for i, c := range g.clients {
		if c.Id == client.Id && g.Info.Owner != client.Id {
			g.clients = append(g.clients[:i], g.clients[i+1:]...)
			client.Close()
			g.Broadcast([]byte("left:" + client.Name))
			break
		}
	}
}
