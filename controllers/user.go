package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-lumen/lumen-api/config"
	"github.com/go-lumen/lumen-api/helpers"
	"github.com/go-lumen/lumen-api/models"
	"github.com/go-lumen/lumen-api/services"
	"github.com/go-lumen/lumen-api/store"
	"github.com/go-lumen/lumen-api/utils"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// UserController holds all controller functions related to the user entity
type UserController struct {
	BaseController
}

// NewUserController instantiates of the controller
func NewUserController() UserController {
	return UserController{}
}

// GetUser allows to retrieve a user from id (in context)
// @Summary Retrieves user based on given ID
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Router /v1/users/{id} [get]
// @Security ApiKeyAuth
// @Tags User
func (uc UserController) GetUser(c *gin.Context) {
	ctx := store.AuthContext(c)
	if loggedUser, loggedGroup, ok := uc.LoggedUser(c); ok {
		switch loggedGroup.GetRole() {
		case store.RoleGod:
			user, err := models.GetUser(ctx, uc.ParamID(c))

			if err != nil {
				uc.AbortWithError(c, helpers.ErrorResourceNotFound(err))
				return
			}
			c.JSON(http.StatusOK, user)
		case store.RoleAdmin, store.RoleUser:
			user, err := models.GetUser(ctx, uc.ParamID(c))
			if user.GroupID.Hex() == loggedUser.GetGroupID() {
				c.JSON(http.StatusOK, user)
				return
			}
			uc.AbortWithError(c, helpers.ErrorResourceNotFound(err))
			return

		case store.RoleCustomer:
			uc.AbortWithError(c, helpers.ErrorUserUnauthorized)
			return
		}
	}
}

// GetUserMe from id (in request)
func (uc UserController) GetUserMe(c *gin.Context) {
	ctx := store.AuthContext(c)
	if user, group, ok := uc.LoggedUser(c); ok {
		organization, err := models.GetOrganization(ctx, utils.ParamID(group.GetOrgID()))
		if err != nil {
			uc.AbortWithError(c, helpers.ErrorResourceNotFound(err))
			return
		}

		var dbUser models.User
		if uc.ErrorInternal(c, ctx.Store.Find(ctx, store.ID(user.GetID()), &dbUser)) {
			return
		}

		c.JSON(http.StatusOK, dbUser.Sanitize(ctx.Role, ctx.Group.GetName(), organization.ID.Hex(), organization.Name))
	} else {
		uc.AbortWithError(c, helpers.ErrorUserUnauthorized)
	}
}

// ChangeUserGroup from id (in request)
func (uc UserController) ChangeUserGroup(c *gin.Context) {
	ctx := store.AuthContext(c)
	if !uc.ShouldBeLogged(ctx) {
		return
	}
	desiredGroupID := c.Param("groupID")

	userOrganization, err := models.GetOrganization(ctx, utils.ParamID(ctx.Group.GetOrgID()))
	if err != nil {
		uc.AbortWithError(c, helpers.ErrorResourceNotFound(err))
		return
	}
	desiredGroup, err := models.GetGroup(ctx, utils.ParamID(desiredGroupID))
	if err != nil {
		uc.AbortWithError(c, helpers.ErrorResourceNotFound(err))
		return
	}

	var dbUser models.User
	if uc.ErrorInternal(c, ctx.Store.Find(ctx, store.ID(ctx.User.GetID()), &dbUser)) {
		return
	}

	switch ctx.Role {
	case store.RoleGod:
		objID, _ := primitive.ObjectIDFromHex(desiredGroupID)
		dbUser.GroupID = objID
		err := models.UpdateUser(ctx, dbUser.ID.Hex(), &dbUser)
		if err != nil {
			uc.AbortWithError(c, helpers.ErrorResourceNotFound(err))
			return
		}
		c.JSON(http.StatusOK, dbUser)
	case store.RoleAdmin, store.RoleUser:
		if desiredGroup.OrganizationID == userOrganization.ID {
			objID, _ := primitive.ObjectIDFromHex(desiredGroupID)
			dbUser.GroupID = objID
			err := models.UpdateUser(ctx, dbUser.ID.Hex(), &dbUser)
			if err != nil {
				uc.AbortWithError(c, helpers.ErrorResourceNotFound(err))
				return
			}
			c.JSON(http.StatusOK, dbUser)
		} else {
			uc.AbortWithError(c, helpers.ErrorUserUnauthorized)
		}
	case store.RoleCustomer:
		uc.AbortWithError(c, helpers.ErrorUserUnauthorized)
	}
}

