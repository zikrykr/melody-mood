package port

import (
	"context"

	"github.com/melody-mood/internal/recommendations/payload"
)

type IRecommendationService interface {
	GenerateRecommendations(ctx context.Context, req payload.GenerateRecommendationsReq) ([]payload.RecommendationResponse, error)
}
