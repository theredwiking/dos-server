package api

import (
	"context"
	"log"
	"net/http"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

func AuthCheck(next http.Handler, client *auth.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header["Authorization"]

		if token == nil {
			token = r.URL.Query()["token"]
			if token == nil {
				http.Error(w, "error getting auth header", http.StatusUnauthorized)
				return
			}
		}

		client, err := client.VerifyIDToken(r.Context(), token[0])
		if err != nil {
			http.Error(w, "error verifying token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", client)
		newReq := r.WithContext(ctx)

		next.ServeHTTP(w, newReq)
	})
}

func FirebaseRoutes(app *firebase.App) *http.ServeMux {
	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	router := http.NewServeMux()

	router.Handle("GET /users", AuthCheck(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {UserList(w, r, client)}), client))

	return router
}

func GameRoutes(app *firebase.App) *http.ServeMux{
	client, err := app.Auth(context.Background())	
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	router := http.NewServeMux()

	router.Handle("GET /create", AuthCheck(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {create(w, r)}), client))
	router.Handle("GET /active", AuthCheck(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {activeGames(w, r)}), client))
	router.Handle("GET /join/{code}", AuthCheck(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {joinGame(w, r)}), client))

	return router
}
