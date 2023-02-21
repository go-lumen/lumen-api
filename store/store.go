package store

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// SortOrder represents a sort direction
type SortOrder string

const (
	// SortAscending for sorting in ascending order
	SortAscending SortOrder = "ascending"
	// SortDescending for sorting in ascending order
	SortDescending SortOrder = "descending"
)

// SortOption represents a sorted field
type SortOption struct {
	Field string
	Order SortOrder
}

// FindOptions represents store common optional parameters
type FindOptions struct {
	Cache        bool
	Limit        int64
	HasLimit     bool
	SortedFields []SortOption
}

// FindOption is an updater function used to update FindOptions
type FindOption func(uo *FindOptions)

// WithCache is used to set the Cache flag
func WithCache(value bool) FindOption {
	return func(uo *FindOptions) {
		uo.Cache = value
	}
}

// WithLimit is used to set the Limit flag
func WithLimit(value int64) FindOption {
	return func(uo *FindOptions) {
		uo.Limit = value
		uo.HasLimit = true
	}
}

// WithSort is used to add a SortedFields
func WithSort(field string, order SortOrder) FindOption {
	return func(uo *FindOptions) {
		uo.SortedFields = append(uo.SortedFields, SortOption{
			Field: field,
			Order: order,
		})
	}
}

// GetFindOptions retrieves common options from varargs
func GetFindOptions(opts ...FindOption) FindOptions {
	defaults := FindOptions{
		Cache: false,
	}
	for _, opt := range opts {
		opt(&defaults)
	}
	return defaults
}

// UpdateOptions represents store update optional parameters
type UpdateOptions struct {
	CreateIfNotExists bool
	OnlyFields        []string // update only these fields, default all field except omitempty ones
}

// UpdateOption is an updater function used to update UpdateOptions
type UpdateOption func(uo *UpdateOptions)

// CreateIfNotExists is used to set the CreateIfNotExists flag
func CreateIfNotExists(value bool) UpdateOption {
	return func(uo *UpdateOptions) {
		uo.CreateIfNotExists = value
	}
}

// OnlyFields is used to set the OnlyFields flag
func OnlyFields(value []string) UpdateOption {
	return func(uo *UpdateOptions) {
		uo.OnlyFields = value
	}
}

// GetUpdateOptions retrieves options from varargs
func GetUpdateOptions(opts ...UpdateOption) UpdateOptions {
	defaults := UpdateOptions{
		CreateIfNotExists: false,
		OnlyFields:        []string{},
	}
	for _, opt := range opts {
		opt(&defaults)
	}
	return defaults
}

// ID is a shortcut for creating an id filter
func ID(id string) bson.M {
	return bson.M{"id": id}
}

// Store interface
type Store interface {
	Create(*Context, string, Model) error
	Find(*Context, bson.M, Model, ...FindOption) error
	FindAll(*Context, bson.M, interface{}, ...FindOption) error
	Update(*Context, bson.M, Model, ...UpdateOption) error
	Delete(*Context, string, Model) error
	DeleteAll(*Context, bson.M, Model) (int64, error)
	GetCollection(*Context, Model) *mongo.Collection
}
