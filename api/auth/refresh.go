package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
)

func Refresh(r *http.Request, refreshToken string) Token {
	form := "grant_type=refresh_token&refresh_token=" + refreshToken + "&client_id=" + os.Getenv("SPOTIFY_CLIENT_ID")
	req, _ := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", strings.NewReader(form))

	authHeaderValue := os.Getenv("SPOTIFY_CLIENT_ID") + ":" + os.Getenv("SPOTIFY_CLIENT_SECRET")
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(authHeaderValue))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", authHeader)

	client := &http.Client{}
	res, _ := client.Do(req)

	body, _ := io.ReadAll(res.Body)

	var token Token

	json.Unmarshal(body, &token)

	return token
}

func FetchWithRefresh(req *http.Request, conn *pgx.Conn, boardToken string) (*http.Response, []byte) {
	var accessToken string
	var refreshToken string

	query := "SELECT spotify_access_token, spotify_refresh_token FROM auth WHERE board_secure_token=$1"
	conn.QueryRow(context.Background(), query, boardToken).Scan(&accessToken, &refreshToken)

	client := &http.Client{}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, _ := client.Do(req)

	if res.StatusCode == 401 {
		refreshed := Refresh(req, refreshToken)

		query = "UPDATE auth SET spotify_access_token=$1 WHERE board_secure_token=$2"

		req.Header.Set("Authorization", "Bearer "+refreshed.AccessToken)
	}

	res, _ = client.Do(req)
	body, _ := io.ReadAll(res.Body)

	return res, body
}
