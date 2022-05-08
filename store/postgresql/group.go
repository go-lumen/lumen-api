package postgresql

import (
	"github.com/go-lumen/lumen-api/helpers/params"
	"github.com/go-lumen/lumen-api/models"
)

// CreateGroup checks if group already exists, and if not, creates it
func (db *Postgresql) CreateGroup(group *models.Group) error {
	return nil
}

// GetGroupByID allows to retrieve a group by its id
func (db *Postgresql) GetGroupByID(id string) (*models.Group, error) {
	return nil, nil
}

// GetGroup allows to retrieve a group by its characteristics
func (db *Postgresql) GetGroup(params params.M) (*models.Group, error) {
	return nil, nil
}

// UpdateGroup allows to update one or more group characteristics
func (db *Postgresql) UpdateGroup(groupID string, params params.M) error {
	return nil
}

// DeleteGroup allows to delete a group by its id
func (db *Postgresql) DeleteGroup(groupID string) error {
	return nil
}

// ChangeGroupOrganization allows to change the organization of a group by its id
func (db *Postgresql) ChangeGroupOrganization(groupID string, organizationID string) error {
	return nil
}

// GetGroups allows to get all groups
func (db *Postgresql) GetGroups() ([]*models.Group, error) {
	return nil, nil
}

// CountGroups allows to count all groups
func (db *Postgresql) CountGroups() (int, error) {
	return 0, nil
}

// GroupExists allows knowing if a group exists through his mail
func (db *Postgresql) GroupExists(groupEmail string) (bool, error) {
	return false, nil
}
