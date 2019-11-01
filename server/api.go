package server

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"go-lumen/lumen-api/services"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

// API structure that holds various necessary services
type API struct {
	Router          *gin.Engine
	Config          *viper.Viper
	MongoDatabase   *mgo.Database
	PostgreDatabase *gorm.DB
	MySQLDatabase   *sql.DB
	EmailSender     services.EmailSender
	TextSender      services.TextSender
}
