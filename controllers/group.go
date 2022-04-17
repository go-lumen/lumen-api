package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-lumen/lumen-api/helpers"
	"github.com/go-lumen/lumen-api/models"
	"github.com/go-lumen/lumen-api/store"
	"github.com/go-lumen/lumen-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"sort"
)

// GroupController holds all controller functions related to the group entity
type GroupController struct {
	BaseController
}

// NewGroupController instantiates the controller
func NewGroupController() GroupController {
	return GroupController{}
}

// CreateGroup to create a new group
func (gc GroupController) CreateGroup(c *gin.Context) {
	ctx := store.AuthContext(c)
	group := &models.Group{}

	if err := c.BindJSON(group); err != nil {
		gc.AbortWithError(c, helpers.ErrorInvalidInput(err))
		return
	}

	if len(group.OrganizationID) < 20 {
		organization, err := models.GetOrganization(ctx, bson.M{"index": group.OrganizationID})
		if err == nil {
			group.OrganizationID = organization.ID
		}
	}

	if _, userGroup, ok := gc.LoggedUser(c); ok {
		switch userGroup.GetRole() {
		case store.RoleGod:
			if err := models.CreateGroup(ctx, group); err != nil {
				gc.AbortWithError(c, helpers.ErrorInternal(err))
				return
			}
			c.JSON(http.StatusCreated, group)
		case store.RoleAdmin, store.RoleUser, store.RoleCustomer:
			gc.AbortWithError(c, helpers.ErrorUserUnauthorized)
		}
	}
}

// GetGroupsIndex to get groups index
func (gc GroupController) GetGroupsIndex(c *gin.Context) {
	ctx := store.AuthContext(c)
	dbGroups, err := models.GetGroups(ctx, bson.M{})
	if err != nil {
		gc.AbortWithError(c, helpers.ErrorResourceNotFound(err))
	}

	groupIndex := int64(0)
	for _, group := range dbGroups {
		if group.Index > groupIndex {
			groupIndex = group.Index
		}
	}

	if _, group, ok := gc.LoggedUser(c); ok {
		userGroup, err := models.GetGroup(ctx, utils.ParamID(group.GetID()))
		if err != nil {
			gc.AbortWithError(c, helpers.ErrorResourceNotFound(err))
			return
		}

		switch userGroup.Role {
		case store.RoleGod, store.RoleAdmin:
			c.JSON(http.StatusOK, groupIndex)
		case store.RoleUser, store.RoleCustomer:
			gc.AbortWithError(c, helpers.ErrorUserUnauthorized)
		}
	}
}

// GetGroups to get all groups
func (gc GroupController) GetGroups(c *gin.Context) {
	ctx := store.AuthContext(c)
	dbGroups, err := models.GetGroups(ctx, bson.M{})
	if gc.ErrorInternal(c, err) {
		return
	}

	dbOrgas, err := models.GetOrganizations(ctx, bson.M{})
	if gc.ErrorInternal(c, err) {
		return
	}

	if _, userGroup, ok := gc.LoggedUser(c); ok {
		switch userGroup.GetRole() {
		case store.RoleGod: // Get all
			for _, group := range dbGroups {
				orga, err := models.FindOrganization(dbOrgas, group.OrganizationID)
				if err == nil {
					group.OrganizationID = orga.Name
				}
			}
			sort.Slice(dbGroups, func(i, j int) bool { return dbGroups[i].Name < dbGroups[j].Name })
			c.JSON(http.StatusOK, dbGroups)
		case store.RoleAdmin, store.RoleUser: // Get all from devices with same group ID (and group for admin)
			var retGroups []*models.Group
			for _, group := range dbGroups {
				if group.OrganizationID == userGroup.GetOrgID() {
					orga, err := models.FindOrganization(dbOrgas, group.OrganizationID)
					if err == nil {
						group.OrganizationID = orga.Name
					}
					retGroups = append(retGroups, group)
				}
			}
			sort.Slice(retGroups, func(i, j int) bool { return retGroups[i].Name < retGroups[j].Name })
			c.JSON(http.StatusOK, retGroups)
		case store.RoleCustomer:
			gc.AbortWithError(c, helpers.ErrorUserUnauthorized)
		}
	}
}

// GetUserGroup to get group from user stored in context
func (gc GroupController) GetUserGroup(c *gin.Context) {
	if _, group, ok := gc.LoggedUser(c); ok {
		c.JSON(http.StatusOK, group)
	}
}

// GetGroup allows to get a specific Group
func (gc GroupController) GetGroup(c *gin.Context) {
	ctx := store.AuthContext(c)
	group, err := models.GetGroup(ctx, gc.ParamID(c))
	if gc.Error(c, err, helpers.ErrorResourceNotFound) {
		return
	}

	if _, userGroup, ok := gc.LoggedUser(c); ok {
		switch userGroup.GetRole() {
		case store.RoleGod:
			c.JSON(http.StatusOK, group)
		case store.RoleAdmin, store.RoleUser:
			if group.OrganizationID == userGroup.GetOrgID() {
				c.JSON(http.StatusOK, group)
				return
			}
			gc.AbortWithError(c, helpers.ErrorUserUnauthorized)
		case store.RoleCustomer:
			gc.AbortWithError(c, helpers.ErrorUserUnauthorized)
		}
	}
}
