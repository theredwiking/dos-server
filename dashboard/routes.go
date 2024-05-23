package dashboard

import (
	"net/http"
	firebase "firebase.google.com/go/v4"
)

// Routes returns a http.Handler that serves the dashboard routes
func Routes(*firebase.App) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "dashboard/login.html")
	})

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "dashboard/index.html")
	})

	return router
}