// CreateUser to create a new user
func (uc UserController) CreateUser(c *gin.Context) {
	ctx := store.AuthContext(c)
	user := &models.User{}

	if err := c.BindJSON(user); err != nil {
		uc.AbortWithError(c, helpers.ErrorInvalidInput(err))
		return
	}

	userGroup, err := models.GetGroup(ctx, utils.ParamID(user.GroupID.Hex()))
	if err != nil {
		uc.AbortWithError(c, helpers.ErrorInvalidInput(err))
		return
	}

	organization, err := models.GetOrganization(ctx, utils.ParamID(userGroup.OrganizationID.Hex()))
	if err != nil {
		uc.AbortWithError(c, helpers.ErrorInvalidInput(err))
		return
	}

	user.LastModification = time.Now().Unix()
	err = models.CreateUser(ctx, user)
	if err != nil {
		uc.AbortWithError(c, helpers.ErrorInvalidInput(err))
		return
	}

	dbUser, err := models.GetUser(ctx, bson.M{"email": user.Email})
	utils.CheckErr(err)

	encodedKey := []byte(config.GetString(c, "rsa_private"))
	privateKey, err := helpers.GetRSAPrivateKey(encodedKey)
	utils.CheckErr(err)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub":           "define",
		"userId":        dbUser.ID,
		"userKey":       dbUser.Key,
		"userEmail":     dbUser.Email,
		"userFirstname": dbUser.FirstName,
		"userLastname":  dbUser.LastName,
	})
	tokenStr, err := token.SignedString(privateKey)
	utils.CheckErr(err)

	apiURL := config.GetString(c, "http_scheme") + `://` + config.GetString(c, "front_url") + `/update-account/` + tokenStr
	frontURL := config.GetString(c, "front_url")
	appName := config.GetString(c, "mail_sender_name")

	s := services.GetEmailSender(c)

	err = s.SendActivationEmail(dbUser, apiURL, appName, frontURL)
	if err != nil {
		logrus.Infoln(err)
		c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("mail_sending_error", "ErrorInternal when sending mail", err))
		return
	}

	c.JSON(http.StatusCreated, user.Sanitize(userGroup.Role, userGroup.Name, organization.ID.Hex(), organization.Name))
}

// DeleteUser to delete an existing user
func (uc UserController) DeleteUser(c *gin.Context) {
	ctx := store.AuthContext(c)
	err := models.DeleteUser(ctx, c.Param("id"))

	if err != nil {
		uc.AbortWithError(c, helpers.ErrorInternal(err))
		return
	}

	c.JSON(http.StatusOK, nil)
}

