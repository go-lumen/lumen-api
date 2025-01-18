package server

import (
	"context"
	"time"

	"github.com/go-lumen/lumen-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SetupMongoIndexes allows to setup MongoDB index
func (a *API) SetupMongoIndexes() error {
	/*database := a.MongoDatabase

	collection := database.C(models.DeviceMessagesCollection)
	err := collection.EnsureIndex(mgo.Index{
		{
			Key:  []string{"$2dsphere:location"},
			Bits: 26,
		},
	})*/

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db := a.MongoDatabase

	indexOpts := options.CreateIndexes().SetMaxTime(time.Second * 10)
	// Index to location 2dsphere type.
	pointIndexModel := mongo.IndexModel{
		Options: options.Index().SetBackground(true),
		Keys:    bson.D{{Key: "location", Value: "2dsphere"}},
	}
	poiIndexes := db.Collection("pois").Indexes()
	deviceMessagesNames, err := poiIndexes.CreateOne(ctx, pointIndexModel, indexOpts)
	if err != nil {
		return err
	}
	utils.Log(nil, "info", "Index successfully created for:", deviceMessagesNames)
	/*

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

		deviceMessages := database.C(models.DeviceMessagesCollection)
		collectionIndexes[deviceMessages] = []mgo.Index{
			{
				Key:    []string{"$2dsphere:location"},
				Bits: 26,
			},
		}

		for collection, indexes := range collectionIndexes {
			for _, index := range indexes {
				err := collection.EnsureIndex(index)

				utils.CheckErr(err)
			}
		}*/
	/*var indexView *mongo.IndexView

	// Specify the MaxTime option to limit the amount of time the operation can run on the server
	opts := options.ListIndexes().SetMaxTime(2 * time.Second)
	cursor, err := indexView.List(context.TODO(), opts)
	if err != nil {
		utils.Log(nil, "error", err)
	}

	// Get a slice of all indexes returned and print them out.
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	utils.Log(nil, "info", "index results:", results)*/

	/*var indexView *mongo.IndexView

	// Create two indexes: {name: 1, email: 1} and {name: 1, age: 1}
	// For the first index, specify no options. The name will be generated as "name_1_email_1" by the driver.
	// For the second index, specify the Name option to explicitly set the name to "nameAge".
	models := []mongo.IndexModel{
	    {
	        Keys: bson.D{{"name", 1}, {"email", 1}},
	    },
	    {
	        Keys:    bson.D{{"name", 1}, {"age", 1}},
	        Options: options.Index().SetName("nameAge"),
	    },
	}

	// Specify the MaxTime option to limit the amount of time the operation can run on the server
	opts := options.CreateIndexes().SetMaxTime(2 * time.Second)
	names, err := indexView.CreateMany(context.TODO(), models, opts)
	if err != nil {
	    log.Fatal(err)
	}

	fmt.Printf("created indexes %v\n", names)*/
	return nil
}

/*func CreateValidator(collection *mgo.Collection, validator bson.M) {
	info := &mgo.CollectionInfo{
		Validator:       validator,
		ValidationLevel: "strict",
	}
	collection.Create(info)
}*/
