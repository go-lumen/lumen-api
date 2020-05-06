package mysql

import (
	"fmt"
	"github.com/go-lumen/lumen-api/helpers/params"
	"github.com/go-lumen/lumen-api/models"
	"github.com/go-lumen/lumen-api/utils"
)

// CreateUser checks if user already exists, and if not, creates it
func (db *Mysql) CreateUser(user *models.User) error {
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

// FindUserByID allows to retrieve a user by its id
func (db *Mysql) FindUserByID(id string) (*models.User, error) {
	return nil, nil
}

// FindUser allows to retrieve a user by its characteristics
func (db *Mysql) FindUser(params params.M) (*models.User, error) {
	var firstName string
	fmt.Println("finding user:", params)
	var rows, err = db.Query("SELECT * FROM users WHERE email = $1", params)
	utils.CheckErr(err)
	utils.CheckErr(rows.Scan(&firstName))

	return nil, nil
}

// DeleteUser allows to delete a user by its id
func (db *Mysql) DeleteUser(user *models.User, userID string) error {
	return nil
}

// ActivateUser allows to activate a user by its id
func (db *Mysql) ActivateUser(activationKey string, id string) error {
	return nil
}

// ChangeLanguage allows to change a user language by its id
func (db *Mysql) ChangeLanguage(id string, language string) error {
	return nil
}

// UpdateUser allows to update one or more user characteristics
func (db *Mysql) UpdateUser(userID string, params params.M) error {
	return nil
}

// GetUsers allows to get all users
func (db *Mysql) GetUsers() ([]*models.User, error) {
	return nil, nil
}

// CountUsers allows to count all users
func (db *Mysql) CountUsers() (int, error) {
	return 0, nil
}

// UserExists allows to know if a user exists through his mail
func (db *Mysql) UserExists(userEmail string) (bool, error) {
	return false, nil
}
