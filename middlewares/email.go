package middlewares

import (
	"github.com/adrien3d/base-api/services"
	"github.com/gin-gonic/gin"
)

// EmailMiddleware allows to retrieve EmailSender
func EmailMiddleware(emailSender services.EmailSender) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("emailSender", emailSender)
		c.Next()
	}
}
