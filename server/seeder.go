package server

import (
	"github.com/adrien3d/base-api/models"
	"github.com/adrien3d/base-api/store/mongodb"
)

func (a *API) SetupSeeds() error {
	store := mongodb.New(a.Database)

	//Mails: 0.10$/1000         Texts: 0.05-0.10$/1       WiFi: 5$/1000

	user := &models.User{
		Firstname: "Adrien",
		Lastname:  "Chapelet",
		Password:  "demo",
		Email:     "adrien@plugblocks.com",
		Phone:     "+33671174203",
		Admin:     true,
	}

	store.CreateUser(user)
	store.ActivateUser(user.ActivationKey, user.Id)
	return nil
}
