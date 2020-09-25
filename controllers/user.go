package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-lumen/lumen-api/config"
	"github.com/go-lumen/lumen-api/helpers"
	"github.com/go-lumen/lumen-api/helpers/params"
	"github.com/go-lumen/lumen-api/models"
	"github.com/go-lumen/lumen-api/services"
	"github.com/go-lumen/lumen-api/store"
	"github.com/go-lumen/lumen-api/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// UserController holds all controller functions related to the user entity
type UserController struct{}

// NewUserController instantiates of the controller
func NewUserController() UserController {
	return UserController{}
}

// GetUser from id (in context)
func (UserController) GetUser(c *gin.Context) {
	user, err := store.GetUserByID(c, c.Param("id"))

	if err != nil {
		c.AbortWithError(http.StatusNotFound, helpers.ErrorWithCode("user_not_found", "The user does not exist", err))
		return
	}

	c.JSON(http.StatusOK, user.Sanitize("", "", ""))
}

// GetUserMe from id (in request)
func (UserController) GetUserMe(c *gin.Context) {
	storeUser, exists := c.Get(store.CurrentKey)
	loggedUser := storeUser.(*models.User)
	if exists {
		user, err := store.GetUserByID(c, loggedUser.ID)

		if err != nil {
			c.AbortWithError(http.StatusNotFound, helpers.ErrorWithCode("user_not_found", "The user does not exist", err))
			return
		}

		c.JSON(http.StatusOK, user.Sanitize("", "", ""))
	} else {
		c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("user_not_logged", "The user does not exist", nil))
	}
}

// CreateUser to create a new user
func (UserController) CreateUser(c *gin.Context) {
	user := &models.User{}

	if err := c.BindJSON(user); err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data", err))
		return
	}

	user.LastModification = time.Now().Unix()
	err := store.CreateUser(c, user)
	utils.CheckErr(err)

	dbUser, err := store.GetUser(c, params.M{"email": user.Email})
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

	apiURL := `https://` + config.GetString(c, "front_url") + `/update-account/` + tokenStr
	frontURL := config.GetString(c, "front_url")
	appName := config.GetString(c, "mail_sender_name")

	s := services.GetEmailSender(c)

	err = s.SendActivationEmail(dbUser, apiURL, appName, frontURL)
	if err != nil {
		logrus.Infoln(err)
		c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("mail_sending_error", "Error when sending mail", err))
		return
	}

	c.JSON(http.StatusCreated, user.Sanitize("", "", ""))
}

// DeleteUser to delete an existing user
func (UserController) DeleteUser(c *gin.Context) {
	err := store.DeleteUser(c, c.Param("id"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, nil)
}

// GetUsers to get all users
func (UserController) GetUsers(c *gin.Context) {
	_, exists := c.Get(store.CurrentKey)
	if exists {
		users, err := store.GetUsers(c, "")
		if err != nil {
			helpers.NewError(http.StatusNotFound, "users_not_found", "Users not found", err)
		}
		c.JSON(http.StatusOK, users)
	}
}

// ResetPasswordRequest allows to send the user an email to reset his password
func (UserController) ResetPasswordRequest(c *gin.Context) {
	dbUser, err := store.GetUser(c, params.M{"email": c.Param("email")})
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
		logrus.Infoln(err)
		c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("mail_sending_error", "Error when sending mail", err))
		return
	}
}

func (UserController) UpdateUser(c *gin.Context) {
	type UserInput struct {
		ID        string `form:"id" json:"id"`
		Key       string `form:"key" json:"key"`
		Password  string `form:"password" json:"password"`
		Firstname string `form:"firstname" json:"firstname"`
		Lastname  string `form:"lastname" json:"lastname"`
	}
	var userInput UserInput
	if c.ShouldBind(&userInput) != nil {
		return
	}

	dbUser, err := store.GetUserByID(c, userInput.ID)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, helpers.ErrorWithCode("user_not_found", "The user does not exist", err))
		return
	}

	if userInput.Key != dbUser.Key {
		c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("invalid_key", "Invalid key", err))
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
		c.AbortWithError(http.StatusInternalServerError, helpers.ErrorWithCode("encryption_failed", "Failed to generate the crypted password", err))
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
	err = store.UpdateUser(c, dbUser.ID, dbUser)
	utils.CheckErr(err)

	c.JSON(http.StatusOK, "User well modified")
}
