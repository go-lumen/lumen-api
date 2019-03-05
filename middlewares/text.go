package middlewares

import (
	"github.com/adrien3d/base-api/services"
	"github.com/gin-gonic/gin"
)

func TextMiddleware(textSender services.TextSender) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("textSender", textSender)
		c.Next()
	}
}
