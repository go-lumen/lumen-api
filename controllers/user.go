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
	user, err := store.GetUserById(c, c.Param("id"))

	if err != nil {
		c.AbortWithError(http.StatusNotFound, helpers.ErrorWithCode("user_not_found", "The user does not exist", err))
		return
	}

	c.JSON(http.StatusOK, user.Sanitize("", ""))
}

// GetUserMe from id (in request)
func (uc UserController) GetUserMe(c *gin.Context) {
	storeUser, exists := c.Get(store.CurrentKey)
	loggedUser := storeUser.(*models.User)
	if exists {
		user, err := store.GetUserById(c, loggedUser.Id)

		if err != nil {
			c.AbortWithError(http.StatusNotFound, helpers.ErrorWithCode("user_not_found", "The user does not exist", err))
			return
		}

		c.JSON(http.StatusOK, user.Sanitize("", ""))
	} else {
		c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("user_not_logged", "The user does not exist", nil))
	}
}

// GetUserDetails from id (in request)
func (uc UserController) GetUserDetails(c *gin.Context) {
	storeUser, exists := c.Get(store.CurrentKey)
	loggedUser := storeUser.(*models.User)
	if exists {
		user, err := store.GetUserById(c, loggedUser.Id)

		if err != nil {
			c.AbortWithError(http.StatusNotFound, helpers.ErrorWithCode("user_not_found", "The user does not exist", err))
			return
		}

		c.JSON(http.StatusOK, user.Detail("", ""))
	} else {
		c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("user_not_logged", "The user does not exist", nil))
	}
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

	databaseUser, err := store.GetUser(c, params.M{"email": user.Email})
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	apiUrl := `https://` + config.GetString(c, "api_url") + `/v1/users/` + string(user.Id) + `/activate/` + user.ActivationKey
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

	c.JSON(http.StatusCreated, user.Sanitize("", ""))
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

// ActivateUser to activate a user (usually via mail)
func (uc UserController) ActivateUser(c *gin.Context) {
	if err := store.ActivateUser(c, c.Param("activationKey"), c.Param("id")); err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	//c.JSON(http.StatusOK, nil)

	/*user, err := store.GetUserById(c, c.Param("id"))
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
func (uc UserController) ResetPasswordRequest(c *gin.Context) {
	databaseUser, err := store.GetUser(c, params.M{"email": c.Param("email")})
	if err != nil {
		c.AbortWithError(http.StatusNotFound, helpers.ErrorWithCode("user_not_found", "The user does not exist", err))
		return
	}

	resetKey := helpers.RandomString(40)

	databaseUser.ResetKey = resetKey
	if err = store.UpdateUser(c, string(databaseUser.Id), databaseUser); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	apiUrl := `https://` + config.GetString(c, "api_url") + `/v1/users/reset_password/` + string(databaseUser.Id) + `/` + resetKey
	frontUrl := config.GetString(c, "front_url")
	appName := config.GetString(c, "mail_sender_name")

	s := services.GetEmailSender(c)

	err = s.SendResetEmail(databaseUser, apiUrl, appName, frontUrl)
	if err != nil {
		logrus.Infoln(err)
		c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("mail_sending_error", "Error when sending mail", err))
		return
	}

	c.JSON(http.StatusCreated, databaseUser.Sanitize("", ""))
}

// ResetPasswordResponse allows the user to change his password
func (uc UserController) ResetPasswordResponse(C *gin.Context) {

}
