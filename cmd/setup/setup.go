package setup

import (
	"github.com/melody-mood/config"
	recommendationHandler "github.com/melody-mood/internal/recommendations/handler"
	recommendationPort "github.com/melody-mood/internal/recommendations/port"
	recommendationService "github.com/melody-mood/internal/recommendations/service"
	openai "github.com/sashabaranov/go-openai"
)

type SetupData struct {
	ConfigData  config.Config
	InternalApp InternalAppStruct
}

type InternalAppStruct struct {
	Services initServicesApp
	Handler  InitHandlerApp
}

// Services
type initServicesApp struct {
	RecommendationService recommendationPort.IRecommendationService
}

// Handler
type InitHandlerApp struct {
	RecommendationHandler recommendationPort.IRecommendationHandler
}

// CloseDB close connection to db
var CloseDB func() error

func InitSetup() SetupData {
	configData := config.GetConfig()
	internalAppVar := initInternalApp()

	return SetupData{
		ConfigData:  configData,
		InternalApp: internalAppVar,
	}
}

func initInternalApp() InternalAppStruct {
	var internalAppVar InternalAppStruct

	openAIClient := InitOpenAIService()

	initAppService(&internalAppVar, openAIClient)
	initAppHandler(&internalAppVar)

	return internalAppVar
}

func initAppService(initializeApp *InternalAppStruct, openAIClient *openai.Client) {
	initializeApp.Services.RecommendationService = recommendationService.NewRecommendationService(openAIClient)
}

func initAppHandler(initializeApp *InternalAppStruct) {
	initializeApp.Handler.RecommendationHandler = recommendationHandler.NewRecommendationHandler(initializeApp.Services.RecommendationService)
}
