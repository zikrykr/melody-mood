package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/melody-mood/internal/playlist/port"
)

type playlistRoutes struct{}

var Routes playlistRoutes

func (r playlistRoutes) NewRoutes(router *gin.RouterGroup, playlistHandler port.IPlaylistHandler) {
	// (POST /api/v1/playlists)
	router.POST("", playlistHandler.GeneratePlaylists)
	// (POST /api/v1/playlists/spotify)
	router.POST("/spotify", playlistHandler.CreateUserSpotifyPlaylist)
}
