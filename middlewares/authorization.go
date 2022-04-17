package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/go-lumen/lumen-api/config"
	"github.com/go-lumen/lumen-api/helpers"
	"github.com/go-lumen/lumen-api/models"
	"github.com/go-lumen/lumen-api/store"
	"github.com/go-lumen/lumen-api/utils"
	"strings"
)

// AuthorizationMiddleware allows granting access to a resource or not
func AuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenReader := c.Request.Header.Get("Authorization")

		authHeaderParts := strings.Split(tokenReader, " ")
		encodedKey := []byte(config.GetString(c, "rsa_private"))
		claims, _ := helpers.ValidateJwtToken(authHeaderParts[1], encodedKey, "access")
		ctx := store.AuthContext(c)
		user, _ := models.GetUser(ctx, utils.ParamID(claims["sub"].(string)))
		c.Set(store.CurrentUserKey, user)

		c.Next()
	}
}
