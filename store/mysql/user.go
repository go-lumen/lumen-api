package mysql

import (
	"fmt"
	"github.com/go-lumen/lumen-api/helpers/params"
	"github.com/go-lumen/lumen-api/models"
)

// CreateUser checks if user already exists, and if not, creates it
func (db *mysql) CreateUser(user *models.User) error {
	return nil
}

// FindUserById allows to retrieve a user by its id
func (db *mysql) FindUserById(id string) (*models.User, error) {
	return nil, nil
}

// FindUser allows to retrieve a user by its characteristics
func (db *mysql) FindUser(params params.M) (*models.User, error) {
	fmt.Println("finding user:", params)
	//rows, err := db.Query("SELECT * FROM users WHERE email = $1", params)

	return nil, nil
}

// DeleteUser allows to delete a user by its id
func (db *mysql) DeleteUser(user *models.User, userId string) error {
	return nil
}

// ActivateUser allows to activate a user by its id
func (db *mysql) ActivateUser(activationKey string, id string) error {
	return nil
}

// ChangeLanguage allows to change a user language by its id
func (db *mysql) ChangeLanguage(id string, language string) error {
	return nil
}

// UpdateUser allows to update one or more user characteristics
func (db *mysql) UpdateUser(userId string, params params.M) error {
	return nil
}

// GetUsers allows to get all users
func (db *mysql) GetUsers() ([]*models.User, error) {
	return nil, nil
}

// CountUsers allows to count all users
func (db *mysql) CountUsers() (int, error) {
	return 0, nil
}

// UserExists allows to know if a user exists through his mail
func (db *mysql) UserExists(userEmail string) (bool, error) {
	return false, nil
}