// GetUsers to get all users
func (uc UserController) GetUsers(c *gin.Context) {
	ctx := store.AuthContext(c)
	if !uc.ShouldBeLogged(ctx) {
		return
	}

	dbGroups, err := models.GetGroups(ctx, bson.M{})
	utils.CheckErr(err)

	switch ctx.Role {
	case store.RoleGod:
		users, err := models.GetUsers(ctx, bson.M{})
		for _, user := range users {
			group, err := models.FindGroup(dbGroups, user.GroupID.Hex())
			if err == nil {
				objID, _ := primitive.ObjectIDFromHex(group.Name)
				user.GroupID = objID
			}
		}
		if err != nil {
			uc.AbortWithError(c, helpers.ErrorResourceNotFound(err))
			return
		}
		c.JSON(http.StatusOK, users)
	case store.RoleAdmin, store.RoleUser:
		users, err := models.GetUsers(ctx, bson.M{"group_id": ctx.User.GetGroupID()})
		for _, user := range users {
			group, err := models.FindGroup(dbGroups, user.GroupID.Hex())
			if err == nil {
				objID, _ := primitive.ObjectIDFromHex(group.Name)
				user.GroupID = objID
			}
		}
		if err != nil {
			uc.AbortWithError(c, helpers.ErrorResourceNotFound(err))
			return
		}
		c.JSON(http.StatusOK, users)
	case store.RoleCustomer:
		uc.AbortWithError(c, helpers.ErrorUserUnauthorized)
	}
}

// ResetPasswordRequest allows to send the user an email to reset his password
func (UserController) ResetPasswordRequest(c *gin.Context) {
	ctx := store.AuthContext(c)
	dbUser, err := models.GetUser(ctx, bson.M{"email": c.Param("email")})
	if err != nil {
		utils.Log(c, "warn", "ResetPasswordRequest failed with input:", c.Param("email"))
		return
	}

	encodedKey := []byte(config.GetString(c, "rsa_private"))
	privateKey, err := helpers.GetRSAPrivateKey(encodedKey)
	utils.CheckErr(err)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub":           "define",
		"userId":        dbUser.ID,
		"userKey":       dbUser.Key,
		"userEmail":     dbUser.Email,
		"userFirstname": dbUser.FirstName,
		"userLastname":  dbUser.LastName,
	})
	tokenStr, err := token.SignedString(privateKey)
	utils.CheckErr(err)

	apiURL := `https://` + config.GetString(c, "front_url") + `/update-account/` + tokenStr
	frontURL := config.GetString(c, "front_url")
	appName := config.GetString(c, "mail_sender_name")

	s := services.GetEmailSender(c)

	err = s.SendResetEmail(dbUser, apiURL, appName, frontURL)
	if err != nil {
		utils.Log(c, "info", "ErrorInternal when sending mail", err)
		c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("mail_sending_error", "ErrorInternal when sending mail", err))
		return
	}
}

// UpdateUser allows to update a user
func (uc UserController) UpdateUser(c *gin.Context) {
	ctx := store.AuthContext(c)

	type UserInput struct {
		ID        string `form:"id" json:"id"`
		Key       string `form:"key" json:"key"`
		Password  string `form:"password" json:"password"`
		Firstname string `form:"firstname" json:"firstname"`
		Lastname  string `form:"lastname" json:"lastname"`
	}
	var userInput UserInput
	if err := c.ShouldBind(&userInput); err != nil {
		uc.AbortWithError(c, helpers.ErrorInvalidInput(err))
		return
	}

	dbUser, err := models.GetUser(ctx, utils.ParamID(userInput.ID))
	if err != nil {
		uc.AbortWithError(c, helpers.ErrorResourceNotFound(err))
		return
	}

	if userInput.Key != dbUser.Key {
		uc.AbortWithError(c, helpers.ErrorUserUnauthorized)
		return
	}

	dbUser.Key = helpers.RandomString(40)

	if dbUser.Status == "created" {
		dbUser.Status = "activated"
	} else if dbUser.Status == "activated" {
		dbUser.Status = "modified"
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	if err != nil {
		uc.AbortWithError(c, helpers.ErrorInternal(err))
		return
	}
	dbUser.Password = string(hashedPassword)

	if userInput.Firstname != "" {
		dbUser.FirstName = userInput.Firstname
	}
	if userInput.Lastname != "" {
		dbUser.LastName = userInput.Lastname
	}

	dbUser.LastModification = time.Now().Unix()
	err = models.UpdateUser(ctx, dbUser.ID.Hex(), dbUser)
	utils.CheckErr(err)

	c.JSON(http.StatusOK, "User well modified")
}
