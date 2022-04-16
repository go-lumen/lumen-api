package store

// RoleChecker defines a standard interface for role checking
type RoleChecker interface {
	// CanBeCreated allows to determine if a user can create a device
	CanBeCreated(user User, group Group) bool

	// CanBeRead allows to determine if a user can read a device
	CanBeRead(user User, group Group) bool

	// CanBeUpdated allows to determine if a user can update a device
	CanBeUpdated(user User, group Group) bool

	// CanBeDeleted allows to determine if a user can delete a device
	CanBeDeleted(User, Group) bool
}

// DefaultRoles implements default role for basic models
type DefaultRoles struct{}

// CanBeCreated allows to determine if a user can create a device
func (dr *DefaultRoles) CanBeCreated(user User, group Group) bool {
	switch group.GetRole() {
	case RoleGod, RoleAdmin:
		return true
	default:
		return false
	}
}

// CanBeRead allows to determine if a user can read a device
func (dr *DefaultRoles) CanBeRead(user User, group Group) bool {
	switch group.GetRole() {
	case RoleGod, RoleAdmin, RoleUser, RoleCustomer:
		return true
	default:
		return false
	}
}

// CanBeUpdated allows to determine if a user can update a device
func (dr *DefaultRoles) CanBeUpdated(user User, group Group) bool {
	return dr.CanBeCreated(user, group)
}

// CanBeDeleted allows to determine if a user can delete a device
func (dr *DefaultRoles) CanBeDeleted(user User, group Group) bool {
	return dr.CanBeCreated(user, group)
}
