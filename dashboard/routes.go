package dashboard

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

func renderHtml(file string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "dashboard/" + file)
	})
}

// Routes returns a http.Handler that serves the dashboard routes
func Routes(app *firebase.App) http.Handler {
	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}
	router := http.NewServeMux()

	router.Handle("/login", renderHtml("login.html"))

	router.Handle("/", AuthCheck(renderHtml("index.html"), client))

	return router
}
