package middlewares

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/go-lumen/lumen-api/store"
	"github.com/go-lumen/lumen-api/store/mongodb"
	"github.com/go-lumen/lumen-api/store/postgresql"
	"github.com/jinzhu/gorm"
	"go.mongodb.org/mongo-driver/mongo"
)

// StoreMongoMiddleware allows to setup MongoDB database
func StoreMongoMiddleware(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		//fmt.Println("Store mongo middleware:", db, config.GetString(c, "mongo_db_name"), c, models.OrganizationsCollection)
		//store.ToContext(c, mongodb.New(db, config.GetString(c, "mongo_db_name"), c))
		store.ToContext(c, mongodb.New(c, db, "lumen"))
		c.Next()
	}
}

// StorePostgreMiddleware allows to setup SQL database
func StorePostgreMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		store.ToContext(c, postgresql.New(db))
		c.Next()
	}
}

// StoreMySQLMiddleware allows to setup SQL database
func StoreMySQLMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		//store.ToContext(c, mysql.New(db))
		c.Next()
	}
}
