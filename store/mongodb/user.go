package mongodb

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"

	"github.com/globalsign/mgo/bson"
	"github.com/go-lumen/lumen-api/helpers"
	"github.com/go-lumen/lumen-api/helpers/params"
	"github.com/go-lumen/lumen-api/models"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateUser checks if user already exists, and if not, creates it
func (db *mngo) CreateUser(user *models.User) error {
	c := db.database.Collection(models.UsersCollection)

	err := user.BeforeCreate()
	user.Id = bson.NewObjectId().Hex()
	if err != nil {
		return err
	}
	cursor, err := c.Find(db.context, bson.M{"email": user.Email})
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}
	if len(results) > 0 {
		return helpers.NewError(http.StatusConflict, "user_already_exists", "User already exists", err)
	}

	_, err = c.InsertOne(db.context, user)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "user_creation_failed", "Failed to insert the user in the database", err)
	}
	//res.InsertedID

	return nil
}

// GetUserById allows to retrieve a user by its id
func (db *mngo) GetUserById(id string) (*models.User, error) {
	c := db.database.Collection(models.UsersCollection)

	user := &models.User{}
	err := c.FindOne(db.context, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "user_not_found", "User not found", err)
	}

	return user, err
}

// GetUser allows to retrieve a user by its characteristics
func (db *mngo) GetUser(params params.M) (*models.User, error) {
	c := db.database.Collection(models.UsersCollection)

	user := &models.User{}
	err := c.FindOne(db.context, params).Decode(&user)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "user_not_found", "User not found", err)
	}

	return user, err
}

// GetUserFromSigfoxId allows to retrieve a user by its Sigfox Id
func (db *mngo) GetUserFromSigfoxId(sigfoxId string) (*models.User, error) {
	c := db.database.Collection(models.UsersCollection)

	user := &models.User{}
	err := c.FindOne(db.context, bson.M{"sigfox_id": sigfoxId}).Decode(&user)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "user_not_found", "User not found", err)
	}

	return user, err
}

// GetUserById allows to retrieve a user by its id
func (db *mngo) UserExists(email string) (bool, *models.User, error) {
	c := db.database.Collection(models.UsersCollection)

	user := &models.User{}
	err := c.FindOne(db.context, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return false, nil, helpers.NewError(http.StatusNotFound, "user_not_found", "User not found", err)
	}

	return true, user, err
}

// UpdateUser allows to update one or more user characteristics
func (db *mngo) UpdateUser(userId string, newUser *models.User) error {
	c := db.database.Collection(models.UsersCollection)

	result, err := c.UpdateOne(context.TODO(), bson.M{"_id": userId}, bson.M{"$set": newUser}, options.Update().SetUpsert(true))
	if err != nil {
		log.Fatal(err)
		return helpers.NewError(http.StatusInternalServerError, "user_update_failed", "Failed to update the user", err)
	}

	if result.MatchedCount != 0 {
		fmt.Println("matched and replaced an existing document")
		return nil
	}
	if result.UpsertedCount != 0 {
		fmt.Printf("UpdateUser: inserted a new document with ID %v\n", result.UpsertedID)
		return nil
	}

	return nil
}

// DeleteUser allows to delete a user by its id
func (db *mngo) DeleteUser(userId string) error {
	c := db.database.Collection(models.UsersCollection)

	_, err := c.DeleteOne(db.context, bson.M{"_id": userId})
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "user_delete_failed", "Failed to delete the user", err)
	}

	//res.DeletedCount

	return nil
}

func (db *mngo) ActivateUser(activationKey string, id string) error {
	c := db.database.Collection(models.UsersCollection)

	result, err := c.UpdateOne(context.TODO(), bson.M{"$and": []bson.M{{"_id": id}, {"activation_key": activationKey}}}, bson.M{"$set": bson.M{"status": "activated"}}, options.Update().SetUpsert(false))

	if err != nil {
		log.Fatal(err)
		return helpers.NewError(http.StatusInternalServerError, "user_activation_failed", "Couldn't find the user to activate", err)
	}

	if result.MatchedCount != 0 {
		fmt.Println("matched and activated a user")
		return nil
	}
	if result.UpsertedCount != 0 {
		fmt.Printf("ActivateUser: inserted a new document with ID %v\n", result.UpsertedID)
		return nil
	}

	return nil
}

// GetUsers allows to get all users
func (db *mngo) GetUsers(groupId string) ([]*models.User, error) {
	c := db.database.Collection(models.UsersCollection)

	list := []*models.User{}

	var filter bson.M
	if groupId != "" {
		filter = bson.M{"group_id": groupId}
	}
	cur, err := c.Find(context.TODO(), filter)
	if err != nil {
		logrus.Warnln("Error on Finding all the documents", err)
	}

	for cur.Next(context.TODO()) {
		var elem models.User
		err = cur.Decode(&elem)
		if err != nil {
			logrus.Warnln("Error on Decoding the document", err)
		}
		list = append(list, &elem)
	}

	return list, err
}

// CountUsers allows to count all users
func (db *mngo) CountUsers() (int, error) {
	c := db.database.Collection(models.UsersCollection)

	count, err := c.CountDocuments(context.TODO(), nil)

	return int(count), err
}
