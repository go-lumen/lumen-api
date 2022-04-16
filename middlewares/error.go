package middlewares

import (
	"github.com/adrien3d/stokelp-poc/helpers"
	"github.com/gin-gonic/gin"
)

// ErrorMiddleware with logging if there is an error
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		errorToPrint := c.Errors.Last()
		if errorToPrint != nil {
			original, ok := errorToPrint.Err.(helpers.Error)
			if ok {
				if original.HTTPCode == 0 {
					original.HTTPCode = -1
				}
				c.JSON(original.HTTPCode, gin.H{"errors": gin.H{"message": original.Message, "code": original.Code}})
			} else {
				c.JSON(-1, gin.H{"errors": gin.H{"message": errorToPrint.Error(), "code": "unknown"}})
			}
		}
	}
}
