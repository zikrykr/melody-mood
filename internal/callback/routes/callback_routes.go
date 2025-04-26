package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/melody-mood/internal/callback/port"
)

type callbackRoutes struct{}

var Routes callbackRoutes

func (r callbackRoutes) NewRoutes(router *gin.RouterGroup, callbackHandler port.ICallbackHandler) {
	// (POST /api/v1/callback)
	router.GET("", callbackHandler.HandleSpotifyCallback)
}
