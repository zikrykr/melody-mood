package constants

const (
	RECOMMENDATION_CACHE_KEY = "session:%s:recommendations:%s:%s:%s" // session_id, personality, genre, occasion
)

const GeneateRecommendationPrompt = `You are a musical expert, and I'm looking for new playlist to listen.
My personality is %s, Mostly I listen to %s genre. And I usually listen to music on %s occasion.

Can you arrange me 5 songs that might fit my personality and make sure that the song available on streaming platforms. 
Not wrapped in any object and without json markdown tag. Format:
[
  {
    "song_name": "...",
	"song_artist": "...",
    "brief_reason": "..."
  },
  ...
]`

const (
	SPOTIFY_SEARCH_Q = "track:%s genre:%s"
)
