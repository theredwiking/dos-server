package game

import (
	"net/http"
)

func Routes() *http.ServeMux{
	router := http.NewServeMux()

	router.HandleFunc("GET /create", create)
	router.HandleFunc("GET /active", activeGames)
	router.HandleFunc("/join/{code}", joinGame)

	return router
}
