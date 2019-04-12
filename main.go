package main

import (
	"github.com/adrien3d/lumen-api/server"
	"github.com/adrien3d/lumen-api/services"
	"github.com/adrien3d/lumen-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	api := &server.API{Router: gin.Default(), Config: viper.New()}

	// Configuration setup
	err := api.SetupViper()
	utils.CheckErr(err)

	// Email sender setup
	api.EmailSender = services.NewEmailSender(api.Config)
	api.TextSender = services.NewTextSender(api.Config)

	// Databases setup
	session, err := api.SetupMongoDatabase()
	utils.CheckErr(err)
	defer session.Close()

	db, err := api.SetupPostgreDatabase()
	utils.CheckErr(err)
	defer db.Close()

	err = api.SetupIndexes()
	utils.CheckErr(err)

	// Seeds setup
	api.SetupSeeds()

	// Router setup
	api.SetupRouter()
	api.Router.Run(api.Config.GetString("host_address"))
}
