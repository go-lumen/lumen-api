package server

import (
	"fmt"
	"github.com/go-lumen/lumen-api/helpers/params"
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
	} else if err := s.Create(ctx, organization); err != nil {
		utils.Log(nil, "err", `ErrorInternal when creating organization:`, err)
	} else {
		utils.Log(nil, "info", `Organization:`, organization.Name, `well created`)
	}
	/*
		poiA := &models.POI{
			Name:    "Fline",
			Code:    "Fline-Paris",
			Type:    "office",
			Country: "FR",
			Location: models.Location{
				Type:        "Point",
				Coordinates: []float64{2.2945, 48.8582},
			},
			OrganizationID: organization.ID,
			MacAdresses:    nil,
		}
		dbPoi := &models.POI{}
		if err := s.Find(ctx, bson.M{"name": poiA.Name}, dbPoi); err == nil {
			utils.Log(nil, "warn", `POI:`, poiA.Name, `already exists`)
		} else if err := s.Create(ctx, poiA); err != nil {
			utils.Log(nil, "err", `ErrorInternal when creating POI A:`, err)
		} else {
			utils.Log(nil, "info", `POI:`, poiA.Name, `well created`)
		}*/

	group := &models.Group{
		Name:           a.Config.GetString("project_name") + " superadmin",
		Role:           store.RoleGod,
		OrganizationID: organization.ID,
	}
	if err := s.Find(ctx, bson.M{"name": group.Name}, group); err == nil {
		utils.Log(nil, "warn", `Group:`, group.Name, `already exists`)
	} else if err := s.Create(ctx, group); err != nil {
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
	}

	userExists, _, err := models.UserExists(ctx, user.Email)
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

	err = models.ActivateUser(ctx, user.Key, user.ID.Hex())
	if err != nil {
		utils.Log(nil, "warn", `ErrorInternal when activating user`, err)
	} else {
		utils.Log(nil, "info", "User well activated")
	}

	ad1 := &models.Ad{
		UserID:               user.ID,
		Name:                 "Ananas IQF 10*10 (100772)",
		Category:             "frozen",
		Place:                "Kervignac",
		Quantity:             1270,
		Price:                2.70,
		Description:          "Ananas préparés à partir d'ananas frais\nCalibre : Cubes 10x10 mm (80% en masse)\nIls sont pelés, étrognés, triés, coupés en cubes puis surgelés individuellement selon le procédé IQF\nPour plus de détails => demander FT",
		PicturesURLs:         []string{"https://cdn.kreezalid.com/kreezalid/554097/catalog/8088/36/1000x1000_juno-jo-hrdnnog-y-a-unsplash_nuq3f_1243822212.jpg"},
		Certifications:       []string{"IFS", "BRC", "FSSC 22000"},
		OriginType:           "other",
		Origin:               "Costa Rica / Vietnam\n",
		VendorLinkToMaterial: "bought-to-transform",
		ExpirationDateType:   "dluo",
		ExpirationDate:       1658959200, //1665093600
		IsOrganic:            false,
		PackingType:          "other",
		Packing:              "Carton doublé d'un sac polyéthylène bleu\n",
		PackingWeight:        10,
		PalletType:           "eur",
		PalletWeight:         640,
		CuttingType:          "10*10",
		SellingReason:        "no-longer-needed",
		DeliveryAvailable:    false,
	}

	/*adExists, _, err := models.AdExists(ctx, ad1.Name)
	if adExists {
		utils.Log(nil, "warn", `Seed ad 1 already exists`)
	} else {
		utils.Log(nil, "info", "Ad 1 doesn't exists already")
	}*/

	err = models.CreateAd(ctx, ad1)
	if err != nil {
		ad1, _ = models.GetAd(ctx, bson.M{"name": ad1.Name})
		utils.Log(nil, "warn", `ErrorInternal when creating ad1:`, err)
		utils.Log(nil, "warn", `ErrorInternal when creating ad1:`, err.Error())
	} else {
		utils.Log(nil, "info", "Ad1 well created")
	}

	/*
		fleet := &models.Fleet{
			Name:               a.Config.GetString("project_name") + " v1",
			SigfoxDeviceTypeID: "5b5726889e93a1464b6e552c",
			UserID:             user.ID,
			GroupID:            group.ID,
			Status:             "available",
			Resolver:           models.TrackerHWWisoliHere,
		}
		if err := models.CreateFleet(ctx, fleet); err != nil {
			dbFleet, _ := models.GetFleet(ctx, bson.M{"sigfox_device_type_id": fleet.SigfoxDeviceTypeID})
			utils.Log(nil, "warn", `Fleet already exists, id:`, dbFleet.ID, `SigfoxDeviceTypeID:`, dbFleet.SigfoxDeviceTypeID, `AppKey:`, dbFleet.AppKey)
		} else {
			utils.Log(nil, "info", "Fleet ", fleet.Name, " well created, Sfx Device Type:", fleet.SigfoxDeviceTypeID, " appKey:", fleet.AppKey)
		}*/

	return nil
}

// SetupPostgreSeeds creates the first user
func (a *API) SetupPostgreSeeds() error {
	utils.Log(nil, "info", "Setup postgre seeds")
	store := postgresql.New(a.PostgreDatabase)

	user := &models.User{
		FirstName: a.Config.GetString("admin_firstname"),
		LastName:  a.Config.GetString("admin_lastname"),
		Password:  a.Config.GetString("admin_password"),
		Email:     a.Config.GetString("admin_email"),
		Phone:     a.Config.GetString("admin_phone"),
	}
	userExists, err := store.UserExists(user.Email)
	if userExists {
		utils.Log(nil, "warn", `Seed user already exists`, err)
	} else {
		if err := store.CreateUser(user); err != nil {
			utils.Log(nil, "warn", `Error when creating user:`, err)
		}
	}

	dbUser, err := store.GetUser(params.M{"email": a.Config.GetString("admin_email")})
	if err != nil {
		utils.Log(nil, "warn", err)
	}
	fmt.Println("Found user", dbUser.ID, ":", dbUser)

	if err := store.ActivateUser(dbUser.Key /*strconv.Itoa(dbUser.ID)*/, dbUser.Email); err != nil {
		utils.Log(nil, "warn", `Error when activating user`, err)
	}
	utils.Log(nil, "info", "Checked")

	return nil
}
