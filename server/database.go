package server

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // For postgre DB
	"go.mongodb.org/mongo-driver/mongo"
)

type dbLogger struct{}

func (a *API) SetupMongoDatabase() (*mongo.Database, error) {
	uri := a.Config.GetString("mongo_db_prefix")

	if (a.Config.GetString("mongo_db_user") != "") && (a.Config.GetString("mongo_db_password") != "") {
		uri += a.Config.GetString("mongo_db_user") + ":" + a.Config.GetString("mongo_db_password") + "@"
	}

	uri += a.Config.GetString("mongo_db_host")

	//utils.Log(nil, "info", uri)

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
