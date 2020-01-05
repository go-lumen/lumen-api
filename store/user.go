package store

import (
	"context"

	"github.com/go-lumen/lumen-api/helpers/params"
	"github.com/go-lumen/lumen-api/models"
)

// CreateUser checks if user already exists, and if not, creates it
func CreateUser(c context.Context, record *models.User) error {
	return FromContext(c).CreateUser(record)
}

// GetUserById allows to retrieve a user by its id
func GetUserById(c context.Context, id string) (*models.User, error) {
	return FromContext(c).GetUserById(id)
}

// GetUser allows to retrieve a user by its characteristics
func GetUser(c context.Context, params params.M) (*models.User, error) {
	return FromContext(c).GetUser(params)
}

// GetUsers allows to get all users
func GetUsers(c context.Context, groupId string) ([]*models.User, error) {
	return FromContext(c).GetUsers(groupId)
}

// ActivateUser allows to activate a user by its id
func ActivateUser(c context.Context, activationKey string, id string) error {
	return FromContext(c).ActivateUser(activationKey, id)
}

// UpdateUser allows to update one or more user characteristics
func UpdateUser(c context.Context, userId string, newUser *models.User) error {
	return FromContext(c).UpdateUser(userId, newUser)
}

// DeleteUser allows to delete a user by its id
func DeleteUser(c context.Context, userId string) error {
	return FromContext(c).DeleteUser(userId) //Current(c)
}

// CountUsers allows to count all users
func CountUsers(c context.Context) (int, error) {
	return FromContext(c).CountUsers()
}

// UserExists checks if user exists
func UserExists(c context.Context, userEmail string) (bool, *models.User, error) {
	return FromContext(c).UserExists(userEmail)
}
