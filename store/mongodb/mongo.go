package mongodb

import "github.com/globalsign/mgo"

type mongo struct {
	*mgo.Database
}

// New creates a database connexion
func New(database *mgo.Database) *mongo {
	return &mongo{database}
}
