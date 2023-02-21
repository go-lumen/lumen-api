package postgresql

import (
	"fmt"
	"github.com/go-lumen/lumen-api/helpers"
	"github.com/go-lumen/lumen-api/helpers/params"
	"github.com/go-lumen/lumen-api/models"
	"github.com/go-lumen/lumen-api/utils"
	"net/http"
)

// CreateUser checks if user already exists, and if not, creates it
func (db *PSQL) CreateUser(user *models.User) error {
	var count int64
	if err := db.database.Model(user).Where("email = ?", user.Email).Count(&count).Error; err != nil || count > 0 {
		return helpers.NewError(http.StatusBadRequest, "user_exists", "the user already exists", err)
	}

	if err := db.database.Create(user).Error; err != nil {
		return helpers.NewError(http.StatusInternalServerError, "user_creation_failed", "could not create the user", err)
	}

	return nil
}

// GetUserByID allows to retrieve a user by its id
func (db *PSQL) GetUserByID(id string) (*models.User, error) {
	var user models.User
	if err := db.database.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "user_not_found", "could not find the user", err)
	}
	return &user, nil
}

// GetUser allows to retrieve a user by its characteristics
func (db *PSQL) GetUser(params params.M) (*models.User, error) {
	session := db.database

	var user models.User
	for key, value := range params {
		session = session.Where(key+" = ?", value)
	}

	if err := session.First(&user).Error; err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "user_not_found", "could not find the user", err)
	}

	return &user, nil
}

func (db *PSQL) GetOrCreateUser(user *models.User) (*models.User, error) {
	if err := db.database.Where("email = ?", user.Email).First(&user).Error; err != nil {
		utils.Log(nil, "warn", `User already exists`, err)
		dbUser, err := db.GetUser(params.M{"email": user.Email})
		if err != nil {
			utils.Log(nil, "warn", err)
		} else {
			dbUser.FirstName = user.FirstName
			dbUser.LastName = user.LastName
			dbUser.Password = user.Password
			dbUser.Email = user.Email
			dbUser.Phone = user.Phone
			if err := db.ActivateUser(dbUser.Key /*strconv.Itoa(dbUser.ID)*/, dbUser.Email); err != nil {
				utils.Log(nil, "warn", `Error when activating user`, err)
			}
			dbUser.BeforeCreate(db.database, true, true, true)
			db.UpdateUser(dbUser.ID, dbUser)
			fmt.Println("Found user", dbUser.ID, ":", dbUser)
		}
	}
	user.BeforeCreate(db.database, false, false, false)
	return user, db.CreateUser(user)
}

// DeleteUser allows to delete a user by its id
func (db *PSQL) DeleteUser(userID string) error {
	return nil
}

// ActivateUser allows to activate a user by its id
func (db *PSQL) ActivateUser(activationKey string, email string) error {
	var user models.User
	fmt.Println("Trying to find user:", email, "with activationKey:", activationKey)
	if err := db.database.Where("email = ?", email).First(&user).Error; err != nil {
		return helpers.NewError(http.StatusNotFound, "user_not_found", "could not find the user", err)
	}

	if user.Key != activationKey {
		return helpers.NewError(http.StatusBadRequest, "invalid_validation_code", "the provided activation code is invalid", nil)
	}

	if err := db.database.Model(&user).Update("status", "activated").Error; err != nil {
		return helpers.NewError(http.StatusInternalServerError, "update_user_failed", "could not update the user", err)
	}
	fmt.Println("Final user:", user)

	return nil
}

// ChangeLanguage allows to change a user language by its id
func (db *PSQL) ChangeLanguage(id string, language string) error {
	return nil
}

// UpdateUser allows to update one or more user characteristics
func (db *PSQL) UpdateUser(userID string, user *models.User) error {
	if err := db.database.Model(user).Updates(&user).Error; err != nil {
		return helpers.NewError(http.StatusInternalServerError, "update_user_failed", "could not update the user", err)
	}
	return nil
}

// GetUsers allows to get all users
func (db *PSQL) GetUsers() ([]*models.User, error) {
	var users []*models.User
	db.database.Find(&users)

	return users, nil
}

// CountUsers allows to count all users
func (db *PSQL) CountUsers() (int, error) {
	return 0, nil
}

// UserExists allows to know if a user exists through his mail
func (db *PSQL) UserExists(userEmail string) (bool, error) {
	var user models.User
	if err := db.database.Where("email = ?", userEmail).First(&user).Error; err == nil {
		return true, nil
	}
	return false, nil
}
