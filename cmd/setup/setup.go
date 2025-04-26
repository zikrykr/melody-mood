package setup

import (
	"github.com/melody-mood/config"
	recommendationHandler "github.com/melody-mood/internal/recommendations/handler"
	recommendationPort "github.com/melody-mood/internal/recommendations/port"
	recommendationService "github.com/melody-mood/internal/recommendations/service"
	sessionHandler "github.com/melody-mood/internal/session/handler"
	sessionPort "github.com/melody-mood/internal/session/port"
	sessionService "github.com/melody-mood/internal/session/service"
	redis "github.com/redis/go-redis/v9"
	openai "github.com/sashabaranov/go-openai"
)

type SetupData struct {
	ConfigData  config.Config
	InternalApp InternalAppStruct
}

type InternalAppStruct struct {
	Services initServicesApp
	Handler  InitHandlerApp

	RedisClient *redis.Client
}

// Services
type initServicesApp struct {
	RecommendationService recommendationPort.IRecommendationService
	SessionService        sessionPort.ISessionService
}

// Handler
type InitHandlerApp struct {
	RecommendationHandler recommendationPort.IRecommendationHandler
	SessionHandler        sessionPort.ISessionHandler
}

func InitSetup() SetupData {
	configData := config.GetConfig()
	redisC := InitRedis()
	internalAppVar := initInternalApp(redisC)

	return SetupData{
		ConfigData:  configData,
		InternalApp: internalAppVar,
	}
}

func initInternalApp(redis *redis.Client) InternalAppStruct {
	var internalAppVar InternalAppStruct

	openAIClient := InitOpenAIService()
	internalAppVar.RedisClient = redis

	initAppService(&internalAppVar, openAIClient)
	initAppHandler(&internalAppVar)

	return internalAppVar
}

func initAppService(initializeApp *InternalAppStruct, openAIClient *openai.Client) {
	initializeApp.Services.RecommendationService = recommendationService.NewRecommendationService(openAIClient, initializeApp.RedisClient)
	initializeApp.Services.SessionService = sessionService.NewSessionService(initializeApp.RedisClient)
}

func initAppHandler(initializeApp *InternalAppStruct) {
	initializeApp.Handler.RecommendationHandler = recommendationHandler.NewRecommendationHandler(initializeApp.Services.RecommendationService)
	initializeApp.Handler.SessionHandler = sessionHandler.NewSessionHandler(initializeApp.Services.SessionService)
}
