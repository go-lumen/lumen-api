package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-lumen/lumen-api/config"
	"github.com/go-lumen/lumen-api/helpers"
	"github.com/go-lumen/lumen-api/models"
	"github.com/go-lumen/lumen-api/store"
	"github.com/go-lumen/lumen-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// AuthController structure
type AuthController struct {
	BaseController
}

// NewAuthController instantiates the controller
func NewAuthController() AuthController {
	return AuthController{}
}

func (ac AuthController) returnToken(c *gin.Context, encodedKey []byte, dbUser *models.User) {
	ctx := store.AuthContext(c)

	accessToken, err := helpers.GenerateToken(encodedKey, dbUser.ID, "access", 4320) //3d in min
	if err != nil {
		utils.CheckErr(err)
		ac.AbortWithError(c, helpers.ErrorTokenGenAccess(err))
		return
	}
	refreshToken, err := helpers.GenerateToken(encodedKey, dbUser.ID, "refresh", 10080) //7d in min
	if err != nil {
		utils.CheckErr(err)
		ac.AbortWithError(c, helpers.ErrorTokenGenRefresh(err))
		return
	}

	// Get group : orga & role
	var group models.Group
	var organization models.Organization
	if ac.Error(c, ctx.Store.Find(ctx, utils.ParamID(dbUser.GroupID), &group), helpers.ErrorResourceNotFound) {
		return
	}
	if ac.Error(c, ctx.Store.Find(ctx, utils.ParamID(group.OrganizationID), &organization), helpers.ErrorResourceNotFound) {
		return
	}

	err = ctx.Store.Update(ctx, store.ID(dbUser.ID), &models.User{LastLogin: time.Now().Unix()},
		store.OnlyFields([]string{"last_login"}))
	if ac.ErrorInternal(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": accessToken, "refresh_token": refreshToken, store.RoleUser: dbUser.Sanitize(group.Role, organization.ID, organization.Name)})
}

// TokensGeneration to authenticate the user and generate a new token
// @Summary Returns a token from username and password
// @Produce json
// @Param userAuth body models.UserAuth true "Query Params"
// @Success 200 {object} models.User
// @Router /v1/auth [post]
// @Security ApiKeyAuth
// @Tags Authentication
func (ac AuthController) TokensGeneration(c *gin.Context) {
	ctx := store.AuthContext(c)
	userInput := models.User{}
	if err := c.Bind(&userInput); err != nil {
		ac.AbortWithError(c, helpers.ErrorInvalidInput(err))
		return
	}

	var dbUser models.User
	if ac.Error(c, ctx.Store.Find(ctx, bson.M{"email": userInput.Email}, &dbUser), helpers.ErrorUserNotExist) {
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(userInput.Password))
	if err != nil {
		ac.AbortWithError(c, helpers.ErrorUserWrongPassword(err))
		utils.Log(c, "info", "CompareHashAndPassword err:", err)
		return
	}

	if dbUser.Status == "created" {
		ac.AbortWithError(c, helpers.ErrorUserNotActivated(nil))
		return
	}

	//Read base64 private key
	encodedKey := []byte(config.GetString(c, "rsa_private"))
	ac.returnToken(c, encodedKey, &dbUser)
}

// TokenRenewal to renew a refresh token
func (ac AuthController) TokenRenewal(c *gin.Context) {
	ctx := store.AuthContext(c)

	type UserInput struct {
		RefreshToken string `json:"refresh_token"`
	}
	var input UserInput
	if err := c.ShouldBind(&input); err != nil {
		ac.AbortWithError(c, helpers.ErrorInvalidInput(err))
		return
	}
	encodedKey := []byte(config.GetString(c, "rsa_private"))
	refreshTokenClaims, err := helpers.ValidateJwtToken(input.RefreshToken, encodedKey, "refresh")
	if err != nil {
		ac.AbortWithError(c, helpers.ErrorTokenRefreshInvalid(err))
		return
	}

	var dbUser models.User
	if ac.Error(c, ctx.Store.Find(ctx, utils.ParamID(fmt.Sprintf("%v", refreshTokenClaims["sub"])), &dbUser), helpers.ErrorUserNotExist) {
		return
	}

	ac.returnToken(c, encodedKey, &dbUser)
}
