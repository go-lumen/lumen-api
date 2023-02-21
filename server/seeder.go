package server

import (
	"github.com/gin-gonic/gin"
	"github.com/go-lumen/lumen-api/models"
	"github.com/go-lumen/lumen-api/store"
	"github.com/go-lumen/lumen-api/store/mongodb"
	"github.com/go-lumen/lumen-api/store/postgresql"
	"github.com/go-lumen/lumen-api/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// SetupMongoSeeds creates the first user
func (a *API) SetupMongoSeeds() error {
	s := mongodb.New(nil, a.MongoDatabase, a.Config.GetString("mongo_db_name"))
	ctx := store.NewGodContext(s)

	//Mails: 0.10$/1000         Texts: 0.05-0.10$/1       WiFi: 5$/1000

	organization := &models.Organization{
		Name:      a.Config.GetString("project_name"),
		LogoURL:   "",
		Siret:     0,
		VATNumber: "",
		Tokens:    100000000000,
		Parent:    "",
	}
	dbOrga := &models.Organization{}
	if err := s.Find(ctx, bson.M{"name": organization.Name}, dbOrga); err == nil {
		utils.Log(nil, "warn", `Organization:`, organization.Name, `already exists`)
	} else if err := s.Create(ctx, "", organization); err != nil {
		utils.Log(nil, "err", `ErrorInternal when creating organization:`, err)
	} else {
		utils.Log(nil, "info", `Organization:`, organization.Name, `well created`)
	}

	group := &models.Group{
		Name:           a.Config.GetString("project_name") + " superadmin",
		Role:           store.RoleGod,
		OrganizationID: organization.ID,
	}
	if err := s.Find(ctx, bson.M{"name": group.Name}, group); err == nil {
		utils.Log(nil, "warn", `Group:`, group.Name, `already exists`)
	} else if err := s.Create(ctx, "", group); err != nil {
		utils.Log(nil, "err", `ErrorInternal when creating group:`, group.Name, err)
	} else {
		utils.Log(nil, "info", "Group well created")
	}

	user := &models.User{
		FirstName: a.Config.GetString("admin_firstname"),
		LastName:  a.Config.GetString("admin_lastname"),
		Password:  a.Config.GetString("admin_password"),
		Email:     a.Config.GetString("admin_email"),
		Phone:     a.Config.GetString("admin_phone"),
		GroupID:   group.ID,
		Balance:   123.45,
	}

	userExists, user, err := models.UserExists(ctx, user.Email)
	if userExists {
		utils.Log(nil, "warn", `Seed user already exists`)
	} else {
		utils.Log(nil, "info", "User doesn't exists already")
	}

	err = models.CreateUser(ctx, user)
	if err != nil {
		utils.Log(nil, "warn", `ErrorInternal when creating user:`, err)
		user, _ = models.GetUser(ctx, bson.M{"email": a.Config.GetString("admin_email")})
	} else {
		utils.Log(nil, "info", "User well created")
	}

	err = models.ActivateUser(ctx, user.Key, user.ID)
	if err != nil {
		utils.Log(nil, "warn", `ErrorInternal when activating user`, err)
	} else {
		utils.Log(nil, "info", "User well activated")
	}

	return nil
}

// SetupPostgreSeeds creates the first user
func (a *API) SetupPostgreSeeds() error {
	utils.Log(nil, "info", "Setup postgre seeds")
	s := postgresql.New(&gin.Context{}, a.PostgreDatabase, a.Config.GetString("POSTGRES_DB_NAME"))
	//ctx := store.NewGodContext(s)

	organization := &models.Organization{
		Name: a.Config.GetString("project_name"),
	}
	s.Create(a.Context, "", organization)

	adminGroup := &models.Group{
		Name:           a.Config.GetString("project_name") + " Admin",
		Role:           "god",
		OrganizationID: organization.ID,
	}
	s.Create(a.Context, "", adminGroup)

	adminUser := &models.User{
		FirstName: a.Config.GetString("admin_firstname"),
		LastName:  a.Config.GetString("admin_lastname"),
		Password:  a.Config.GetString("admin_password"),
		Email:     a.Config.GetString("admin_email"),
		Phone:     a.Config.GetString("admin_phone"),
		Status:    "activated",
		GroupID:   adminGroup.ID,
		Balance:   123.45,
	}
	s.GetOrCreateUser(adminUser)
	/*err := models.ActivateUser(ctx, adminUser.Key, adminUser.ID)
	if err != nil {
		utils.Log(nil, "warn", `ErrorInternal when activating user`, err)
	} else {
		utils.Log(nil, "info", "User well activated")
	}*/

	user1 := &models.User{
		FirstName: a.Config.GetString("user1_firstname"),
		LastName:  a.Config.GetString("user1_lastname"),
		Password:  a.Config.GetString("user1_password"),
		Email:     a.Config.GetString("user1_email"),
		GroupID:   adminGroup.ID,
	}
	s.GetOrCreateUser(user1)
	user2 := &models.User{
		FirstName: a.Config.GetString("user2_firstname"),
		LastName:  a.Config.GetString("user2_lastname"),
		Password:  a.Config.GetString("user2_password"),
		Email:     a.Config.GetString("user2_email"),
		GroupID:   adminGroup.ID,
	}
	s.GetOrCreateUser(user2)

	return nil
}
