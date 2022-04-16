package middlewares

import (
	"github.com/adrien3d/stokelp-poc/config"
	"github.com/adrien3d/stokelp-poc/helpers"
	"github.com/adrien3d/stokelp-poc/models"
	"github.com/adrien3d/stokelp-poc/store"
	"github.com/adrien3d/stokelp-poc/utils"
	"github.com/gin-gonic/gin"
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
