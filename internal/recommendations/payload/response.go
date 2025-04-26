package payload

type RecommendationResponse struct {
	SpotifyTrackID  string `json:"spotify_track_id"`
	SongName        string `json:"song_name"`
	SongArtist      string `json:"song_artist"`
	SongAlbum       string `json:"song_album"`
	ReleaseDate     string `json:"release_date"`
	SpotifyCoverArt string `json:"spotify_cover_art"`
	BriefReason     string `json:"brief_reason"`
}
