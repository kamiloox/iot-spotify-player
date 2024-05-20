package main

import (
	"api/auth"
	"api/database"
	"api/player"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Println(err)
	}

	conn := database.Connect()

	http.HandleFunc("/auth/login", auth.Authorize)
	http.HandleFunc("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		auth.Callback(w, r, conn)
	})
	http.HandleFunc("/playback", func(w http.ResponseWriter, r *http.Request) {
		player.Playback(w, r, conn)
	})

	http.ListenAndServe(getAppPort(), nil)
}

func getAppPort() string {
	if port, ok := os.LookupEnv("PORT"); ok {
		return ":" + port
	}

	return ":http"
}
