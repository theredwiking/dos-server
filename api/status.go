package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func activeGames(w http.ResponseWriter, r *http.Request) {
	jsonData, err := json.Marshal(gameList)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error reading list of active games", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
