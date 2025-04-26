package payload

type GeneratePlaylistReq struct {
	PickedSongs  []PickedSongReq `json:"picked_songs"`
	Genre        string          `json:"genre"`
	IsRegenerate bool            `json:"is_regenerate"`

	SessionID string
}

type CreateUserSpotifyPlaylistReq struct {
	Name      string `json:"name"`
	SessionID string
}

type PickedSongReq struct {
	SongName   string `json:"song_name"`
	SongArtist string `json:"song_artist"`
}
