package server

import (
	"fmt"
	"github.com/go-lumen/lumen-api/models"
	"github.com/go-lumen/lumen-api/store/mongodb"
	"github.com/go-lumen/lumen-api/store/mysql"
	"github.com/go-lumen/lumen-api/store/postgresql"
	"github.com/go-lumen/lumen-api/utils"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/sirupsen/logrus"
)

// SetupMongoSeeds creates the first user
func (a *API) SetupMongoSeeds() error {
	store := mongodb.New(a.MongoDatabase)

	//Mails: 0.10$/1000         Texts: 0.05-0.10$/1       WiFi: 5$/1000

	user := &models.User{
		FirstName: "Adrien",
		LastName:  "Chapelet",
		Password:  "adminpwd",
		Email:     "adrien@plugblocks.com",
		Phone:     "+33671174203",
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

	if store.ActivateUser(user.ActivationKey, user.Id) != nil {
		logrus.Warnln(`Error when activating user`)
	}

	return nil
}

func setupPostgreTables(db *pg.DB, models []interface{}) error {
	for _, model := range models {
		err := db.DropTable(model, &orm.DropTableOptions{
			IfExists: true,
			Cascade:  true,
		})
		if err != nil {
			return err
		}

		err = db.CreateTable(model, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetupPostgreSeeds creates the first user
func (a *API) SetupPostgreSeeds() error {
	store := postgresql.New(a.PostgreDatabase)

	mods := []interface{}{
		(*models.User)(nil),
	}
	utils.CheckErr(setupPostgreTables(store.DB, mods))

	user := &models.User{
		FirstName: "Adrien",
		LastName:  "Chapelet",
		Password:  "adminpwd",
		Email:     "adrien@plugblocks.com",
		Phone:     "+33671174203",
		Admin:     true,
	}
	utils.CheckErr(user.BeforeCreate())

	fmt.Println("Checking")
	userExists, _ := store.UserExists(user.Email)
	if userExists {
		logrus.Infoln(`Seed user already exists`)
		return nil
	}

	if store.CreateUser(user) != nil {
		logrus.Warnln(`Error when creating user`)
	}

	if store.ActivateUser(user.ActivationKey, user.Id) != nil {
		logrus.Warnln(`Error when activating user`)
	}
	fmt.Println("Checked")

	return nil
}

// SetupMySQLSeeds creates the first user
func (a *API) SetupMySQLSeeds() error {
	store := mysql.New(a.MySQLDatabase)

	user := &models.User{
		FirstName: "Adrien",
		LastName:  "Chapelet",
		Password:  "adminpwd",
		Email:     "adrien@plugblocks.com",
		Phone:     "+33671174203",
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

	if store.ActivateUser(user.ActivationKey, user.Id) != nil {
		logrus.Warnln(`Error when activating user`)
	}

	return nil
}
