package store

import (
	"github.com/adrien3d/lumen-api/helpers/params"
	"github.com/adrien3d/lumen-api/models"
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
}
