package postgresql

import (
	"fmt"
	"github.com/go-lumen/lumen-api/helpers/params"
	"github.com/go-lumen/lumen-api/models"
)

// CreateUser checks if user already exists, and if not, creates it
func (db *postgres) CreateUser(user *models.User) error {
	if err := db.Create(&user).Error; err != nil {
		// handle err
	}

	return nil
}

// FindUserById allows to retrieve a user by its id
func (db *postgres) FindUserById(id string) (*models.User, error) {
	fmt.Println("finding user:", id)
	user := &models.User{}

	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		//handle error
	}

	return user, nil
}

// FindUser allows to retrieve a user by its characteristics
func (db *postgres) FindUser(params params.M) (res *models.User, err error) {
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
	fmt.Println("finding user:", userEmail)
	user := &models.User{Email: userEmail}

	db.Where("email = ?", userEmail).Find(&user)

	fmt.Println("User:", user)

	return false, nil
}
