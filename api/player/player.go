package player

import (
	"api/auth"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func Playback(w http.ResponseWriter, r *http.Request, conn *pgx.Conn) {
	boardToken := r.URL.Query().Get("token")

	req, _ := http.NewRequest("GET", "https://api.spotify.com/v1/me/player", nil)

	res, body := auth.FetchWithRefresh(req, conn, boardToken)

	var rawPlayer RawPlayer

	_ = json.Unmarshal(body, &rawPlayer)

	if res.StatusCode == 204 {
		var player InactivePlayer

		player.State = PLAYBACK_STATE_INACTIVE

		output, _ := json.Marshal(player)
		w.Write(output)

		return
	}

	if res.StatusCode != 200 {
		w.WriteHeader(res.StatusCode)
		w.Write(body)

		return
	}

	var player ActivePlayer

	if rawPlayer.IsPlaying {
		player.State = PLAYBACK_STATE_PLAYING
	} else {
		player.State = PLAYBACK_STATE_PAUSED
	}

	player.Name = rawPlayer.Item.Name
	player.Artist = rawPlayer.Item.Artists[0].Name
	player.DurationMs = rawPlayer.Item.DurationMs
	player.ProgressMs = rawPlayer.ProgressMs

	w.Header().Set("Content-Type", "application/json")

	output, _ := json.Marshal(player)
	w.Write(output)
}
