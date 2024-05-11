package player

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

const access_token = "BQDi7se9IT8s2zhbAecm2429nsNmZxzym7HJXUOUc2zIS_BhdX4gg6wNy65KuFnogWbozQ5JTBlnoXnXnJJxfDThQ6GtvY6ixh6mnODG1aret8nlXg_msS6pM2pMW2kZg2T3iHQLISx528MRBr0pFb4hs1SBlaASWzqljamhD_BQF--6OQCN6NLzuewD0fNL0MZidpb9zGePogQ"

type PlaybackState string

const (
	PLAYBACK_STATE_PLAYING  PlaybackState = "playing"
	PLAYBACK_STATE_PAUSED   PlaybackState = "paused"
	PLAYBACK_STATE_INACTIVE PlaybackState = "inactive"
)

type InactivePlayer struct {
	State PlaybackState `json:"state"`
}

type ActivePlayer struct {
	State      PlaybackState `json:"state"`
	Name       string        `json:"name"`
	Artist     string        `json:"artists"`
	DurationMs int           `json:"durationMs"`
	ProgressMs int           `json:"progressMs"`
}

type RawArtist struct {
	Name string `json:"name"`
}

type RawItem struct {
	Name       string      `json:"name"`
	Artists    []RawArtist `json:"artists"`
	DurationMs int         `json:"duration_ms"`
}

type RawPlayer struct {
	IsPlaying  bool    `json:"is_playing"`
	Item       RawItem `json:"item"`
	ProgressMs int     `json:"progress_ms"`
}

func GetPlayer(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/player", nil)

	authHeader := "Bearer " + access_token

	req.Header.Add("Authorization", authHeader)

	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

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
