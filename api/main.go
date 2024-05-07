package main

import (
	"api/auth"
	"encoding/json"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

type playbackState string

const (
	PLAYING       playbackState = "PLAYING"
	PAUSED        playbackState = "PAUSED"
	NOT_AVAILABLE playbackState = "NOT_AVAILABLE"
)

type Player struct {
	State playbackState `json:"state"`
}

func main() {
	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/auth/callback", auth.Callback)
	http.HandleFunc("/auth/login", auth.Authorize)
	http.ListenAndServe(":8080", nil)
}

func GetPlaybackState(w http.ResponseWriter, r *http.Request) {
	player := Player{State: NOT_AVAILABLE}

	res, err := json.Marshal(player)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
