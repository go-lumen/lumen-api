package postgresql

import (
	"github.com/go-lumen/lumen-api/helpers/params"
	"github.com/go-lumen/lumen-api/models"
)

// CreateGroup checks if group already exists, and if not, creates it
func (db *PSQL) CreateGroup(group *models.Group) error {
	return nil
}

// GetGroupByID allows to retrieve a group by its id
func (db *PSQL) GetGroupByID(id string) (*models.Group, error) {
	return nil, nil
}

// GetGroup allows to retrieve a group by its characteristics
func (db *PSQL) GetGroup(params params.M) (*models.Group, error) {
	return nil, nil
}

// UpdateGroup allows to update one or more group characteristics
func (db *PSQL) UpdateGroup(groupID string, params params.M) error {
	return nil
}

// DeleteGroup allows to delete a group by its id
func (db *PSQL) DeleteGroup(groupID string) error {
	return nil
}

// ChangeGroupOrganization allows to change the organization of a group by its id
func (db *PSQL) ChangeGroupOrganization(groupID string, organizationID string) error {
	return nil
}

// GetGroups allows to get all groups
func (db *PSQL) GetGroups() ([]*models.Group, error) {
	return nil, nil
}

// CountGroups allows to count all groups
func (db *PSQL) CountGroups() (int, error) {
	return 0, nil
}

// GroupExists allows to know if a group exists through his mail
func (db *PSQL) GroupExists(groupEmail string) (bool, error) {
	return false, nil
}
