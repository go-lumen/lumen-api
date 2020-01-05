package models

import (
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/go-lumen/lumen-api/helpers"
	"golang.org/x/crypto/bcrypt"
)

// User type holds all required informations
type User struct {
	Id                 string `json:"id" bson:"_id,omitempty" valid:"-"`
	FirstName          string `json:"first_name" bson:"first_name" valid:"-"`
	LastName           string `json:"last_name" bson:"last_name" valid:"-"`
	Password           string `json:"password" bson:"password" valid:"required"`
	Email              string `json:"email" bson:"email" valid:"email,required"`
	Role               string `json:"role" bson:"role" valid:"required"`
	Address            string `json:"address,omitempty" bson:"address,omitempty" valid:"-"`
	Status             string `json:"status" bson:"status" valid:"-"`
	Phone              string `json:"phone" bson:"phone" valid:"-"`
	Language           string `json:"language,omitempty" bson:"language,omitempty" valid:"-"`
	ActivationKey      string `json:"activation_key" bson:"activation_key" valid:"-"`
	ResetKey           string `json:"reset_key" bson:"reset_key" valid:"-"`
	LastModification   int64  `json:"last_modification" bson:"last_modification" valid:"-"`
	LastPasswordUpdate int64  `json:"last_password_update" bson:"last_password_update" valid:"-"`
	GroupId            string `json:"group_id" bson:"group_id" valid:"-"`
}

type UserDetails struct {
	Id        string `json:"id" bson:"_id,omitempty" valid:"-"`
	FirstName string `json:"first_name" bson:"first_name" valid:"-"`
	LastName  string `json:"last_name" bson:"last_name" valid:"-"`
	Email     string `json:"email" bson:"email" valid:"email,required"`
	Role      string `json:"role" bson:"role" valid:"required"`
	Address   string `json:"address" bson:"address" valid:"-"`
	Status    string `json:"status" bson:"status" valid:"-"`
	Phone     string `json:"phone" bson:"phone" valid:"-"`
	Language  string `json:"language" bson:"language" valid:"-"`
}

// SanitizedUser allows to expose only few characteristics
type SanitizedUser struct {
	Id        string `json:"id" bson:"_id,omitempty" valid:"-"`
	FirstName string `json:"first_name" bson:"first_name" valid:"-"`
	LastName  string `json:"last_name" bson:"last_name" valid:"-"`
	Email     string `json:"email" bson:"email" valid:"-"`
	Status    string `json:"status" bson:"status" valid:"-"`
	Role      string `json:"role" bson:"role" valid:"required"`
}

// Sanitize allows to create a lightweight user
func (user *User) Sanitize(role string, organization string) SanitizedUser {
	return SanitizedUser{user.Id, user.FirstName, user.LastName, user.Email, user.Status, role}
}

func (user *User) Detail(role string, organization string) UserDetails {
	return UserDetails{user.Id, user.FirstName, user.LastName, user.Email, role, user.Address, user.Status, user.Phone, user.Language}
}

// BeforeCreate is here to check inputs and generating an encrypted password
func (user *User) BeforeCreate() error {
	//user.Id = bson.NewObjectId().Hex()
	user.ActivationKey = helpers.RandomString(40)
	user.Email = strings.ToLower(user.Email)
	user.Status = "created"

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
