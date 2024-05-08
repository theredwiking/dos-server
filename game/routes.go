package game

import (
	"net/http"
)

func Routes() *http.ServeMux{
	router := http.NewServeMux()

	router.HandleFunc("/create", Create)

	return router
}
