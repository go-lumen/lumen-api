package server

import (
	"database/sql"
	"github.com/adrien3d/lumen-api/services"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/spf13/viper"
)

// API structure that holds various necessary services
type API struct {
	Router          *gin.Engine
	Config          *viper.Viper
	MongoDatabase   *mgo.Database
	PostgreDatabase *sql.DB
	EmailSender     services.EmailSender
	TextSender      services.TextSender
}
