package mysql

import (
	"fmt"
	"go-lumen/lumen-api/helpers/params"
	"go-lumen/lumen-api/models"
	"go-lumen/lumen-api/utils"
)

// CreateUser checks if user already exists, and if not, creates it
func (db *mysql) CreateUser(user *models.User) error {
	/*
	   CREATE TABLE `userinfo` (
	       `uid` INT(10) NOT NULL AUTO_INCREMENT,
	       `username` VARCHAR(64) NULL DEFAULT NULL,
	       `departname` VARCHAR(64) NULL DEFAULT NULL,
	       `created` DATE NULL DEFAULT NULL,
	       PRIMARY KEY (`uid`)
	   );
	*/
	return nil
}

// FindUserById allows to retrieve a user by its id
func (db *mysql) FindUserById(id string) (*models.User, error) {
	return nil, nil
}

// FindUser allows to retrieve a user by its characteristics
func (db *mysql) FindUser(params params.M) (*models.User, error) {
	var firstName string
	fmt.Println("finding user:", params)
	var rows, err = db.Query("SELECT * FROM users WHERE email = $1", params)
	utils.CheckErr(err)
	utils.CheckErr(rows.Scan(&firstName))

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
