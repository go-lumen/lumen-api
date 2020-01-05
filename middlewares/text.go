package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/go-lumen/lumen-api/services"
)

// TextMiddleware allows to retrieve the TextSender
func TextMiddleware(textSender services.TextSender) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("textSender", textSender)
		c.Next()
	}
}
