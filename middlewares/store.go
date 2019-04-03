package middlewares

import (
	"github.com/adrien3d/lumen-api/store"
	"github.com/adrien3d/lumen-api/store/mongodb"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
)

// StoreMiddleware allows to setup database
func StoreMiddleware(db *mgo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		store.ToContext(c, mongodb.New(db))
		c.Next()
	}
}
