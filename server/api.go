package server

import (
	"github.com/adrien3d/base-api/services"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/spf13/viper"
)

type API struct {
	Router      *gin.Engine
	Config      *viper.Viper
	Database    *mgo.Database
	EmailSender services.EmailSender
	TextSender  services.TextSender
}
