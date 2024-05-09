package game

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/theredwiking/dos-server/socket"
)

type Games map[uint32]*socket.Game

var gameList = make(Games)

func addGame(game socket.GameInfo) {
	gameList[game.Id] = socket.NewGame(game)
}

func removeGame(id uint32) {
	delete(gameList, id)
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

	client := socket.NewClient(w, r)
	if client == nil {
		return
	}

	game.AddClient(*client)
	log.Println("Client joined game: %d, Userid: %d", gameId, client.Id)
}

func activeGames(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting list of active games:", gameList)
	jsonData, err := json.Marshal(gameList)
	if err != nil {
		http.Error(w, "Error reading list of active games", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
