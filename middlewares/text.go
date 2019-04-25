package middlewares

import (
	"github.com/go-lumen/lumen-api/services"
	"github.com/gin-gonic/gin"
)

// TextMiddleware allows to retrieve the TextSender
func TextMiddleware(textSender services.TextSender) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("textSender", textSender)
		c.Next()
	}
}
