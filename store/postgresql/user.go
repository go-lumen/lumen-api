package postgresql

import (
	"fmt"
	"github.com/go-lumen/lumen-api/helpers/params"
	"github.com/go-lumen/lumen-api/models"
	"github.com/go-lumen/lumen-api/utils"
)

// CreateUser checks if user already exists, and if not, creates it
func (db *postgres) CreateUser(user *models.User) error {
	err := db.Insert(user)
	utils.CheckErr(err)
	return err
}

// FindUserById allows to retrieve a user by its id
func (db *postgres) FindUserById(id string) (*models.User, error) {
	fmt.Println("finding user:", id)
	user := &models.User{}
	err := db.Model(models.User{}).Where("id=?", id).Select(user) //.Order("id")
	utils.CheckErr(err)

	return user, err
}

// FindUser allows to retrieve a user by its characteristics
func (db *postgres) FindUser(params params.M) (*models.User, error) {
	//rows, err := db.Query("SELECT * FROM users WHERE email = $1", params)

	return nil, nil
}

// DeleteUser allows to delete a user by its id
func (db *postgres) DeleteUser(user *models.User, userId string) error {
	return nil
}

// ActivateUser allows to activate a user by its id
func (db *postgres) ActivateUser(activationKey string, id string) error {
	return nil
}

// ChangeLanguage allows to change a user language by its id
func (db *postgres) ChangeLanguage(id string, language string) error {
	return nil
}

// UpdateUser allows to update one or more user characteristics
func (db *postgres) UpdateUser(userId string, params params.M) error {
	return nil
}

// GetUsers allows to get all users
func (db *postgres) GetUsers() ([]*models.User, error) {
	users := make([]*models.User, 0)
	err := db.Model(models.User{}).Order("id").Select(users)
	utils.CheckErr(err)

	return users, err
}

// CountUsers allows to count all users
func (db *postgres) CountUsers() (int, error) {
	return 0, nil
}

// UserExists allows to know if a user exists through his mail
func (db *postgres) UserExists(userEmail string) (bool, error) {
	return false, nil
}
