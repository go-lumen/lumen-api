package postgresql

import (
	"github.com/adrien3d/stokelp-poc/helpers/params"
	"github.com/adrien3d/stokelp-poc/models"
)

// CreateGroup checks if group already exists, and if not, creates it
func (db *postgresql) CreateGroup(group *models.Group) error {
	return nil
}

// GetGroupByID allows to retrieve a group by its id
func (db *postgresql) GetGroupByID(id string) (*models.Group, error) {
	return nil, nil
}

// GetGroup allows to retrieve a group by its characteristics
func (db *postgresql) GetGroup(params params.M) (*models.Group, error) {
	return nil, nil
}

// UpdateGroup allows to update one or more group characteristics
func (db *postgresql) UpdateGroup(groupID string, params params.M) error {
	return nil
}

// DeleteGroup allows to delete a group by its id
func (db *postgresql) DeleteGroup(groupID string) error {
	return nil
}

// ChangeGroupOrganization allows to change the organization of a group by its id
func (db *postgresql) ChangeGroupOrganization(groupID string, organizationID string) error {
	return nil
}

// GetGroups allows to get all groups
func (db *postgresql) GetGroups() ([]*models.Group, error) {
	return nil, nil
}

// CountGroups allows to count all groups
func (db *postgresql) CountGroups() (int, error) {
	return 0, nil
}

// GroupExists allows to know if a group exists through his mail
func (db *postgresql) GroupExists(groupEmail string) (bool, error) {
	return false, nil
}
