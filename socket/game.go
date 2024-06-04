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
	g.Broadcast([]byte("joined:" + client.Name + "," + client.Id))
}

func (g *Game) clientList(client Client) {
	list := "list:"
	for _, c := range g.clients {
		if c.Id != client.Id {
			list += c.Name + "," + c.Id + ";"
		}
	}
	client.send <- []byte("clients:" + list)
}

func (g *Game) Start() {
	g.Broadcast([]byte("game:started"))
}

func (g *Game) End() {
	g.Broadcast([]byte("game:ended"))
	for i, client := range g.clients {
		g.clients = append(g.clients[:i], g.clients[i+1:]...)
		client.send <- []byte("leave")
		client.Close()
	}
}

func (g *Game) Close() {
	for _, client := range g.clients {
		client.Close()
	}
}

func (g *Game) clientName(id string) string {
	for _, client := range g.clients {
		if client.Id == id {
			return client.Name
		}
	}
	return ""
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
			} else {
				switch string(action) {
				case "ready":
					g.ownerMessage([]byte("turn:now"))
				case "played":
					g.Broadcast([]byte("played:" + g.clientName(id) + ":" + string(message)))
				case "turn":
					g.playerMessage([]byte("turn:now"), string(message))
				case "left":
					g.RemoveClient(string(message))
				case "dos":
					g.Broadcast([]byte("dos:" + string(message)))
				case "card":
					if string(message) == "reset" {
						g.Broadcast([]byte("count:reset"))
					} else {
						g.Broadcast([]byte("pulled:" + g.clientName(id)))
					}
				case "color":
					g.Broadcast([]byte("color:" + string(message)))
				case "win":
					g.Broadcast([]byte("won:" + string(message)))
				default:
					log.Printf("Invalid action received: %s\n", action)
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

func (g *Game) ownerMessage(message []byte) {
	for _, client := range g.clients {
		if client.Id == g.Info.Owner {
			log.Printf("Sending message to client %s: %s\n", client.Name, message)
			client.send <- message
		}
	}
}

func (g *Game) playerMessage(message []byte, name string) {
	for _, client := range g.clients {
		if client.Name == name {
			log.Printf("Sending message to client %s: %s\n", client.Name, message)
			client.send <- message
		}
	}
}

func (g *Game) RemoveClient(id string) {
	for i, c := range g.clients {
		if c.Id == id {
			g.clients = append(g.clients[:i], g.clients[i+1:]...)
			c.send <- []byte("leave")
			c.Close()
			g.Broadcast([]byte("left:" + c.Name))
			break
		}
	}
}
