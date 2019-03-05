package middlewares

import (
	"github.com/adrien3d/base-api/store"
	"github.com/adrien3d/base-api/store/mongodb"
	"github.com/gin-gonic/gin"
	mgo "github.com/globalsign/mgo"
)

func StoreMiddleware(db *mgo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		store.ToContext(c, mongodb.New(db))
		c.Next()
	}
}
