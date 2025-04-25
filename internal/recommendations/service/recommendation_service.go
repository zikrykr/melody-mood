package service

import (
	"context"
	"fmt"

	"github.com/melody-mood/constants"
	"github.com/melody-mood/internal/recommendations/payload"
	"github.com/melody-mood/internal/recommendations/port"
	"github.com/melody-mood/pkg"
	openai "github.com/sashabaranov/go-openai"
)

type RecommendationService struct {
	openAIClient *openai.Client
}

func NewRecommendationService(openAIClient *openai.Client) port.IRecommendationService {
	return RecommendationService{
		openAIClient: openAIClient,
	}
}

func (r RecommendationService) GenerateRecommendations(ctx context.Context, req payload.GenerateRecommendationsReq) ([]payload.RecommendationResponse, error) {
	recommendationResp, err := pkg.CreateOpenAIMessage(ctx, r.openAIClient, fmt.Sprintf(constants.GeneateRecommendationPrompt, req.Personality, req.Genre, req.Occasion))
	if err != nil {
		return nil, err
	}

	recommendations, err := pkg.ParseToStruct[[]payload.RecommendationResponse](recommendationResp)
	if err != nil {
		return nil, err
	}
	return recommendations, nil
}
