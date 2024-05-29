package api

import (
	"context"
	"log"
	"net/http"

	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/iterator"
)

func UserList(w http.ResponseWriter, r *http.Request, client *auth.Client) {
	iter := client.Users(context.Background(), "")
	for {
		user, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Printf("error listing users: %v\n", err)
		}

		log.Printf("User: %s\n", user.Email)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Users listed"))
}
