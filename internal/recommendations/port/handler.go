package port

import "github.com/gin-gonic/gin"

type IRecommendationHandler interface {
	// (GET /v1/recommendations)
	GenerateRecommendations(c *gin.Context)
}
