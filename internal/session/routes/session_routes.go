package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/melody-mood/internal/session/port"
)

type sessionRoutes struct{}

var Routes sessionRoutes

func (r sessionRoutes) NewRoutes(router *gin.RouterGroup, sessionHandler port.ISessionHandler) {
	// (POST /api/v1/sessions)
	router.POST("", sessionHandler.GenerateSessionID)
	// (GET /api/v1/sessions/auth/spotify)
	router.GET("/auth/spotify", sessionHandler.GenerateAuthSpotifyURL)
}
