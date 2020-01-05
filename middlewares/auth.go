package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-lumen/lumen-api/config"
	"github.com/go-lumen/lumen-api/helpers"
	"github.com/go-lumen/lumen-api/helpers/params"
	"github.com/go-lumen/lumen-api/store"
	"net/http"
	"strings"
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

		user, _ := store.GetUser(c, params.M{"_id": claims["sub"].(string)})
		//logrus.Infoln("looking for: " + claims["sub"].(string) + " Got user: " + fmt.Sprint(user))
		c.Set(store.CurrentKey, user)
		/*if err := store.UpdateUser(c, string(user.Id), params.M{"$set": params.M{"last_access": time.Now().Unix()}}); err != nil {
			println(err)
		}*/

		/*if user.LastPasswordUpdate > int64(claims["iat"].(float64)) {
			c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_token_new_password", "the given token is invalid due to new password", err))
		}*/

		c.Next()
	}
}
