package controllers

import (
	"fmt"
	"go-lumen/lumen-api/utils"
	"net/http"

	"go-lumen/lumen-api/config"
	"go-lumen/lumen-api/helpers"
	"go-lumen/lumen-api/helpers/params"
	"go-lumen/lumen-api/models"
	"go-lumen/lumen-api/store"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

// AuthController structure
type AuthController struct {
}

// NewAuthController instantiates of the controller
func NewAuthController() AuthController {
	return AuthController{}
}

// UserAuthentication for authenticating user
func (ac AuthController) UserAuthentication(c *gin.Context) {
	userInput := models.User{}
	if err := c.Bind(&userInput); err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data", err))
		return
	}

	user, err := store.FindUser(c, params.M{"email": userInput.Email})
	if err != nil {
		c.AbortWithError(http.StatusNotFound, helpers.ErrorWithCode("user_does_not_exist", "User does not exist", err))
		return
	}

	fmt.Println("Comparing", string([]byte(user.Password)), "and", string([]byte(userInput.Password)))
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password))
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("incorrect_password", "Password is not correct", err))
		fmt.Println("CompareHashAndPassword err:", err)
		return
	}

	if !user.Active {
		c.AbortWithError(http.StatusNotFound, helpers.ErrorWithCode("user_needs_activation", "User needs to be activated via email", nil))
		return
	}

	//Read base64 private key
	encodedKey := []byte(config.GetString(c, "rsa_private"))
	accessToken, err := helpers.GenerateAccessToken(encodedKey, string(user.Id), user.LastPasswordUpdate)
	if err != nil {
		utils.CheckErr(err)
		c.AbortWithError(http.StatusInternalServerError, helpers.ErrorWithCode("token_generation_failed", "Could not generate the access token", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": accessToken, "user": user.Sanitize()})
}
