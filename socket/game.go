package socket

import "log"

type GameInfo struct {
	Code string `json:"code"`
	Owner uint32 `json:"owner"`
}

type Game struct {
	Info GameInfo `json:"info"`
	Clients []Client `json:"clients"`
}

func NewGame(game GameInfo) *Game {
	return &Game{
		Info: game,
		Clients: []Client{},
	}
}

func (g *Game) AddClient(client Client) {
	g.Clients = append(g.Clients, client)
	g.Broadcast([]byte(client.Name + " has joined the game"))
}

func (g *Game) OwnerMessage(message []byte) {
	for _, client := range g.Clients {
		if client.Id == g.Info.Owner {
			client.Write(message)
			break
		}
	}
}

func (g *Game) IsFull() bool {
	return len(g.Clients) >= 10
}

func (g *Game) Broadcast(message []byte) {
	for _, client := range g.Clients {
		log.Printf("Broadcasting message to client %d: %s\n", client.Id, message)
		client.Write(message)
	}
}

func (g *Game) RemoveClient(client Client) {
	for i, c := range g.Clients {
		if c.Id == client.Id && g.Info.Owner != client.Id {
			g.Clients = append(g.Clients[:i], g.Clients[i+1:]...)
			client.Close()
			g.Broadcast([]byte(client.Name + " has left the game"))
			break
		}
	}
}
