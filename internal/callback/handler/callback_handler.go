package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/melody-mood/internal/callback/port"
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
		c.Redirect(http.StatusFound, "https://melody-mood.com/oauth-spotify?status=error")
	}

	c.Redirect(http.StatusFound, "https://melody-mood.com/oauth-spotify?status=success")
}
