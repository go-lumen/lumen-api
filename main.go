package main

import (
	"github.com/go-lumen/lumen-api/server"
	"github.com/go-lumen/lumen-api/services"
	"github.com/go-lumen/lumen-api/utils"
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

	// Database setup
	session, err := api.SetupMongoDatabase()
	utils.CheckErr(err)
	defer session.Close()

	err = api.SetupMongoIndexes()
	utils.CheckErr(err)

	// Seeds setup
	err = api.SetupMongoSeeds()
	utils.CheckErr(err)

	/*db, err := api.SetupPostgreDatabase()
	utils.CheckErr(err)
	defer db.Close()

	err = api.SetupPostgreSeeds()
	utils.CheckErr(err)*/

	// Router setup
	api.SetupRouter()
	err = api.Router.Run(api.Config.GetString("host_address"))
	utils.CheckErr(err)
}
