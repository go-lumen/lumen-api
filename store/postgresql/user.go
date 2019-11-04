package postgresql

import (
	"fmt"
	"go-lumen/lumen-api/helpers"
	"go-lumen/lumen-api/helpers/params"
	"go-lumen/lumen-api/models"
	"net/http"
)

// CreateUser checks if user already exists, and if not, creates it
func (db *postgres) CreateUser(user *models.User) error {
	var count int
	if err := db.Model(user).Where("email = ?", user.Email).Count(&count).Error; err != nil || count > 0 {
		//return helpers.NewError(http.StatusBadRequest, "user_exists", "the user already exists")
	}
	fmt.Println("User exists Passed")

	if err := db.Create(user).Error; err != nil {
		fmt.Println("user_creation_failed", err)
		return helpers.NewError(http.StatusInternalServerError, "user_creation_failed", "could not create the user", err)
	}
	fmt.Println("Users creation Passed")

	return nil
}

// FindUserById allows to retrieve a user by its id
func (db *postgres) FindUserById(id string) (*models.User, error) {
	var user models.User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "user_not_found", "could not find the user", err)
	}
	return &user, nil
}

// FindUser allows to retrieve a user by its characteristics
/*func (db *postgres) FindUser(params params.M) (res *models.User, err error) {
	//rows, err := db.Query("SELECT * FROM users WHERE email = $1", params)

	reqStr := ""
	availableParams := []string{"first_name", "last_name", "email", "phone", "activation_key", "reset_key"}
	reqParams := []string{}
	i := 0
	for key, value := range params {
		//fmt.Println(key, ":", value)
		for _, availParam := range availableParams {
			if key == availParam {
				reqParams = append(reqParams, fmt.Sprintf("%v", value))
				if i != 0 {
					reqStr += "&"
				}
				reqStr += key + "=?"
				i++
			}
		}
	}
	fmt.Println(reqStr, reqParams)
	db.Where(reqStr, reqParams).Find(&res)
	fmt.Println("err", err)
	fmt.Println(res)

	return res, err
}*/
func (db *postgres) FindUser(params params.M) (*models.User, error) {
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
func (db *postgres) DeleteUser(user *models.User, userId string) error {
	return nil
}

// ActivateUser allows to activate a user by its id
func (db *postgres) ActivateUser(activationKey string, id string) error {
	var user models.User
	if err := db.Where("id = ?", id).First(&user); err != nil {
		return helpers.NewError(http.StatusNotFound, "user_not_found", "could not find the user", err.Error)
	}

	if user.ActivationKey != activationKey {
		return helpers.NewError(http.StatusBadRequest, "invalid_validation_code", "the provided activation code is invalid", nil)
	}

	if err := db.Model(&user).Update("active = ?", true).Error; err != nil {
		return helpers.NewError(http.StatusInternalServerError, "update_user_failed", "could not update the user", err)
	}

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
	var users []*models.User
	db.Find(&users)

	return users, nil
}

// CountUsers allows to count all users
func (db *postgres) CountUsers() (int, error) {
	return 0, nil
}

// UserExists allows to know if a user exists through his mail
func (db *postgres) UserExists(userEmail string) (bool, error) {
	var user models.User
	if err := db.Where("email = ?", userEmail).First(&user).Error; err != nil {
		return true, helpers.NewError(http.StatusNotFound, "user_exists", "found the user", err)
	}
	fmt.Println("User:", user)
	return false, nil
}
