package store

import (
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Model represents a generic store model
type Model interface {
	RoleChecker
}

// MongoModel represents a generic mongo model
type MongoModel interface {
	Model
	GetCollection() string
}

// MongoFindAllOptioner represents default mongo FindOptions
type MongoFindAllOptioner interface {
	ApplyOptions(*options.FindOptions)
}
