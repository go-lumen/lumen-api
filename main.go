package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-lumen/lumen-api/migrations"
	"go-lumen/lumen-api/server"
	"go-lumen/lumen-api/services"
	"go-lumen/lumen-api/utils"
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
	dbType := api.Config.GetString("db_type")
	switch dbType {
	case "mongo":
		session, err := api.SetupMongoDatabase()
		utils.CheckErr(err)
		defer session.Close()

		err = api.SetupMongoIndexes()
		utils.CheckErr(err)

		// Seeds setup
		err = api.SetupMongoSeeds()
		utils.CheckErr(err)

	case "postgresql":
		db, err := api.SetupPostgreDatabase()
		utils.CheckErr(err)
		defer db.Close()

		err = api.SetupPostgreSeeds()
		utils.CheckErr(err)

		migrator := migrations.New(api)
		err = migrator.Migrate()
		if err != nil {
			fmt.Println("migration error", err)
		}
		//utils.CheckErr(err)

	case "mysql":
		db, err := api.SetupMySQLDatabase()
		utils.CheckErr(err)
		defer db.Close()

		err = api.SetupMySQLSeeds()
		utils.CheckErr(err)
	}

	// Router setup
	api.SetupRouter()
	err = api.Router.Run(api.Config.GetString("host_address"))
	utils.CheckErr(err)
}
