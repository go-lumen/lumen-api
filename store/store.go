package store

import (
	"go-lumen/lumen-api/helpers/params"
	"go-lumen/lumen-api/models"
)

// Store interface
type Store interface {
	CreateUser(*models.User) error
	DeleteUser(*models.User, string) error
	FindUserById(string) (*models.User, error)
	ActivateUser(string, string) error
	FindUser(params.M) (*models.User, error)
	UpdateUser(string, params.M) error
	ChangeLanguage(string, string) error
	GetUsers() ([]*models.User, error)
	CountUsers() (int, error)
	UserExists(string) (bool, error)
}
