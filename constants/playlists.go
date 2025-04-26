package constants

const (
	PLAYLIST_CACHE_KEY = "session:%s:playlist"
)

const GeneratePlaylistPrompt = `You are a musical expert, and I'm looking for new playlist to listen.
Based on this picked songs: %s and the genre: %s

Can you arrange me 20 songs (including the picked songs) that might fit my personality and make sure that the song available on streaming platforms. 
Please only response in JSON RAW array, not wrapped in any object and without json markdown tag. Format:
[
  {
    "song_name": "...",
	"song_artist": "..."
  },
  ...
]`
