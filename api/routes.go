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
		cookie, err := r.Cookie("idToken")
		if err != nil {
			http.Redirect(w, r, "/dashboard/login", http.StatusTemporaryRedirect)
			return
		}
		
		if cookie == nil {
			http.Redirect(w, r, "/dashboard/login", http.StatusTemporaryRedirect)
			return
		}
		
		_, err = client.VerifyIDToken(r.Context(), cookie.Value)
		if err != nil {
			http.Error(w, "error verifying token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Routes(app *firebase.App) *http.ServeMux{
	client, err := app.Auth(context.Background())	
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	router := http.NewServeMux()

	//router.HandleFunc("GET /create", create)
	//router.HandleFunc("GET /active", activeGames)
	//router.HandleFunc("/join/{code}", joinGame)
	router.Handle("GET /create", AuthCheck(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {create(w, r)}), client))
	router.Handle("GET /active", AuthCheck(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {activeGames(w, r)}), client))
	router.Handle("GET /join/{code}", AuthCheck(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {joinGame(w, r)}), client))

	return router
}
