package models

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/go-lumen/lumen-api/helpers"
	"github.com/go-lumen/lumen-api/store"
	"github.com/go-lumen/lumen-api/utils"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

// UserAuth type holds required information for authentication
type UserAuth struct {
	Email    string `json:"email" bson:"email" valid:"email,required"`
	Password string `json:"password" bson:"password" valid:"required"`
}

// User type holds all required information
type User struct {
	store.DefaultRoles `bson:"-,omitempty"`
	ID                 string `json:"id" bson:"_id,omitempty" valid:"-"`
	FirstName          string `json:"first_name" bson:"first_name" valid:"-"`
	LastName           string `json:"last_name" bson:"last_name" valid:"-"`
	Password           string `json:"password" bson:"password" valid:"required"`
	Email              string `json:"email" bson:"email" valid:"email,required"`
	Address            string `json:"address,omitempty" bson:"address,omitempty" valid:"-"`
	Status             string `json:"status" bson:"status" valid:"-"`
	Phone              string `json:"phone" bson:"phone" valid:"-"`
	Language           string `json:"language,omitempty" bson:"language,omitempty" valid:"-"`
	Key                string `json:"key" bson:"key" valid:"-"`
	LastLogin          int64  `json:"last_login" bson:"last_login" valid:"-"`
	LastModification   int64  `json:"last_modification" bson:"last_modification" valid:"-"`
	GroupID            string `json:"group_id" bson:"group_id" valid:"-"`
}

// GetID returns ID
func (user *User) GetID() string {
	return user.ID
}

// GetGroupID returns organization ID
func (user *User) GetGroupID() string {
	return user.GroupID
}

// GetCollection returns mongodb collection
func (user *User) GetCollection() string {
	return UsersCollection
}

// UserDetails type holds user with details
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

// Sanitize allows to create a lightweight user
func (user *User) Sanitize(role string, organizationID string, organizationName string) SanitizedUser {
	return SanitizedUser{user.ID, user.FirstName, user.LastName, user.Email, user.Status, user.GroupID, role, organizationID, organizationName}
}

// Detail to detail a user
func (user *User) Detail(role string, organization string) UserDetails {
	return UserDetails{user.ID, user.FirstName, user.LastName, user.Email, user.Address, user.Status, user.Phone, user.Language, user.GroupID, role, organization}
}

// FindUser is used to find a user in a users list (for performance purposes, only 1 db request)
func FindUser(dbUsers []*User, userID string) (ret *User, err error) {
	for _, user := range dbUsers {
		if userID == user.ID {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

// BeforeCreate validates object struct
func (user *User) BeforeCreate() error {
	//TODO: Check if we can remove this old dependency
	//user.ID = mgobson.NewObjectId().Hex()
	user.Key = helpers.RandomString(40)
	user.Email = strings.ToLower(user.Email)
	user.LastModification = time.Now().Unix()
	user.Status = "created"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "encryption_failed", "Failed to generate the encrypted password", err)
	}
	user.Password = string(hashedPassword)

	if _, err = govalidator.ValidateStruct(user); err != nil {
		return helpers.NewError(http.StatusBadRequest, "input_not_valid", err.Error(), err)
	}
	return nil
}

// UsersCollection represents a specific MongoDB collection
const UsersCollection = "users"

// UsersTableName represents a SQL table
const UsersTableName = "users_user"

// CreateUser checks if user already exists, and if not, creates it
func CreateUser(c *store.Context, user *User) error {
	var existingUsers []*User
	if err := c.Store.FindAll(c, bson.M{"email": user.Email}, &existingUsers); err != nil {
		return err
	}

	if len(existingUsers) > 0 {
		return helpers.NewError(http.StatusConflict, "user_already_exists", "User already exists", nil)
	}

	err := c.Store.Create(c, user)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "user_creation_failed", "Failed to insert the user in the database", err)
	}

	return nil
}

// GetUser allows to retrieve a user by its characteristics
func GetUser(c *store.Context, filter bson.M) (*User, error) {
	var user User
	err := c.Store.Find(c, filter, &user)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "user_not_found", "User not found", err)
	}

	return &user, err
}

// UserExists allows checking if a user exists by its email
func UserExists(c *store.Context, email string) (bool, *User, error) {
	var user User
	err := c.Store.Find(c, bson.M{"email": email}, &user)
	if err != nil {
		return false, nil, helpers.NewError(http.StatusNotFound, "user_not_found", "User not found", err)
	}

	return true, &user, err
}

// GetUsers allows to get all users
func GetUsers(c *store.Context, filter bson.M) ([]*User, error) {
	var list []*User

	err := c.Store.FindAll(c, filter, &list, store.WithSort("last_name", store.SortAscending))
	if err != nil {
		logrus.Warnln("ErrorInternal on Finding all the documents", err)
	}

	return list, err
}

// UpdateUser allows to update one or more user characteristics
func UpdateUser(c *store.Context, userID string, newUser *User) error {
	err := c.Store.Update(c, store.ID(userID), newUser, store.CreateIfNotExists(true))
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "user_update_failed", "Failed to update the user", err)
	}

	return nil
}

// ActivateUser allows to activate a user by its id
func ActivateUser(c *store.Context, activationKey string, id string) error {
	err := c.Store.Update(c, bson.M{"$and": []bson.M{utils.ParamID(id), {"key": activationKey}}}, &User{Status: "activated"}, store.OnlyFields([]string{"status"}))

	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "user_activation_failed", "Couldn't find the user to activate", err)
	}

	return nil
}

// DeleteUser allows to delete a user by its id
func DeleteUser(c *store.Context, userID string) error {
	err := c.Store.Delete(c, userID, &User{})
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "user_delete_failed", "Failed to delete the user", err)
	}

	return nil
}
