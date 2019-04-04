package server

import (
	"github.com/adrien3d/lumen-api/models"
	"github.com/adrien3d/lumen-api/store/mongodb"
	"github.com/sirupsen/logrus"
)

// SetupSeeds creates the first user
func (a *API) SetupSeeds() error {
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
