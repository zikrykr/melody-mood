package setup

import (
	"github.com/melody-mood/config"
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
}

// Handler
type InitHandlerApp struct {
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

	initAppService(&internalAppVar)
	initAppHandler(&internalAppVar)

	return internalAppVar
}

func initAppService(initializeApp *InternalAppStruct) {
}

func initAppHandler(initializeApp *InternalAppStruct) {
}
