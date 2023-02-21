package models

import (
	"errors"
	"github.com/asaskevich/govalidator"
	mgobson "github.com/globalsign/mgo/bson"
	"github.com/go-lumen/lumen-api/helpers"
	"github.com/go-lumen/lumen-api/store"
	"github.com/go-lumen/lumen-api/utils"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
)

// Organization type holds all required information
type Organization struct {
	store.DefaultRoles `bson:"-,omitempty"`
	ID                 string `json:"id" bson:"id,omitempty" valid:"-"`
	Name               string `json:"name" bson:"name" valid:"-"`
	LogoURL            string `json:"logo_url,omitempty" bson:"logo_url,omitempty" valid:"-"`
	Siret              int64  `json:"siret,omitempty" bson:"siret,omitempty" valid:"-"`
	VATNumber          string `json:"vat_number,omitempty" bson:"vat_number,omitempty" valid:"-"`
	Tokens             int64  `json:"tokens,omitempty" bson:"tokens,omitempty" valid:"-"`
	Country            string `json:"country" bson:"country" valid:"-"`
	AppKey             string `json:"app_key,omitempty" bson:"app_key,omitempty" valid:"-"`
	Parent             string `json:"parent_id,omitempty" bson:"parent_id,omitempty" valid:"-"`
	DefaultGroupID     string `json:"default_group_id" bson:"default_group_id" valid:"-"`
}

// GetCollection returns mongodb collection
func (organization *Organization) GetCollection() string {
	return "organizations"
}

// OrgByAppKey returns an org filter
func OrgByAppKey(key string) bson.M { return bson.M{"app_key": key} }

// SanitizedOrganization type holds only essential informations
type SanitizedOrganization struct {
	ID      string `json:"id" bson:"_id,omitempty" valid:"-"`
	Name    string `json:"name" bson:"name" valid:"-"`
	LogoURL string `json:"logo_url" bson:"logo_url" valid:"-"`
	Parent  string `json:"parent_id" bson:"parent_id" valid:"-"`
}

// SanitizedOrganizationWithParent type holds only essential informations with nested parent
type SanitizedOrganizationWithParent struct {
	ID      string                `json:"id" bson:"_id,omitempty" valid:"-"`
	Name    string                `json:"name" bson:"name" valid:"-"`
	LogoURL string                `json:"logo_url" bson:"logo_url" valid:"-"`
	Parent  SanitizedOrganization `json:"parent_organization" bson:"parent_organization" valid:"-"`
}

// Sanitize allows to generate a SanitizedOrganization
func (organization *Organization) Sanitize() SanitizedOrganization {
	return SanitizedOrganization{organization.ID, organization.Name, organization.LogoURL, organization.Parent}
}

// SanitizeWithParent allows to generate a SanitizedOrganizationWithParent
func (organization *Organization) SanitizeWithParent(parentOrganization SanitizedOrganization) SanitizedOrganizationWithParent {
	return SanitizedOrganizationWithParent{organization.ID, organization.Name, organization.LogoURL, parentOrganization}
}

// FindOrganization is used to find an organization in a organizations list (for performance purposes, only 1 db request)
func FindOrganization(dbOrganizations []*Organization, organizationID string) (ret *Organization, err error) {
	for _, organization := range dbOrganizations {
		if organization.ID == organizationID {
			return organization, nil
		}
	}
	return nil, errors.New("Organization not found")
}

// BeforeCreate validates object struct
func (organization *Organization) BeforeCreate() error {
	organization.ID = mgobson.NewObjectId().Hex()
	organization.AppKey = helpers.RandomString(40)

	_, err := govalidator.ValidateStruct(organization)
	if err != nil {
		return helpers.NewError(http.StatusBadRequest, "input_not_valid", err.Error(), err)
	}
	return nil
}

// ApplyOptions set default find options
func (organization *Organization) ApplyOptions(o *options.FindOptions) {
	o.SetSort(bson.D{{Key: "name", Value: 1}})
}

// CreateOrganization checks if organization already exists, and if not, creates it
func CreateOrganization(c *store.Context, organization *Organization) error {
	err := organization.BeforeCreate()
	if err != nil {
		return err
	}

	var existingOrgs []*Organization
	err = c.Store.FindAll(c, bson.M{"name": organization.Name}, &existingOrgs)
	if err != nil {
		return err
	}

	if len(existingOrgs) > 0 {
		return helpers.NewError(http.StatusConflict, "organization_already_exists", "Organization already exists", err)
	}

	err = c.Store.Create(c, "organizations", organization)
	if err != nil {
		utils.Log(nil, "warn", err)
		return helpers.NewError(http.StatusInternalServerError, "organization_creation_failed", "Failed to insert the organization in the database", err)
	}

	return nil
}

// GetOrganization allows to retrieve a organization by its characteristics
func GetOrganization(c *store.Context, filter bson.M) (*Organization, error) {
	var organization Organization
	err := c.Store.Find(c, filter, &organization)
	if err != nil {
		return nil, helpers.NewError(http.StatusNotFound, "organization_not_found", "Organization not found", err)
	}

	return &organization, err
}

// IsOrganizationParent allows to know if an organization is a parent, and retrieve its parent if not
func IsOrganizationParent(c *store.Context, organizationID string) (bool, string, error) {
	var organization Organization
	err := c.Store.Find(c, bson.M{"_id": organizationID}, &organization)
	if err != nil {
		return false, "", helpers.NewError(http.StatusNotFound, "organization_not_found", "Organization not found", err)
	}

	if len(organization.Parent) > 1 {
		return false, organization.Parent, nil
	}

	return true, "", err
}

// IsOrganizationChildren allows to know if an organization is a children
func IsOrganizationChildren(c *store.Context, organizationID string, comparedOne string) (bool, string, error) {
	var organization Organization
	err := c.Store.Find(c, bson.M{"_id": organizationID}, &organization)
	if err != nil {
		return false, "", helpers.NewError(http.StatusNotFound, "organization_not_found", "Organization not found", err)
	}

	if organization.Parent == comparedOne {
		return true, organization.Parent, nil
	}

	return false, organization.Parent, err
}

// GetOrganizations allows to get all organizations
func GetOrganizations(c *store.Context, filter bson.M) ([]*Organization, error) {
	var list []*Organization

	err := c.Store.FindAll(c, filter, &list)
	if err != nil {
		logrus.Warnln("ErrorInternal on Finding all the documents", err)
	}

	return list, err
}

// UpdateOrganization allows to update one or more organization characteristics
func UpdateOrganization(c *store.Context, organizationID string, newOrganization *Organization) error {
	err := c.Store.Update(c, store.ID(organizationID), newOrganization, store.CreateIfNotExists(true))
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "organization_update_failed", "Failed to update the organization", err)
	}

	return nil
}

// ChangeParent allows to change an organization parent by their IDs
func ChangeParent(c *store.Context, organizationID string, newParent string) error {
	err := c.Store.Update(c, store.ID(organizationID), &Organization{Parent: newParent},
		store.CreateIfNotExists(true),
		store.OnlyFields([]string{"parent_id"}))
	if err != nil {
		return helpers.NewError(http.StatusInternalServerError, "organization_update_failed", "Failed to update the organization", err)
	}

	return nil
}
