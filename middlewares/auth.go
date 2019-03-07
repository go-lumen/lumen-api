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

		user.LastAccess = time.Now().Unix()

		store.UpdateUser(c, params.M{"$set": params.M{"last_access": user.LastAccess}})

		c.Next()
	}
}
