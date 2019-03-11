package middlewares

import (
	"errors"
	"github.com/adrien3d/base-api/config"
	"github.com/adrien3d/base-api/helpers"
	"github.com/adrien3d/base-api/helpers/params"
	"github.com/adrien3d/base-api/store"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

// AuthMiddleware allows to analyze the token and check that it is valid
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenReader := c.Request.Header.Get("Authorization")

		authHeaderParts := strings.Split(tokenReader, " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			c.AbortWithError(http.StatusBadRequest, errors.New("Authorization header format must be Bearer {token}"))
			return
		}

		encodedKey := []byte(config.GetString(c, "rsa_private"))
		claims, err := helpers.ValidateJwtToken(authHeaderParts[1], encodedKey, "access")
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_token", "the given token is invalid", err))
			return
		}

		user, _ := store.FindUserById(c, claims["sub"].(string))
		c.Set(store.CurrentKey, user)

		if err := store.UpdateUser(c, user.Id, params.M{"$set": params.M{"last_access": time.Now().Unix()}}); err != nil {
			println(err)
		}

		c.Next()
	}
}
