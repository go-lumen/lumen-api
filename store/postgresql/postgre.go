package postgresql

import (
	"github.com/jinzhu/gorm"
)

// Postgres holds db
type Postgres struct {
	*gorm.DB
}

// New creates a database connexion
func New(database *gorm.DB) *Postgres {
	return &Postgres{database}
}
