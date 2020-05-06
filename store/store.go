package store

import (
	"github.com/go-lumen/lumen-api/helpers/params"
	"github.com/go-lumen/lumen-api/models"
)

// Store interface
type Store interface {
	CreateUser(*models.User) error
	GetUserByID(string) (*models.User, error)
	GetUser(params.M) (*models.User, error)
	UpdateUser(string, *models.User) error
	DeleteUser(string) error
	ActivateUser(string, string) error
	GetUsers(string) ([]*models.User, error)
	CountUsers() (int, error)
	UserExists(string) (bool, *models.User, error)
}
