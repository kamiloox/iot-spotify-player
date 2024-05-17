package main

import (
	"api/auth"
	"api/database"
	"api/player"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	dotEnvError := godotenv.Load("../.env")

	if dotEnvError != nil {
		log.Fatal("Error loading .env file")
	}

	conn := database.Connect()

	http.HandleFunc("/auth/login", auth.Authorize)
	http.HandleFunc("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		auth.Callback(w, r, conn)
	})
	http.HandleFunc("/playback", func(w http.ResponseWriter, r *http.Request) {
		player.Playback(w, r, conn)
	})

	http.ListenAndServe("127.0.0.1:8080", nil)
}
