package game

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/theredwiking/dos-server/socket"
)

type ActiveGame map[string]socket.GameInfo

var activeGameInfo = make(ActiveGame)

func addActiveGame(game socket.GameInfo) {
	activeGameInfo[game.Code] = game
}

func removeActiveGame(code string) {
	delete(activeGameInfo, code)
}

func activeGames(w http.ResponseWriter, r *http.Request) {
	jsonData, err := json.Marshal(activeGameInfo)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error reading list of active games", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
