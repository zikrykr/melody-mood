package constants

const (
	OpenAIURL   = "https://api.openai.com/v1/chat/completions"
	OpenAIModel = "gpt-3.5-turbo"
)

const GeneateRecommendationPrompt = `You are a musical expert, and I'm looking for new playlist to listen.
My personality is %s, Mostly I listen to %s genre. And I usually listen to music on %s occasion.

Can you arrange me 5 songs that might fit my personality and don't put any mainstream songs. 
not wrapped in any object. Format:
[
  {
    "spotify_track_id": "...",
    "spotify_cover_art": "...",
    "song_name": "...",
    "song_artist": "...",
    "song_album": "...",
    "released_year": "1986",
    "brief_reason": "..."
  },
  ...
]`
