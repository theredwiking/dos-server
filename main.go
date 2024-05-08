package main

import (
	"log"
	"net/http"

	"github.com/theredwiking/dos-server/game"
)

func main() {
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running"))
	})
	http.HandleFunc("/create", game.Create)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
