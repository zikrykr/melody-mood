package port

import "github.com/gin-gonic/gin"

type ICallbackHandler interface {
	// (GET /v1/callback)
	HandleSpotifyCallback(c *gin.Context)
}
