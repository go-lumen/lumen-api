package store

import (
	"context"

	"github.com/adrien3d/base-api/helpers/params"
	"github.com/adrien3d/base-api/models"
)

func CreateUser(c context.Context, record *models.User) error {
	return FromContext(c).CreateUser(record)
}

func DeleteUser(c context.Context, userId string) error {
	return FromContext(c).DeleteUser(Current(c), userId)
}

func FindUserById(c context.Context, id string) (*models.User, error) {
	return FromContext(c).FindUserById(id)
}

func FindUser(c context.Context, params params.M) (*models.User, error) {
	return FromContext(c).FindUser(params)
}

func GetUsers(c context.Context) ([]*models.User, error) {
	return FromContext(c).GetUsers()
}

func ActivateUser(c context.Context, activationKey string, id string) error {
	return FromContext(c).ActivateUser(activationKey, id)
}

func ChangeLanguage(c context.Context, id string, language string) error {
	return FromContext(c).ChangeLanguage(id, language)
}

func UpdateUser(c context.Context, params params.M) error {
	return FromContext(c).UpdateUser(Current(c), params)
}

func CountUsers(c context.Context) (int, error) {
	return FromContext(c).CountUsers()
}
