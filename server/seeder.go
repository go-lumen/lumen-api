package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"go-lumen/lumen-api/models"
	"go-lumen/lumen-api/store/mongodb"
	"go-lumen/lumen-api/store/mysql"
	"go-lumen/lumen-api/store/postgresql"
)

// SetupMongoSeeds creates the first user
func (a *API) SetupMongoSeeds() error {
	store := mongodb.New(a.MongoDatabase)

	//Mails: 0.10$/1000         Texts: 0.05-0.10$/1       WiFi: 5$/1000

	user := &models.User{
		FirstName: a.Config.GetString("admin_firstname"),
		LastName:  a.Config.GetString("admin_lastname"),
		Password:  a.Config.GetString("admin_password"),
		Email:     a.Config.GetString("admin_email"),
		Phone:     a.Config.GetString("admin_phone"),
		Admin:     true,
	}

	userExists, _ := store.UserExists(user.Email)
	if userExists {
		logrus.Infoln(`Seed user already exists`)
		return nil
	}

	if store.CreateUser(user) != nil {
		logrus.Warnln(`Error when creating user`)
	}

	if store.ActivateUser(user.ActivationKey, string(user.Id)) != nil {
		logrus.Warnln(`Error when activating user`)
	}

	return nil
}

// SetupPostgreSeeds creates the first user
func (a *API) SetupPostgreSeeds() error {
	fmt.Println("Setup postgre seeds")
	store := postgresql.New(a.PostgreDatabase)

	user := &models.User{
		FirstName: a.Config.GetString("admin_firstname"),
		LastName:  a.Config.GetString("admin_lastname"),
		Password:  a.Config.GetString("admin_password"),
		Email:     a.Config.GetString("admin_email"),
		Phone:     a.Config.GetString("admin_phone"),
		Admin:     true,
	}
	/*userExists, _ := store.UserExists(user.Email)
	if userExists {
		logrus.Infoln(`Seed user already exists`)
		return nil
	}*/

	if store.CreateUser(user) != nil {
		logrus.Warnln(`Error when creating user`)
	}

	if store.ActivateUser(user.ActivationKey, string(user.Id)) != nil {
		logrus.Warnln(`Error when activating user`)
	}
	fmt.Println("Checked")

	return nil
}

// SetupMySQLSeeds creates the first user
func (a *API) SetupMySQLSeeds() error {
	store := mysql.New(a.MySQLDatabase)

	user := &models.User{
		FirstName: a.Config.GetString("admin_firstname"),
		LastName:  a.Config.GetString("admin_lastname"),
		Password:  a.Config.GetString("admin_password"),
		Email:     a.Config.GetString("admin_email"),
		Phone:     a.Config.GetString("admin_phone"),
		Admin:     true,
	}

	userExists, _ := store.UserExists(user.Email)
	if userExists {
		logrus.Infoln(`Seed user already exists`)
		return nil
	}

	if store.CreateUser(user) != nil {
		logrus.Warnln(`Error when creating user`)
	}

	if store.ActivateUser(user.ActivationKey, string(user.Id)) != nil {
		logrus.Warnln(`Error when activating user`)
	}

	return nil
}
