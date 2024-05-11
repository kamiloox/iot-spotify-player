package main

import (
	"api/auth"
	"api/player"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/auth/callback", auth.Callback)
	http.HandleFunc("/auth/login", auth.Authorize)
	http.HandleFunc("/playback", player.GetPlayer)
	http.ListenAndServe(":8080", nil)
}
