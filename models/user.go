package models

import (
	"github.com/globalsign/mgo/bson"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"go-lumen/lumen-api/helpers"
	"golang.org/x/crypto/bcrypt"
)

// User type holds all required informations
type User struct {
	Id                 string `json:"id" bson:"_id,omitempty" valid:"-"`
	FirstName          string `json:"first_name" bson:"first_name"`
	LastName           string `json:"last_name" bson:"last_name"`
	Password           string `json:"password" bson:"password" valid:"required"`
	Email              string `json:"email" bson:"email" valid:"email,required"`
	Phone              string `json:"phone" bson:"phone"`
	Language           string `json:"language" bson:"language"`
	ActivationKey      string `json:"activation_key" bson:"activation_key"`
	ResetKey           string `json:"reset_key" bson:"reset_key"`
	Active             bool   `json:"active" bson:"active"`
	Admin              bool   `json:"admin" bson:"admin"`
	LastAccess         int64  `json:"last_access" bson:"last_access" valid:"-"`
	LastPasswordUpdate int64  `json:"last_password_update" bson:"last_password_update"`
}

// SanitizedUser allows to expose only few characteristics
type SanitizedUser struct {
	Id        string `json:"id" bson:"_id,omitempty"`
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	Email     string `json:"email" bson:"email"`
	Admin     bool   `json:"admin" bson:"admin"`
}

// Sanitize allows to create a lightweight user
func (user *User) Sanitize() SanitizedUser {
	return SanitizedUser{user.Id, user.FirstName, user.LastName, user.Email, user.Admin}
}

// BeforeCreate is here to check inputs and generating an encrypted password
func (user *User) BeforeCreate() error {
	user.Id = bson.NewObjectId().Hex()
	user.Active = false
	user.ActivationKey = helpers.RandomString(40)
	user.Email = strings.ToLower(user.Email)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "encryption_failed", "Failed to generate the crypted password", err)
	}
	user.Password = string(hashedPassword)

	_, err = govalidator.ValidateStruct(user)
	if err != nil {
		return helpers.NewError(http.StatusBadRequest, "input_not_valid", err.Error(), err)
	}

	return nil
}

// UsersCollection represents a specific MongoDB collection
const UsersCollection = "users"
