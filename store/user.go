package store

import (
	"context"

	"github.com/adrien3d/lumen-api/helpers/params"
	"github.com/adrien3d/lumen-api/models"
)

// CreateUser checks if user already exists, and if not, creates it
func CreateUser(c context.Context, record *models.User) error {
	return FromContext(c).CreateUser(record)
}

// FindUserById allows to retrieve a user by its id
func FindUserById(c context.Context, id string) (*models.User, error) {
	return FromContext(c).FindUserById(id)
}

// FindUser allows to retrieve a user by its characteristics
func FindUser(c context.Context, params params.M) (*models.User, error) {
	return FromContext(c).FindUser(params)
}

// DeleteUser allows to delete a user by its id
func DeleteUser(c context.Context, userId string) error {
	return FromContext(c).DeleteUser(Current(c), userId)
}

// ActivateUser allows to activate a user by its id
func ActivateUser(c context.Context, activationKey string, id string) error {
	return FromContext(c).ActivateUser(activationKey, id)
}

// ChangeLanguage allows to change a user language by its id
func ChangeLanguage(c context.Context, id string, language string) error {
	return FromContext(c).ChangeLanguage(id, language)
}

// UpdateUser allows to update one or more user characteristics
func UpdateUser(c context.Context, userId string, params params.M) error {
	return FromContext(c).UpdateUser(userId, params)
}

// GetUsers allows to get all users
func GetUsers(c context.Context) ([]*models.User, error) {
	return FromContext(c).GetUsers()
}

// CountUsers allows to count all users
func CountUsers(c context.Context) (int, error) {
	return FromContext(c).CountUsers()
}

func UserExists(c context.Context, userEmail string) (bool, error) {
	return FromContext(c).UserExists(userEmail)
}
