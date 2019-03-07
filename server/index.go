package server

import (
	"github.com/adrien3d/base-api/models"

	"github.com/globalsign/mgo"
)

// SetupIndexes allows to setup MongoDB index
func (a *API) SetupIndexes() error {
	database := a.Database

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

			if err != nil {
				return err
			}
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
