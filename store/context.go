package store

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

const (
	// CurrentUserKey for user
	CurrentUserKey = "currentUser"
	// CurrentUserGroupKey for user group
	CurrentUserGroupKey = "currentUserGroup"
	// StoreKey for storing
	StoreKey = "store"
	// AppKey for external access
	AppKey = "appKey"
)

// Setter interface
type Setter interface {
	Set(string, interface{})
}

// Current allows to retrieve user from context
func Current(c context.Context) *User {
	return c.Value(CurrentUserKey).(*User)
}

// ToContext allows to set a value in store
func ToContext(c Setter, store Store) {
	c.Set(StoreKey, store)
}

// FromContext allows to get store from context
func FromContext(c context.Context) Store {
	return c.Value(StoreKey).(Store)
}

// UserRole represents an user role
type UserRole = string

const (
	// RoleGod have access to everything
	RoleGod UserRole = "god"
	// RoleAdmin have access to everything on its group
	RoleAdmin = "admin"
	// RoleUser have access to everything that was created by himself
	RoleUser = "user"
	// RoleCustomer have access
	RoleCustomer = "customer"
)

// User is a generic store user
type User interface {
	GetID() string
	GetGroupID() string
}

// Group is a generic store group
type Group interface {
	GetID() string
	GetRole() UserRole
	GetOrgID() UserRole
}

// Context is an advanced context for authenticated users
type Context struct {
	C        *gin.Context
	Store    Store
	IsLogged bool
	User     User
	Role     UserRole
	Group    Group
}

// AuthContext retrieves the authenticated context for current store user
func AuthContext(c *gin.Context) *Context {
	storedUser, userExists := c.Get(CurrentUserKey)
	storedUserGroup, userGroupExists := c.Get(CurrentUserGroupKey)

	var user User
	if userExists {
		user = storedUser.(User)
	}

	var group Group
	role := ""
	if userGroupExists {
		group = storedUserGroup.(Group)
		role = group.GetRole()
	}

	return &Context{
		C:        c,
		Store:    FromContext(c),
		IsLogged: userExists && userGroupExists,
		User:     user,
		Role:     role,
		Group:    group,
	}
}

// GetCache returns the cached value and true or nil and false if not found (request scoped cache)
func (c *Context) GetCache(key string) (interface{}, bool) {
	if cachedValue, found := c.C.Get("store:" + key); found {
		return cachedValue, true
	}
	return nil, false
}

// SetCache sets a cache value (request scoped cache)
func (c *Context) SetCache(key string, value interface{}) {
	c.C.Set("store:"+key, value)
}

type user struct{}

func (u user) GetID() string      { return "" }
func (u user) GetGroupID() string { return "fline" }

type group struct{}

func (g group) GetID() string      { return "" }
func (g group) GetRole() UserRole  { return RoleGod }
func (g group) GetOrgID() UserRole { return "fline" }

// NewGodContext creates a god in memory context
func NewGodContext(s Store) *Context {
	return &Context{
		C:        new(gin.Context),
		Store:    s,
		IsLogged: true,
		User:     user{},
		Role:     RoleGod,
		Group:    group{},
	}
}
