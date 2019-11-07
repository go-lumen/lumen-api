package postgresql

import (
	"github.com/jinzhu/gorm"
)

type postgres struct {
	*gorm.DB
}

// New creates a database connexion
func New(database *gorm.DB) *postgres {
	return &postgres{database}
}
