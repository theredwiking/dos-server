package api

import (
	"log"
	"net/http"

	"github.com/theredwiking/dos-server/socket"
)

type Games map[string]*socket.Game

var gameList = make(Games)

func addGame(game socket.GameInfo) {
	newGame := socket.NewGame(game)
	go newGame.ReadMessages()
	gameList[game.Code] = newGame
}

func removeGame(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	delete(gameList, code)
}

func joinGame(w http.ResponseWriter, r *http.Request) {
	gameId := r.PathValue("code")

	game := gameList[gameId]
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
