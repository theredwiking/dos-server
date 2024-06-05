package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/iterator"
)


func UserList(w http.ResponseWriter, r *http.Request, client *auth.Client) {
	userList := []auth.UserInfo{}

	iter := client.Users(context.Background(), "")
	for {
		user, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Printf("error listing users: %v\n", err)
		}
		userList = append(userList, *user.UserInfo)
	}
	
	users, err := json.Marshal(userList)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error reading list of users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(users)
}
