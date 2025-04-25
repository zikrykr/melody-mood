package service

import (
	"context"
	"fmt"
	"time"

	"github.com/melody-mood/constants"
	"github.com/melody-mood/internal/recommendations/payload"
	"github.com/melody-mood/internal/recommendations/port"
	"github.com/melody-mood/pkg"
	"github.com/redis/go-redis/v9"
	openai "github.com/sashabaranov/go-openai"
)

type RecommendationService struct {
	openAIClient *openai.Client
	rds          *redis.Client
}

func NewRecommendationService(openAIClient *openai.Client, rds *redis.Client) port.IRecommendationService {
	return RecommendationService{
		openAIClient: openAIClient,
		rds:          rds,
	}
}

func (r RecommendationService) GenerateRecommendations(ctx context.Context, req payload.GenerateRecommendationsReq) ([]payload.RecommendationResponse, error) {
	var (
		recommendations []payload.RecommendationResponse
		cacheKey        = fmt.Sprintf(constants.RECOMMENDATION_CACHE_KEY, req.SessionID, req.Personality, req.Genre, req.Occasion)
	)
	// Try read from cache
	cached, err := r.rds.Get(ctx, cacheKey).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if !req.IsRegenerate && cached != "" {
		recommendations, err = pkg.ParseToStruct[[]payload.RecommendationResponse](cached)
		if err != nil {
			return nil, err
		}
	} else {
		recommendationResp, err := pkg.CreateOpenAIMessage(ctx, r.openAIClient, fmt.Sprintf(constants.GeneateRecommendationPrompt, req.Personality, req.Genre, req.Occasion))
		if err != nil {
			return nil, err
		}

		// Save to Redis
		err = r.rds.Set(ctx, cacheKey, recommendationResp, 2*time.Hour).Err()
		if err != nil {
			return nil, err
		}

		recommendations, err = pkg.ParseToStruct[[]payload.RecommendationResponse](recommendationResp)
		if err != nil {
			return nil, err
		}
	}

	return recommendations, nil
}
