package auth

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"os"
)

func Authorize(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest("GET", "https://accounts.spotify.com/authorize", nil)

	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := req.URL.Query()
	q.Add("response_type", "code")
	q.Add("client_id", os.Getenv("SPOTIFY_CLIENT_ID"))
	q.Add("scope", os.Getenv("SPOTIFY_SCOPE"))
	q.Add("redirect_uri", os.Getenv("SPOTIFY_REDIRECT_URI"))
	req.URL.RawQuery = q.Encode()

	http.Redirect(w, r, req.URL.String(), http.StatusTemporaryRedirect)
}

func Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", nil)

	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(os.Getenv("SPOTIFY_CLIENT_ID")+":"+os.Getenv("SPOTIFY_CLIENT_SECRET")))

	form := req.URL.Query()
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)
	form.Add("redirect_uri", os.Getenv("SPOTIFY_REDIRECT_URI"))

	req.Header.Add("Authorization", authHeader)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = form.Encode()

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
