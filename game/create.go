package game

import (
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/google/uuid"
	"github.com/theredwiking/dos-server/socket"
)


func create(w http.ResponseWriter, r *http.Request) {
	game := socket.GameInfo{
		Code: generateCode(),
		Owner: uuid.New().ID(),
	}

	jsonData, err := json.Marshal(game)
	if err != nil {
		http.Error(w, "Error creating game", http.StatusInternalServerError)
		return
	}

	addGame(game)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

func generateCode() string {
	possibleChars := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGLMNOPQRSTUVWXYZ"
	code := ""
	for i := 0; i < 6; i++ {
		code += string(possibleChars[rand.Intn(len(possibleChars))])
	}
	return code
}
