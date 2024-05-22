package socket

import (
	"log"
	"strconv"
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
	g.Broadcast([]byte("joined:" + client.Name))
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

func splitMessage(message []byte) (uint32, []byte, []byte) {
	id := uint32(0)
	for i, byte := range message {
		if byte == ':' {
			id, err := strconv.ParseUint(string(message[:i]), 10, 32)
			if err != nil {
				log.Println(err)
				id = 0
			}
			for j, byte := range message[i+1:] {
				if byte == ':' {
					return uint32(id), message[i+1:i+j+1], message[i+j+2:]
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
			if id == 0 {
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
			log.Printf("Message received from client %d: %s\n", id, message)
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
			g.Broadcast([]byte("left:" + client.Name))
			break
		}
	}
}
