package main

import (
	"log"
	"net/http"

	"github.com/theredwiking/dos-server/game"
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("GET /status", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running"))
	})

	router.HandleFunc("GET /game/create", game.Create)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
