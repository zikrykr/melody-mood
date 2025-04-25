package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/melody-mood/internal/recommendations/port"
)

type recommendationRoutes struct{}

var Routes recommendationRoutes

func (r recommendationRoutes) NewRoutes(router *gin.RouterGroup, recommendationHandler port.IRecommendationHandler) {
	// (POST /api/v1/recommendations)
	router.POST("/", recommendationHandler.GenerateRecommendations)
}
