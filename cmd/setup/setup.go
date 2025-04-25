package setup

import (
	"github.com/melody-mood/config"
	"github.com/openai/openai-go"
)

type SetupData struct {
	ConfigData   config.Config
	InternalApp  InternalAppStruct
	OpenAIClient *openai.Client
}

type InternalAppStruct struct {
	Services initServicesApp
	Handler  InitHandlerApp
}

// Services
type initServicesApp struct {
}

// Handler
type InitHandlerApp struct {
}

// CloseDB close connection to db
var CloseDB func() error

func InitSetup() SetupData {
	configData := config.GetConfig()

	internalAppVar := initInternalApp()

	openAPIClient := InitOpenAIService()

	return SetupData{
		ConfigData:   configData,
		InternalApp:  internalAppVar,
		OpenAIClient: openAPIClient,
	}
}

func initInternalApp() InternalAppStruct {
	var internalAppVar InternalAppStruct

	initAppService(&internalAppVar)
	initAppHandler(&internalAppVar)

	return internalAppVar
}

func initAppService(initializeApp *InternalAppStruct) {
}

func initAppHandler(initializeApp *InternalAppStruct) {
}
