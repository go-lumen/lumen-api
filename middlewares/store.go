package middlewares

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/go-lumen/lumen-api/store"
	"github.com/go-lumen/lumen-api/store/mongodb"
	"github.com/go-lumen/lumen-api/store/postgresql"
)

// StoreMongoMiddleware allows to setup MongoDB database
func StoreMongoMiddleware(db *mgo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		store.ToContext(c, mongodb.New(db))
		c.Next()
	}
}

// StoreSQLMiddleware allows to setup SQL database
func StoreSQLMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		store.ToContext(c, postgresql.New(db))
		c.Next()
	}
}
