package middlewares

import (
	"github.com/adrien3d/base-api/config"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func ConfigMiddleware(viper *viper.Viper) gin.HandlerFunc {
	return func(c *gin.Context) {
		config.ToContext(c, config.New(viper))
		c.Next()
	}
}
