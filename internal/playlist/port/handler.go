package port

import "github.com/gin-gonic/gin"

type IPlaylistHandler interface {
	// (GET /v1/playlists)
	GeneratePlaylists(c *gin.Context)
	// (POST /v1/playlists/spotify)
	CreateUserSpotifyPlaylist(c *gin.Context)
}
