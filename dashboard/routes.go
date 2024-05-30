package dashboard

import (
	"context"
	"embed"
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
			log.Printf("error verifying token: %v\n", err)
			http.Error(w, "error verifying token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func renderHtml(file string, fs embed.FS) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := fs.ReadFile("dashboard/" + file)
		if err != nil {
			http.Error(w, "error reading file", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
}

// Routes returns a http.Handler that serves the dashboard routes
func Routes(app *firebase.App, files embed.FS) http.Handler {
	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}
	router := http.NewServeMux()

	router.Handle("/login", renderHtml("login.html", files))

	router.Handle("/", AuthCheck(renderHtml("index.html", files), client))

	return router
}
