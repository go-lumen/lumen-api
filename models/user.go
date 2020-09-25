package models

import (
	"github.com/asaskevich/govalidator"
	"github.com/globalsign/mgo/bson"
	"net/http"
	"strings"
	"time"

	"github.com/go-lumen/lumen-api/helpers"
	"golang.org/x/crypto/bcrypt"
)

// User type holds all required informations
type User struct {
	ID               string `json:"id" bson:"_id,omitempty" valid:"-"`
	FirstName        string `json:"first_name" bson:"first_name" valid:"-"`
	LastName         string `json:"last_name" bson:"last_name" valid:"-"`
	Password         string `json:"password" bson:"password" valid:"required"`
	Email            string `json:"email" bson:"email" valid:"email,required"`
	Role             string `json:"role" bson:"role" valid:"role,required"`
	Address          string `json:"address,omitempty" bson:"address,omitempty" valid:"-"`
	Status           string `json:"status" bson:"status" valid:"-"`
	Phone            string `json:"phone" bson:"phone" valid:"-"`
	Language         string `json:"language,omitempty" bson:"language,omitempty" valid:"-"`
	Key              string `json:"key" bson:"key" valid:"-"`
	LastLogin        int64  `json:"last_login" bson:"last_login" valid:"-"`
	LastModification int64  `json:"last_modification" bson:"last_modification" valid:"-"`
	GroupID          string `json:"group_id" bson:"group_id" valid:"-"`
}

type UserDetails struct {
	ID           string `json:"id" bson:"_id,omitempty" valid:"-"`
	FirstName    string `json:"first_name" bson:"first_name" valid:"-"`
	LastName     string `json:"last_name" bson:"last_name" valid:"-"`
	Email        string `json:"email" bson:"email" valid:"email,required"`
	Address      string `json:"address" bson:"address" valid:"-"`
	Status       string `json:"status" bson:"status" valid:"-"`
	Phone        string `json:"phone" bson:"phone" valid:"-"`
	Language     string `json:"language" bson:"language" valid:"-"`
	GroupID      string `json:"group_id" bson:"group_id" valid:"-"`
	Role         string `json:"role" bson:"role" valid:"-"`
	Organization string `json:"organization" bson:"organization" valid:"-"`
}

// SanitizedUser allows to expose only few characteristics
type SanitizedUser struct {
	ID               string `json:"id" bson:"_id,omitempty" valid:"-"`
	FirstName        string `json:"first_name" bson:"first_name" valid:"-"`
	LastName         string `json:"last_name" bson:"last_name" valid:"-"`
	Email            string `json:"email" bson:"email" valid:"-"`
	Status           string `json:"status" bson:"status" valid:"-"`
	GroupID          string `json:"group_id" bson:"group_id" valid:"-"`
	Role             string `json:"role" bson:"role" valid:"-"`
	OrganizationID   string `json:"organization_id" bson:"organization_id" valid:"-"`
	OrganizationName string `json:"organization_name" bson:"organization_name" valid:"-"`
}

type APIKeys struct {
	HereKey     string `json:"here_key"`
	TomtomKey   string `json:"tomtom_key"`
	MichelinKey string `json:"michelin_key"`
}

// Sanitize allows to create a lightweight user
func (user *User) Sanitize(role string, organizationID string, organizationName string) SanitizedUser {
	return SanitizedUser{user.ID, user.FirstName, user.LastName, user.Email, user.Status, user.GroupID, role, organizationID, organizationName}
}

func (user *User) Detail(role string, organization string) UserDetails {
	return UserDetails{user.ID, user.FirstName, user.LastName, user.Email, user.Address, user.Status, user.Phone, user.Language, user.GroupID, role, organization}
}

// BeforeCreate is here to check inputs and generating an encrypted password
func (user *User) BeforeCreate() error {
	user.ID = bson.NewObjectId().Hex()
	user.Key = helpers.RandomString(40)
	user.Email = strings.ToLower(user.Email)
	user.LastModification = time.Now().Unix()
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
