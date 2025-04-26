package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/melody-mood/internal/playlist/payload"
	"github.com/melody-mood/internal/playlist/port"
	"github.com/melody-mood/pkg"
)

type PlaylistHandler struct {
	playlistService port.IPlaylistService
}

func NewPlaylistHandler(service port.IPlaylistService) port.IPlaylistHandler {
	return PlaylistHandler{
		playlistService: service,
	}
}

func (h PlaylistHandler) GeneratePlaylists(c *gin.Context) {
	var req payload.GeneratePlaylistReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	ctx := c.Request.Context()

	sessionID := c.GetHeader("X-Session-ID")
	if sessionID == "" {
		pkg.ResponseError(c, http.StatusBadRequest, fmt.Errorf("Missing X-Session-ID"))
		return
	}

	req.SessionID = sessionID

	resp, err := h.playlistService.GeneratePlaylists(ctx, req)
	if err != nil {
		pkg.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, pkg.HTTPResponse{
		Success: true,
		Message: "Playlist generated successfully",
		Data:    resp,
	})
}

func (h PlaylistHandler) CreateUserSpotifyPlaylist(c *gin.Context) {
	var req payload.CreateUserSpotifyPlaylistReq
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	ctx := c.Request.Context()

	sessionID := c.GetHeader("X-Session-ID")
	if sessionID == "" {
		pkg.ResponseError(c, http.StatusBadRequest, fmt.Errorf("Missing X-Session-ID"))
		return
	}

	req.SessionID = sessionID

	err := h.playlistService.CreateUserSpotifyPlaylist(ctx, req)
	if err != nil {
		pkg.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, pkg.HTTPResponse{
		Success: true,
		Message: "Spotify playlist created successfully",
	})
}
