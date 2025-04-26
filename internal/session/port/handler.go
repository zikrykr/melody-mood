package port

import "github.com/gin-gonic/gin"

type ISessionHandler interface {
	// (GET /v1/sessions)
	GenerateSessionID(c *gin.Context)
	// (GET /v1/sessions/auth/spotify)
	GenerateAuthSpotifyURL(c *gin.Context)
}
