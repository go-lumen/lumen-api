package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

// Mngo struct holds informations about MongoDB database
type Mngo struct {
	database *mongo.Database
	dbName   string
	context  context.Context
}

// New creates a database connexion
func New(context context.Context, database *mongo.Database, dbName string) *Mngo {
	return &Mngo{database, dbName, context}
}
