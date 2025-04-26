package payload

type PlaylistTrack struct {
	SongName        string `json:"song_name"`
	SongArtist      string `json:"song_artist"`
	SongAlbum       string `json:"song_album"`
	ReleaseDate     string `json:"release_date"`
	SpotifyCoverArt string `json:"spotify_cover_art"`
	SpotifyTrackID  string `json:"spotify_track_id"`
	SpotifyTrackURI string `json:"spotify_track_uri"`
}

type GeneratePlaylistResp struct {
	Tracks []PlaylistTrack `json:"tracks"`
}
