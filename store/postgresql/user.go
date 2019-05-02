package postgresql

import (
	"fmt"
	"github.com/go-lumen/lumen-api/helpers/params"
	"github.com/go-lumen/lumen-api/models"
)

// CreateUser checks if user already exists, and if not, creates it
func (db *postgre) CreateUser(user *models.User) error {
	return nil
}

// FindUserById allows to retrieve a user by its id
func (db *postgre) FindUserById(id string) (*models.User, error) {
	return nil, nil
}

// FindUser allows to retrieve a user by its characteristics
func (db *postgre) FindUser(params params.M) (*models.User, error) {
	fmt.Println("finding user:", params)
	//rows, err := db.Query("SELECT * FROM users WHERE email = $1", params)

	return nil, nil
}

// DeleteUser allows to delete a user by its id
func (db *postgre) DeleteUser(user *models.User, userId string) error {
	return nil
}

// ActivateUser allows to activate a user by its id
func (db *postgre) ActivateUser(activationKey string, id string) error {
	return nil
}

// ChangeLanguage allows to change a user language by its id
func (db *postgre) ChangeLanguage(id string, language string) error {
	return nil
}

// UpdateUser allows to update one or more user characteristics
func (db *postgre) UpdateUser(userId string, params params.M) error {
	return nil
}

// GetUsers allows to get all users
func (db *postgre) GetUsers() ([]*models.User, error) {
	return nil, nil
}

// CountUsers allows to count all users
func (db *postgre) CountUsers() (int, error) {
	return 0, nil
}

// UserExists allows to know if a user exists through his mail
func (db *postgre) UserExists(userEmail string) (bool, error) {
	return false, nil
}
