package socket

import (
	"log"

	"github.com/theredwiking/dos-server/game"
)

type GameInfo struct {
	Code string `json:"code"`
	Owner uint32 `json:"owner"`
}

type Game struct {
	Info GameInfo `json:"info"`
	clients []Client
	Connections uint32
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
	g.Broadcast([]byte(client.Name + " has joined the game"))
}

func (g *Game) Start() {
	g.Broadcast([]byte("Game starting!"))
}

func (g *Game) End() {
	g.Broadcast([]byte("Game ending!"))
}

func (g *Game) Close() {
	for _, client := range g.clients {
		client.Close()
	}
}

func (g *Game) ReadMessages() {
	for {
		select {
		case msg := <-g.messages:
			log.Printf("Message received from client: %s\n", msg)
		}
	}
}

func (g *Game) IsFull() bool {
	return len(g.clients) >= 10
}

func (g *Game) Broadcast(message []byte) {
	for _, client := range g.clients {
		log.Printf("Broadcasting message to client %d: %s\n", client.Id, message)
		client.send <- message
	}
}

func (g *Game) RemoveClient(client Client) {
	for i, c := range g.clients {
		if c.Id == client.Id && g.Info.Owner != client.Id {
			g.clients = append(g.clients[:i], g.clients[i+1:]...)
			client.Close()
			g.Broadcast([]byte(client.Name + " has left the game"))
			break
		}
	}
}
