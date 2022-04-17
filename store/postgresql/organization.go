package postgresql

import (
	"github.com/go-lumen/lumen-api/helpers/params"
	"github.com/go-lumen/lumen-api/models"
)

// CreateOrganization checks if organization already exists, and if not, creates it
func (db *postgresql) CreateOrganization(organization *models.Organization) error {
	return nil
}

// GetOrganizationByID allows to retrieve a organization by its id
func (db *postgresql) GetOrganizationByID(id string) (*models.Organization, error) {
	return nil, nil
}

// GetOrganization allows to retrieve a organization by its characteristics
func (db *postgresql) GetOrganization(params params.M) (*models.Organization, error) {
	return nil, nil
}

// UpdateOrganization allows to update one or more organization characteristics
func (db *postgresql) UpdateOrganization(organizationID string, params params.M) error {
	return nil
}

// DeleteOrganization allows to delete a organization by its id
func (db *postgresql) DeleteOrganization(organizationID string) error {
	return nil
}

// IsOrganizationParent allows to know if an organization is a parent, and retrieve its parent if not
func (db *postgresql) IsOrganizationParent(organizationID string) (bool, string, error) {
	return false, "", nil
}

// ChangeParent allows to change an organization parent by its id
func (db *postgresql) ChangeParent(organizationID string, parentID string) error {
	return nil
}

// GetOrganizations allows to get all organizations
func (db *postgresql) GetOrganizations() ([]*models.Organization, error) {
	return nil, nil
}

// CountOrganizations allows to count all organizations
func (db *postgresql) CountOrganizations() (int, error) {
	return 0, nil
}

// OrganizationExists allows to know if a organization exists through his mail
func (db *postgresql) OrganizationExists(organizationEmail string) (bool, error) {
	return true, nil
}
