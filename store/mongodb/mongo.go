package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type mngo struct {
	database *mongo.Database
	dbName   string
	context  context.Context
}

// New creates a database connexion
func New(database *mongo.Database, dbName string, context context.Context) *mngo {
	return &mngo{database, dbName, context}
}
