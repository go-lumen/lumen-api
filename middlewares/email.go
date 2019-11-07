package middlewares

import (
	"github.com/gin-gonic/gin"
	"go-lumen/lumen-api/services"
)

// EmailMiddleware allows to retrieve EmailSender
func EmailMiddleware(emailSender services.EmailSender) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("emailSender", emailSender)
		c.Next()
	}
}
