package server

import (
	"github.com/adrien3d/lumen-api/models"
	"github.com/adrien3d/lumen-api/utils"

	"github.com/globalsign/mgo"
)

// SetupMongoIndexes allows to setup MongoDB index
func (a *API) SetupMongoIndexes() error {
	database := a.MongoDatabase

	// Creates a list of indexes to ensure
	collectionIndexes := make(map[*mgo.Collection][]mgo.Index)

	// User indexes
	users := database.C(models.UsersCollection)
	collectionIndexes[users] = []mgo.Index{
		{
			Key:    []string{"email"},
			Unique: true,
		},
	}

	for collection, indexes := range collectionIndexes {
		for _, index := range indexes {
			err := collection.EnsureIndex(index)

			utils.CheckErr(err)
		}
	}
	return nil
}

/*func CreateValidator(collection *mgo.Collection, validator bson.M) {
	info := &mgo.CollectionInfo{
		Validator:       validator,
		ValidationLevel: "strict",
	}
	collection.Create(info)
}*/
