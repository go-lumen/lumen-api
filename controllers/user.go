package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-lumen/lumen-api/config"
	"github.com/go-lumen/lumen-api/helpers"
	"github.com/go-lumen/lumen-api/helpers/params"
	"github.com/go-lumen/lumen-api/models"
	"github.com/go-lumen/lumen-api/services"
	"github.com/go-lumen/lumen-api/store"
	"github.com/sirupsen/logrus"
	"net/http"
)

// UserController holds all controller functions related to the user entity
type UserController struct{}

// NewUserController instantiates of the controller
func NewUserController() UserController {
	return UserController{}
}

// GetUser from id (in context)
func (uc UserController) GetUser(c *gin.Context) {
	user, err := store.FindUserById(c, c.Param("id"))

	if err != nil {
		c.AbortWithError(http.StatusNotFound, helpers.ErrorWithCode("user_not_found", "The user does not exist", err))
		return
	}

	c.JSON(http.StatusOK, user.Sanitize())
}

// CreateUser to create a new user
func (uc UserController) CreateUser(c *gin.Context) {
	user := &models.User{}

	if err := c.BindJSON(user); err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data", err))
		return
	}

	if err := store.CreateUser(c, user); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	databaseUser, err := store.FindUser(c, params.M{"email": user.Email})
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	apiUrl := `https://` + config.GetString(c, "api_url") + `/v1/users/` + user.Id + `/activate/` + user.ActivationKey
	frontUrl := config.GetString(c, "front_url")
	appName := config.GetString(c, "mail_sender_name")

	s := services.GetEmailSender(c)

	fmt.Println("DB User:", databaseUser)
	err = s.SendActivationEmail(databaseUser, apiUrl, appName, frontUrl)
	if err != nil {
		logrus.Infoln(err)
		c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("mail_sending_error", "Error when sending mail", err))
		return
	}

	c.JSON(http.StatusCreated, user.Sanitize())
}

// DeleteUser to delete an existing user
func (uc UserController) DeleteUser(c *gin.Context) {
	err := store.DeleteUser(c, c.Param("id"))

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, nil)
}

// ChangeLanguage changes user language
func (uc UserController) ChangeLanguage(c *gin.Context) {
	if err := store.ChangeLanguage(c, c.Param("id"), c.Param("language")); err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, nil)
}

// ActivateUser to activate a user (usually via mail)
func (uc UserController) ActivateUser(c *gin.Context) {
	if err := store.ActivateUser(c, c.Param("activationKey"), c.Param("id")); err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	//c.JSON(http.StatusOK, nil)

	/*user, err := store.FindUserById(c, c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusNotFound, helpers.ErrorWithCode("user_not_found", "The user does not exist", err))
		return
	}

	vars := gin.H{
		"User":    user,
		"AppName": config.GetString(c, "mail_sender_name"),
		"AppUrl":  config.GetString(c, "front_url"),
	}

	c.HTML(http.StatusOK, "./templates/html/page_account_activated.html", vars)*/

	c.Redirect(http.StatusMovedPermanently, "https://"+config.GetString(c, "front_url"))
}

// GetUsers to get all users
func (uc UserController) GetUsers(c *gin.Context) {
	users, err := store.GetUsers(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, users)
}

// ResetPasswordRequest allows to send the user an email to reset his password
func (uc UserController) ResetPasswordRequest(c *gin.Context) {
	databaseUser, err := store.FindUser(c, params.M{"email": c.Param("email")})
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	resetKey := helpers.RandomString(40)

	if err = store.UpdateUser(c, databaseUser.Id, params.M{"resetKey": resetKey}); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	apiUrl := `https://` + config.GetString(c, "api_url") + `/v1/users/reset_password/` + databaseUser.Id + `/` + resetKey
	frontUrl := config.GetString(c, "front_url")
	appName := config.GetString(c, "mail_sender_name")

	s := services.GetEmailSender(c)

	err = s.SendResetEmail(databaseUser, apiUrl, appName, frontUrl)
	if err != nil {
		logrus.Infoln(err)
		c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("mail_sending_error", "Error when sending mail", err))
		return
	}

	c.JSON(http.StatusCreated, databaseUser.Sanitize())
}

// ResetPasswordResponse allows the user to change his password
func (uc UserController) ResetPasswordResponse(C *gin.Context) {

}
