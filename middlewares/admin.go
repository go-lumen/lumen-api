package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-lumen/lumen-api/helpers"
	"github.com/go-lumen/lumen-api/store"
	"net/http"
)

// AdminMiddleware allows to filter admin only users
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := store.Current(c)

		if !(user.Role == "admin") {
			c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("admin_required", "The user is not administrator", errors.New("The user is not administrator")))
			return
		}

		c.Next()
	}
}
