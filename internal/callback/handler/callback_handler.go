package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/melody-mood/internal/callback/port"
	"github.com/melody-mood/pkg"
)

type CallbackHandler struct {
	callbackService port.ICallbackService
}

func NewCallbackHandler(service port.ICallbackService) port.ICallbackHandler {
	return CallbackHandler{
		callbackService: service,
	}
}

func (h CallbackHandler) HandleSpotifyCallback(c *gin.Context) {
	ctx := c.Request.Context()

	code := c.Query("code")
	errMsg := c.Query("error")
	state := c.Query("state")

	err := h.callbackService.HandleSpotifyCallback(ctx, code, errMsg, state)
	if err != nil {
		pkg.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, pkg.HTTPResponse{
		Success: true,
		Message: "Spotify callback handled successfully",
	})
}
