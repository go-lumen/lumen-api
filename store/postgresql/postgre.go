package postgresql

import (
	"context"
	"gorm.io/gorm"
)

type PSQL struct {
	database *gorm.DB
	dbName   string
	context  context.Context
}

// New creates a database connexion
func New(context context.Context, database *gorm.DB, dbName string) *PSQL {
	return &PSQL{database, dbName, context}
}
