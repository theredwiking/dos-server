package game

import (
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/google/uuid"
)

type Game struct {
	Id uint32 `json:"id"`
	Code string `json:"code"`
	Owner string `json:"owner"`
}

func Create(w http.ResponseWriter, r *http.Request) {
	game := Game{
		Id: uuid.New().ID(),
		Code: GenerateCode(),
		Owner: "John",
	}

	jsonData, err := json.Marshal(game)
	if err != nil {
		http.Error(w, "Error creating game", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

func GenerateCode() string {
	possibleChars := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGLMNOPQRSTUVWXYZ"
	code := ""
	for i := 0; i < 6; i++ {
		code += string(possibleChars[rand.Intn(len(possibleChars))])
	}
	return code
}
