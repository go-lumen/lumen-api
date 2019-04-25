package server

import (
	"github.com/go-lumen/lumen-api/services"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/go-pg/pg"
	"github.com/spf13/viper"
)

// API structure that holds various necessary services
type API struct {
	Router          *gin.Engine
	Config          *viper.Viper
	MongoDatabase   *mgo.Database
	PostgreDatabase *pg.DB
	EmailSender     services.EmailSender
	TextSender      services.TextSender
}
