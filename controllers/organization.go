package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-lumen/lumen-api/helpers"
	"github.com/go-lumen/lumen-api/models"
	"github.com/go-lumen/lumen-api/store"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

// OrganizationController holds all controller functions related to the organization entity
type OrganizationController struct {
	BaseController
}

// NewOrganizationController instantiates the controller
func NewOrganizationController() OrganizationController {
	return OrganizationController{}
}

// CreateOrganization to create a new organization
func (oc OrganizationController) CreateOrganization(c *gin.Context) {
	ctx := store.AuthContext(c)
	organization := &models.Organization{}

	if oc.BindJSONError(c, organization) {
		return
	}

	if _, group, ok := oc.LoggedUser(c); ok {
		switch group.GetRole() {
		case store.RoleGod:
			if oc.ErrorInternal(c, models.CreateOrganization(ctx, organization)) {
				return
			}
			c.JSON(http.StatusCreated, organization)
		case store.RoleAdmin, store.RoleUser, store.RoleCustomer:
			oc.AbortWithError(c, helpers.ErrorUserUnauthorized)
		}
	}
}

// GetOrganizations to get all organizations
func (oc OrganizationController) GetOrganizations(c *gin.Context) {
	ctx := store.AuthContext(c)
	dbOrganizations, err := models.GetOrganizations(ctx, bson.M{})
	if oc.ErrorInternal(c, err) {
		return
	}

	if _, group, ok := oc.LoggedUser(c); ok {
		switch group.GetRole() {
		case store.RoleGod: // Get all
			c.JSON(http.StatusOK, dbOrganizations)
		case store.RoleAdmin, store.RoleUser: // Get all from devices with same group ID (and organization for admin)
			userOrganization, err := models.GetOrganization(ctx, bson.M{"_id": group.GetOrgID()})
			if oc.Error(c, err, helpers.ErrorResourceNotFound) {
				return
			}
			c.JSON(http.StatusOK, []*models.Organization{userOrganization})
		case store.RoleCustomer:
			oc.AbortWithError(c, helpers.ErrorUserUnauthorized)
		}
	}
}

// GetOrganization allows to get a specific Organization
func (oc OrganizationController) GetOrganization(c *gin.Context) {
	ctx := store.AuthContext(c)
	organization, err := models.GetOrganization(ctx, oc.ParamID(c))
	if oc.Error(c, err, helpers.ErrorResourceNotFound) {
		return
	}

	if _, group, ok := oc.LoggedUser(c); ok {
		switch group.GetRole() {
		case store.RoleGod:
			c.JSON(http.StatusOK, organization)
		case store.RoleAdmin, store.RoleUser:
			if organization.ID == group.GetOrgID() {
				c.JSON(http.StatusOK, organization)
				return
			}
			oc.AbortWithError(c, helpers.ErrorUserUnauthorized)
		case store.RoleCustomer:
			oc.AbortWithError(c, helpers.ErrorUserUnauthorized)
		}
	}
}

// GetOrganizationGroups allows to get groups from an organization
func (oc OrganizationController) GetOrganizationGroups(c *gin.Context) {
	ctx := store.AuthContext(c)
	organization, err := models.GetOrganization(ctx, oc.ParamID(c))
	if oc.Error(c, err, helpers.ErrorResourceNotFound) {
		return
	}

	groups, err := models.GetGroups(ctx, bson.M{"organization_id": c.Param("id")})
	if oc.Error(c, err, helpers.ErrorResourceNotFound) {
		return
	}

	if _, group, ok := oc.LoggedUser(c); ok {
		switch group.GetRole() {
		case store.RoleGod:
			c.JSON(http.StatusOK, groups)
		case store.RoleAdmin, store.RoleUser:
			if organization.ID == group.GetOrgID() {
				c.JSON(http.StatusOK, groups)
				return
			}
			oc.AbortWithError(c, helpers.ErrorUserUnauthorized)
		case store.RoleCustomer:
			oc.AbortWithError(c, helpers.ErrorUserUnauthorized)
		}
	}
}

// GetUserOrganization to get organization from user stored in context
func (oc OrganizationController) GetUserOrganization(c *gin.Context) {
	ctx := store.AuthContext(c)
	if _, group, ok := oc.LoggedUser(c); ok {
		userOrganization, err := models.GetOrganization(ctx, bson.M{"_id": group.GetOrgID()})
		if oc.Error(c, err, helpers.ErrorResourceNotFound) {
			return
		}

		switch group.GetRole() {
		case store.RoleGod: // Get all
			c.JSON(http.StatusOK, userOrganization)
		case store.RoleAdmin, store.RoleUser: // Get all from devices with same group ID (and organization for admin)
			c.JSON(http.StatusOK, userOrganization.Sanitize())
		case store.RoleCustomer:
			c.JSON(http.StatusOK, userOrganization.Sanitize())
		}
	}
}

// GetOrganizationByAppKey to get organization from AppKey
func (oc OrganizationController) GetOrganizationByAppKey(c *gin.Context) {
	ctx := store.AuthContext(c)
	var org models.Organization
	if oc.Error(c, ctx.Store.Find(ctx, models.OrgByAppKey(c.Param("id")), &org), helpers.ErrorResourceNotFound) {
		return
	}

	c.JSON(http.StatusOK, org)
}
