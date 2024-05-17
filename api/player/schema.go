package player

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
