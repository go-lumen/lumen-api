package postgresql

import (
	"fmt"
	"github.com/go-lumen/lumen-api/helpers"
	"github.com/go-lumen/lumen-api/helpers/params"
	"github.com/go-lumen/lumen-api/models"
	"net/http"
)

// CreateUser checks if user already exists, and if not, creates it
func (db *Postgres) CreateUser(user *models.User) error {
	var count int
	if err := db.Model(user).Where("email = ?", user.Email).Count(&count).Error; err != nil || count > 0 {
		fmt.Println("user_exists", err)
		return helpers.NewError(http.StatusBadRequest, "user_exists", "the user already exists", err)
	}

	if err := db.Create(user).Error; err != nil {
		fmt.Println("user_creation_failed", err)
		return helpers.NewError(http.StatusInternalServerError, "user_creation_failed", "could not create the user", err)
	}

	return nil
}

// GetUserByID allows to retrieve a user by its id
func (db *Postgres) GetUserByID(id string) (*models.User, error) {
	var user models.User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "user_not_found", "could not find the user", err)
	}
	return &user, nil
}

// GetUser allows to retrieve a user by its characteristics
func (db *Postgres) GetUser(params params.M) (*models.User, error) {
	session := db.New()

	var user models.User
	for key, value := range params {
		session = session.Where(key+" = ?", value)
	}

	if err := session.First(&user).Error; err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "user_not_found", "could not find the user", err)
	}

	return &user, nil
}

// DeleteUser allows to delete a user by its id
func (db *Postgres) DeleteUser(userID string) error {
	return nil
}

// ActivateUser allows to activate a user by its id
func (db *Postgres) ActivateUser(activationKey string, email string) error {
	var user models.User
	fmt.Println("Trying to find user:", email, "with activationKey:", activationKey)
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return helpers.NewError(http.StatusNotFound, "user_not_found", "could not find the user", err)
	}

	if user.Key != activationKey {
		return helpers.NewError(http.StatusBadRequest, "invalid_validation_code", "the provided activation code is invalid", nil)
	}

	if err := db.Model(&user).Update("active", true).Error; err != nil {
		return helpers.NewError(http.StatusInternalServerError, "update_user_failed", "could not update the user", err)
	}
	fmt.Println("Final user:", user)

	return nil
}

// ChangeLanguage allows to change a user language by its id
func (db *Postgres) ChangeLanguage(id string, language string) error {
	return nil
}

// UpdateUser allows to update one or more user characteristics
func (db *Postgres) UpdateUser(userID string, newUser *models.User) error {
	return nil
}

// UpdateUserFields allows to update one or more user characteristics
func (db *Postgres) UpdateUserFields(userID string, fields params.M) error {
	return nil
}

// GetUsers allows to get all users
func (db *Postgres) GetUsers(groupID string) ([]*models.User, error) {
	var users []*models.User
	db.Find(&users)

	return users, nil
}

// CountUsers allows to count all users
func (db *Postgres) CountUsers() (int, error) {
	return 0, nil
}

// UserExists allows to know if a user exists through his mail
func (db *Postgres) UserExists(userEmail string) (bool, *models.User, error) {
	var user models.User
	if err := db.Where("email = ?", userEmail).First(&user).Error; err == nil {
		return true, &user, nil
	}
	return false, nil, nil
}
