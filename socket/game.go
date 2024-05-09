package socket

type GameInfo struct {
	Id uint32 `json:"id"`
	Code string `json:"code"`
	Owner uint32 `json:"owner"`
}

type Game struct {
	Info GameInfo `json:"info"`
	Clients []Client `json:"players"`
}

func NewGame(game GameInfo) *Game {
	return &Game{
		Info: game,
		Clients: []Client{},
	}
}

func (g *Game) AddClient(client Client) {
	g.Clients = append(g.Clients, client)
	g.OwnerMessage([]byte(client.Name + " has joined the game"))
}

func (g *Game) OwnerMessage(message []byte) {
	for _, client := range g.Clients {
		if client.Id == g.Info.Owner {
			client.Writer <- message
			break
		}
	}
}

func (g *Game) Broadcast(message []byte) {
	for _, client := range g.Clients {
		client.Writer <- message
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
