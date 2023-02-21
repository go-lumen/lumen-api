package store

import (
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Model represents a generic store model
type Model interface {
	//RoleChecker
}

// GenericModel represents a generic mongo model
type GenericModel interface {
	Model
	GetCollection() string
}

// EnsureGenericModel ensures that a model implement GenericModel interface
func EnsureGenericModel(model Model) GenericModel {
	genericModel, ok := model.(GenericModel)
	if !ok {
		panic("when using driver, you should implement MongoModel interface on your model")
	}
	return genericModel
}

// MongoFindAllOptioner represents default mongo FindOptions
type MongoFindAllOptioner interface {
	ApplyOptions(*options.FindOptions)
}
