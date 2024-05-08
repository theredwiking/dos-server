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

	gameRoutes := game.Routes()
	router.Handle("/game/", http.StripPrefix("/game", gameRoutes))


	server := http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	log.Println("Server is running on port 3000")
	log.Fatal(server.ListenAndServe())
}
