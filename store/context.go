package store

import (
	"go-lumen/lumen-api/models"
	"golang.org/x/net/context"
)

const (
	// CurrentKey for user
	CurrentKey = "currentUser"
	// StoreKey for storing
	StoreKey = "store"
)

// Setter interface
type Setter interface {
	Set(string, interface{})
}

// Current allows to retrieve user from context
func Current(c context.Context) *models.User {
	return c.Value(CurrentKey).(*models.User)
}

// ToContext allows to set a value in store
func ToContext(c Setter, store Store) {
	c.Set(StoreKey, store)
}

// FromContext allows to get store from context
func FromContext(c context.Context) Store {
	return c.Value(StoreKey).(Store)
}
