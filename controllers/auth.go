package controllers

import (
	"fmt"
	"github.com/go-lumen/lumen-api/config"
	"github.com/go-lumen/lumen-api/helpers"
	"github.com/go-lumen/lumen-api/helpers/params"
	"github.com/go-lumen/lumen-api/models"
	"github.com/go-lumen/lumen-api/store"
	"github.com/go-lumen/lumen-api/utils"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthController structure
type AuthController struct {
}

// NewAuthController instantiates the controller
func NewAuthController() AuthController {
	return AuthController{}
}

func returnToken(c *gin.Context, encodedKey []byte, dbUser *models.User) {
	accessToken, err := helpers.GenerateToken(encodedKey, dbUser.ID, "access", 4320) //3d in min
	if err != nil {
		utils.CheckErr(err)
		c.AbortWithError(http.StatusInternalServerError, helpers.ErrorWithCode("token_generation_failed", "Could not generate the access token", err))
		return
	}
	refreshToken, err := helpers.GenerateToken(encodedKey, dbUser.ID, "refresh", 10080) //7d in min
	if err != nil {
		utils.CheckErr(err)
		c.AbortWithError(http.StatusInternalServerError, helpers.ErrorWithCode("token_generation_failed", "Could not generate the refresh token", err))
		return
	}

	err = store.UpdateUserFields(c, dbUser.ID, params.M{"$set": params.M{"last_login": time.Now().Unix()}})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, helpers.ErrorWithCode("update_user_failed", "Could not update the user", err))
	}

	c.JSON(http.StatusOK, gin.H{"token": accessToken, "refresh_token": refreshToken, "user": dbUser.Sanitize(dbUser.Role, "", "")})
}

// UserAuthentication for authenticating user
func (ac AuthController) TokensGeneration(c *gin.Context) {
	userInput := models.User{}
	if err := c.Bind(&userInput); err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("invalid_input", "Failed to bind the body data", err))
		return
	}

	dbUser, err := store.GetUser(c, params.M{"email": userInput.Email})
	if err != nil {
		c.AbortWithError(http.StatusNotFound, helpers.ErrorWithCode("user_does_not_exist", "User does not exist", err))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(userInput.Password))
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, helpers.ErrorWithCode("incorrect_password", "Password is not correct", err))
		fmt.Println("CompareHashAndPassword err:", err)
		return
	}

	if dbUser.Status == "created" {
		c.AbortWithError(http.StatusNotFound, helpers.ErrorWithCode("user_needs_activation", "User needs to be activated via email", nil))
		return
	}

	//Read base64 private key
	encodedKey := []byte(config.GetString(c, "rsa_private"))
	returnToken(c, encodedKey, dbUser)
}

func (ac AuthController) TokenRenewal(c *gin.Context) {
	type UserInput struct {
		RefreshToken string `json:"refresh_token"`
	}
	var input UserInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		fmt.Println(err)
		return
	}
	encodedKey := []byte(config.GetString(c, "rsa_private"))
	refreshTokenClaims, err := helpers.ValidateJwtToken(input.RefreshToken, encodedKey, "refresh")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("refresh_token_invalid", "Refresh token invalid", err))
		return
	}

	dbUser, err := store.GetUserByID(c, fmt.Sprintf("%v", refreshTokenClaims["sub"]))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, helpers.ErrorWithCode("user_not_found", "User not found", err))
		return
	}

	returnToken(c, encodedKey, dbUser)
}
