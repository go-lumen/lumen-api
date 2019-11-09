package mongodb

import (
	"net/http"

	"github.com/globalsign/mgo/bson"
	"go-lumen/lumen-api/helpers"
	"go-lumen/lumen-api/helpers/params"
	"go-lumen/lumen-api/models"
)

// CreateUser checks if user already exists, and if not, creates it
func (db *mongo) CreateUser(user *models.User) error {
	session := db.Session.Copy()
	defer session.Close()
	users := db.C(models.UsersCollection).With(session)

	err := user.BeforeCreate()
	if err != nil {
		return err
	}

	if count, _ := users.Find(bson.M{"email": user.Email}).Count(); count > 0 {
		return helpers.NewError(http.StatusConflict, "user_already_exists", "User already exists", err)
	}

	err = users.Insert(user)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "user_creation_failed", "Failed to insert the user in the database", err)
	}

	return nil
}

// FindUserById allows to retrieve a user by its id
func (db *mongo) FindUserById(id string) (*models.User, error) {
	session := db.Session.Copy()
	defer session.Close()
	users := db.C(models.UsersCollection).With(session)

	user := &models.User{}
	err := users.FindId(id).One(user)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "user_not_found", "User not found", err)
	}

	return user, err
}

// FindUser allows to retrieve a user by its characteristics
func (db *mongo) FindUser(params params.M) (*models.User, error) {
	session := db.Session.Copy()
	defer session.Close()
	users := db.C(models.UsersCollection).With(session)

	user := &models.User{}

	err := users.Find(params).One(user)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "user_not_found", "User not found", err)
	}

	return user, err
}

// DeleteUser allows to delete a user by its id
func (db *mongo) DeleteUser(user *models.User, userId string) error {
	session := db.Session.Copy()
	defer session.Close()
	users := db.C(models.UsersCollection).With(session)

	err := users.Remove(bson.M{"_id": userId})
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "user_delete_failed", "Failed to delete the user", err)
	}

	return nil
}

// ActivateUser allows to activate a user by its id
func (db *mongo) ActivateUser(activationKey string, id string) error {
	session := db.Session.Copy()
	defer session.Close()
	users := db.C(models.UsersCollection).With(session)

	err := users.Update(bson.M{"$and": []bson.M{{"_id": id}, {"activation_key": activationKey}}}, bson.M{"$set": bson.M{"active": true}})
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "user_activation_failed", "Couldn't find the user to activate", err)
	}
	return nil
}

// ChangeLanguage allows to change a user language by its id
func (db *mongo) ChangeLanguage(id string, language string) error {
	session := db.Session.Copy()
	defer session.Close()
	users := db.C(models.UsersCollection).With(session)

	if err := users.UpdateId(id, bson.M{"$set": bson.M{"language": language}}); err != nil {
		return helpers.NewError(http.StatusInternalServerError, "user_activation_failed", "Couldn't find the user to change language", err)
	}
	return nil
}

// UpdateUser allows to update one or more user characteristics
func (db *mongo) UpdateUser(userId string, params params.M) error {
	session := db.Session.Copy()
	defer session.Close()
	users := db.C(models.UsersCollection).With(session)

	if err := users.UpdateId(userId, params); err != nil {
		return helpers.NewError(http.StatusInternalServerError, "user_update_failed", "Failed to update the user", err)
	}

	return nil
}

// GetUsers allows to get all users
func (db *mongo) GetUsers() ([]*models.User, error) {
	session := db.Session.Copy()
	defer session.Close()

	users := db.C(models.UsersCollection).With(session)

	list := []*models.User{}
	if err := users.Find(params.M{}).All(&list); err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "users_not_found", "Users not found", err)
	}

	return list, nil
}

// CountUsers allows to count all users
func (db *mongo) CountUsers() (int, error) {
	session := db.Session.Copy()
	defer session.Close()

	users := db.C(models.UsersCollection).With(session)

	nbr, err := users.Find(params.M{}).Count()
	if err != nil {
		return -1, helpers.NewError(http.StatusNotFound, "users_not_found", "Users not found", err)
	}
	return nbr, nil
}

// UserExists allows to know if a user exists through his mail
func (db *mongo) UserExists(userEmail string) (bool, error) {
	session := db.Session.Copy()
	defer session.Close()
	users := db.C(models.UsersCollection).With(session)

	user := &models.User{}

	err := users.Find(params.M{"email": userEmail}).One(user)
	if err != nil {
		return false, helpers.NewError(http.StatusNotFound, "user_not_found", "User not found", err)
	}

	return true, err
}
