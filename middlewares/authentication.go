package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-lumen/lumen-api/config"
	"github.com/go-lumen/lumen-api/helpers"
	"github.com/go-lumen/lumen-api/models"
	"github.com/go-lumen/lumen-api/store"
	"github.com/go-lumen/lumen-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"strings"
)

// AuthenticationMiddleware allows to analyze the token and check that it is valid
func AuthenticationMiddleware() gin.HandlerFunc {
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

		ctx := store.AuthContext(c)
		user, _ := models.GetUser(ctx, bson.M{"id": claims["sub"].(string)})
		//logrus.Infoln("looking for: " + claims["sub"].(string) + " Got user: " + fmt.Sprint(user))
		c.Set(store.CurrentUserKey, user)
		group, err := models.GetGroup(ctx, bson.M{"id": user.GroupID})
		if err != nil {
			utils.Log(c, "error", "Group not found")
		}
		c.Set(store.CurrentUserGroupKey, group)
		/*if err := store.UpdateUser(c, string(user.ID), bson.M{"$set": bson.M{"last_access": time.Now().Unix()}}); err != nil {
			println(err)
		}*/

		/*if user.LastPasswordUpdate > int64(claims["iat"].(float64)) {
			c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_token_new_password", "the given token is invalid due to new password", err))
		}*/

		c.Next()
	}
}
