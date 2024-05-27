package main

import (
	"context"
	"log"
	"net/http"

	firebase "firebase.google.com/go/v4"
	"github.com/theredwiking/dos-server/api"
	"github.com/theredwiking/dos-server/dashboard"
)

func main() {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v\n", err)
	}

	router := http.NewServeMux()
	router.HandleFunc("GET /status", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running"))
	})

	dashboard := dashboard.Routes(app)
	router.Handle("/dashboard/", http.StripPrefix("/dashboard", dashboard))

	gameRoutes := api.Routes(app)
	router.Handle("/game/", http.StripPrefix("/game", gameRoutes))


	server := http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	log.Println("Server is running on port 3000")
	log.Fatal(server.ListenAndServe())
}
