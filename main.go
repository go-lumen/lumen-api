package main

import (
	"github.com/adrien3d/base-api/server"
	"github.com/adrien3d/base-api/services"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	api := &server.API{Router: gin.Default(), Config: viper.New()}

	// Configuration setup
	err := api.SetupViper()
	if err != nil {
		panic(err)
	}

	// Email sender setup
	api.EmailSender = services.NewEmailSender(api.Config)
	api.TextSender = services.NewTextSender(api.Config)

	// Database setup
	session, err := api.SetupDatabase()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err = api.SetupIndexes()
	if err != nil {
		panic(err)
	}

	// Seeds setup
	api.SetupSeeds()

	// Router setup
	api.SetupRouter()
	api.Router.Run(api.Config.GetString("host_address"))
}
