package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func Authorize(w http.ResponseWriter, r *http.Request) {
	req, _ := http.NewRequest("GET", "https://accounts.spotify.com/authorize", nil)

	q := req.URL.Query()
	q.Add("response_type", "code")
	q.Add("client_id", os.Getenv("SPOTIFY_CLIENT_ID"))
	q.Add("scope", os.Getenv("SPOTIFY_SCOPE"))
	q.Add("redirect_uri", os.Getenv("SPOTIFY_REDIRECT_URI"))
	q.Add("state", r.URL.Query().Get("token"))
	req.URL.RawQuery = q.Encode()

	http.Redirect(w, r, req.URL.String(), http.StatusTemporaryRedirect)
}

func Callback(w http.ResponseWriter, r *http.Request, conn *pgx.Conn) {
	code := r.URL.Query().Get("code")
	boardToken := r.URL.Query().Get("state")

	req, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token", nil)

	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(os.Getenv("SPOTIFY_CLIENT_ID")+":"+os.Getenv("SPOTIFY_CLIENT_SECRET")))

	form := req.URL.Query()
	form.Add("grant_type", "authorization_code")
	form.Add("code", code)
	form.Add("redirect_uri", os.Getenv("SPOTIFY_REDIRECT_URI"))

	req.Header.Add("Authorization", authHeader)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = form.Encode()

	client := &http.Client{}
	res, _ := client.Do(req)

	body, _ := io.ReadAll(res.Body)

	var token Token
	json.Unmarshal(body, &token)

	query := "INSERT INTO auth (board_secure_token, spotify_access_token, spotify_refresh_token) VALUES ($1, $2, $3) ON CONFLICT (board_secure_token) DO UPDATE SET spotify_access_token=$2, spotify_refresh_token=$3"
	conn.Exec(context.Background(), query, boardToken, token.AccessToken, token.RefreshToken)

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Credentials saved for token: " + boardToken + "\nYou can close this tab now."))
}
