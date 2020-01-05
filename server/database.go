package server

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-lumen/lumen-api/utils"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type dbLogger struct{}

// SetupMongoDatabase establishes the connexion with the mongo database
func (a *API) SetupMongoDatabase() (*mongo.Database, error) {
	uri := "mongodb://"

	if (a.Config.GetString("mongo_db_user") != "") && (a.Config.GetString("mongo_db_password") != "") {
		uri += a.Config.GetString("mongo_db_user") + ":" + a.Config.GetString("mongo_db_password") + "@"
	}
	uri += a.Config.GetString("mongo_db_host")
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("Mongo client couldn't connect with background context: %v", err)
	}
	database := client.Database(a.Config.GetString("mongo_db_name"))
	a.MongoDatabase = database

	return database, nil
}

// SetupPostgreDatabase establishes the connexion with the PostgreSQL database
func (a *API) SetupPostgreDatabase() (*gorm.DB, error) {
	connectionURI := fmt.Sprintf(
		"sslmode=disable dbname=%s host=%s port=%s user=%s password=%s",
		a.Config.GetString("postgres_db_name"),
		a.Config.GetString("postgres_db_addr"),
		a.Config.GetString("postgres_db_port"),
		a.Config.GetString("postgres_db_user"),
		a.Config.GetString("postgres_db_password"),
	)

	db, err := gorm.Open("postgres", connectionURI)
	if err != nil {
		return nil, err
	}

	// Debug database logs
	debugDatabase := a.Config.GetBool("debug_database")
	db.LogMode(debugDatabase)

	db.DB().SetConnMaxLifetime(time.Minute * 5)
	db.DB().SetMaxIdleConns(5)
	db.DB().SetMaxOpenConns(5)

	a.PostgreDatabase = db
	return db, nil
}

// SetupMySQLDatabase establishes the connexion with the MySQL database
func (a *API) SetupMySQLDatabase() (*sql.DB, error) {
	db, err := sql.Open("mysql", a.Config.GetString("mysql_db_user")+":"+a.Config.GetString("mysql_db_password")+
		"@/"+a.Config.GetString("mysql_db_name"))
	utils.CheckErr(err)

	a.MySQLDatabase = db

	return db, nil
}
