package main

import (
	"github.com/melody-mood/cmd/rest"
	appSetup "github.com/melody-mood/cmd/setup"
	"github.com/melody-mood/config"
)

func main() {
	// config init
	config.InitConfig()

	// app setup init
	setup := appSetup.InitSetup()

	// starting REST server
	rest.StartServer(setup)
}
