package server

import (
	"fmt"
	"github.com/go-lumen/lumen-api/helpers/params"
	"github.com/go-lumen/lumen-api/models"
	"github.com/go-lumen/lumen-api/store/mongodb"
	"github.com/go-lumen/lumen-api/store/mysql"
	"github.com/go-lumen/lumen-api/store/postgresql"
	"github.com/go-lumen/lumen-api/utils"
	"github.com/sirupsen/logrus"
)

// SetupMongoSeeds creates the first user
func (a *API) SetupMongoSeeds() error {
	store := mongodb.New(a.MongoDatabase, a.Config.GetString("mongo_db_name"), nil)

	user := &models.User{
		FirstName: a.Config.GetString("admin_firstname"),
		LastName:  a.Config.GetString("admin_lastname"),
		Password:  a.Config.GetString("admin_password"),
		Email:     a.Config.GetString("admin_email"),
		Phone:     a.Config.GetString("admin_phone"),
		Role: "admin",
	}

	userExists, _, err := store.UserExists(user.Email)
	if userExists {
		utils.Log(nil, "warn", `Seed user already exists`, err)
	} else {
		utils.Log(nil, "info", "User doesn't exists already")
	}

	err = store.CreateUser(user)
	if err != nil {
		utils.Log(nil, "warn", `Error when creating user:`, err)
		user, _ = store.GetUser(params.M{"email": a.Config.GetString("admin_email")})
	} else {
		utils.Log(nil, "info", "User well created")

		err = store.ActivateUser(user.ActivationKey, user.Id)
		if err != nil {
			utils.Log(nil, "warn", `Error when activating user`, err)
		} else {
			utils.Log(nil, "info", "User well activated")
		}

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
		Role:      "admin",
	}
	userExists, userFound, err := store.UserExists(user.Email)
	if userExists {
		logrus.Infoln(`Seed user already exists`, err, userFound)
	} else {
		if err := store.CreateUser(user); err != nil {
			logrus.Warnln(`Error when creating user:`, err)
		}
	}

	dbUser, err := store.GetUser(params.M{"email": a.Config.GetString("admin_email")})
	if err != nil {
		logrus.Warnln(err)
	}
	fmt.Println("Found user", dbUser.Id, ":", dbUser)

	if err := store.ActivateUser(dbUser.ActivationKey, /*strconv.Itoa(dbUser.Id)*/ dbUser.Email); err != nil {
		logrus.Warnln(`Error when activating user`, err)
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
		Role:      "admin",
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
