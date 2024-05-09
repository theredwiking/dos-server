package game

import (
	"log"
	"net/http"
	"strconv"

	"github.com/theredwiking/dos-server/socket"
)

type Games map[uint32]*socket.Game

var gameList = make(Games)

func addGame(game socket.GameInfo) {
	gameList[game.Id] = socket.NewGame(game)
	addActiveGame(game)
}

func removeGame(id uint32) {
	delete(gameList, id)
	removeActiveGame(id)
}

func joinGame(w http.ResponseWriter, r *http.Request) {
	gameId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid game id", http.StatusBadRequest)
		return
	}

	game := gameList[uint32(gameId)]
	if game == nil {
		http.Error(w, "Game not found", http.StatusNotFound)
		log.Println("Game not found:", gameId)
		return
	}
	
	if game.IsFull() {
		http.Error(w, "Game is full", http.StatusForbidden)
		log.Println("Game is full:", gameId)
		return
	}

	client := socket.NewClient(w, r)
	if client == nil {
		return
	}

	game.AddClient(*client)
}
