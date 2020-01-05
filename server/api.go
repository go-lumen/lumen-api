package server

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/go-lumen/lumen-api/services"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

// API structure that holds various necessary services
type API struct {
	Router          *gin.Engine
	Config          *viper.Viper
	MongoDatabase   *mongo.Database
	PostgreDatabase *gorm.DB
	MySQLDatabase   *sql.DB
	EmailSender     services.EmailSender
	TextSender      services.TextSender
}
