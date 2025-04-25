package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/melody-mood/internal/recommendations/payload"
	"github.com/melody-mood/internal/recommendations/port"
	"github.com/melody-mood/pkg"
)

type RecommendationHandler struct {
	recommendationService port.IRecommendationService
}

func NewRecommendationHandler(service port.IRecommendationService) port.IRecommendationHandler {
	return RecommendationHandler{
		recommendationService: service,
	}
}

func (h RecommendationHandler) GenerateRecommendations(c *gin.Context) {
	var req payload.GenerateRecommendationsReq
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

	resp, err := h.recommendationService.GenerateRecommendations(ctx, req)
	if err != nil {
		pkg.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, pkg.HTTPResponse{
		Success: true,
		Message: "Recommendations retrieved successfully",
		Data:    resp,
	})
}
