package models

import (
	"errors"
	"github.com/asaskevich/govalidator"
	mgobson "github.com/globalsign/mgo/bson"
	"github.com/go-lumen/lumen-api/helpers"
	"github.com/go-lumen/lumen-api/store"
	"github.com/sahilm/fuzzy"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

// Group type holds all required informations
type Group struct {
	//store.DefaultRoles `bson:"-,omitempty"`
	ID             string `json:"id" bson:"_id,omitempty" valid:"-"`
	Name           string `json:"name" bson:"name" valid:"-"`
	Role           string `json:"role" bson:"role" valid:"-"`
	OrganizationID string `json:"organization_id" bson:"organization_id" valid:"-"`
}

// GetID returns ID
func (group *Group) GetID() string {
	return group.ID
}

// GetRole returns role
func (group *Group) GetRole() store.UserRole {
	return group.Role
}

// GetOrgID returns organization ID
func (group *Group) GetOrgID() store.UserRole {
	return group.OrganizationID
}

// GetCollection returns mongodb collection
func (group *Group) GetCollection() string {
	return GroupsCollection
}

// BeforeCreate validates object struct
func (group *Group) BeforeCreate() error {
	group.ID = mgobson.NewObjectId().Hex()

	_, err := govalidator.ValidateStruct(group)
	if err != nil {
		return helpers.NewError(http.StatusBadRequest, "input_not_valid", err.Error(), err)
	}
	return nil
}

// FindGroup is used to find a group in a groups list (for performance purposes, only 1 db request)
func FindGroup(dbGroups []*Group, groupID string) (ret *Group, err error) {
	for _, group := range dbGroups {
		if group.ID == groupID {
			return group, nil
		}
	}
	return nil, errors.New("group not found")
}

// FindGroupByFuzzyName is used to find a group in a groups list by fuzzy name matching (for performance purposes, only 1 db request)
func FindGroupByFuzzyName(dbGroups []*Group, name string) (ret *Group, err error) {
	var groupNames []string
	for _, group := range dbGroups {
		groupNames = append(groupNames, group.Name)
	}
	matches := fuzzy.Find(name, groupNames)

	for _, group := range dbGroups {
		if matches.Len() > 0 && group.Name == matches[0].Str {
			return group, nil
		}
	}
	return nil, errors.New("group not found")
}

// GroupsCollection represents a specific MongoDB collection
const GroupsCollection = "groups"

// CreateGroup checks if group already exists, and if not, creates it
func CreateGroup(c *store.Context, group *Group) error {

	err := group.BeforeCreate()
	if err != nil {
		return err
	}

	var existingGroups []*Group
	err = c.Store.FindAll(c, bson.M{"name": group.Name}, &existingGroups)
	if err != nil {
		return err
	}

	if len(existingGroups) > 0 {
		return helpers.NewError(http.StatusConflict, "group_already_exists", "Group already exists", err)
	}

	err = c.Store.Create(c, "groups", group)
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "group_creation_failed", "Failed to insert the group in the database", err)
	}

	return nil
}

// GetGroup allows to retrieve a group by its characteristics
func GetGroup(c *store.Context, filter bson.M) (*Group, error) {
	var group Group
	err := c.Store.Find(c, filter, &group)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "group_not_found", "Group not found", err)
	}

	return &group, err
}

// GetGroups allows to get all groups
func GetGroups(c *store.Context, filter bson.M) ([]*Group, error) {
	var list []*Group
	err := c.Store.FindAll(c, filter, &list)
	if err != nil {
		logrus.Warnln("ErrorInternal on Finding all the documents", err)
	}

	return list, err
}

// ChangeGroupOrganization allows to change the organization of a group by its id
func ChangeGroupOrganization(c *store.Context, groupID string, organizationID string) error {
	err := c.Store.Update(c, store.ID(groupID), &Group{OrganizationID: organizationID},
		store.OnlyFields([]string{"organization_id"}),
		store.CreateIfNotExists(true))

	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "group_group_change_failed", "Couldn't find the group to change group", err)
	}

	return nil
}
