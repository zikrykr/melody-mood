package port

import "github.com/gin-gonic/gin"

type IPlaylistHandler interface {
	// (GET /v1/playlists)
	GeneratePlaylists(c *gin.Context)
}
