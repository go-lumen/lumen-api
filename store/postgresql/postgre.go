package postgresql

import (
	"gorm.io/gorm"
)

// Postgresql contains default DB structure
type Postgresql struct {
	*gorm.DB
}

// New creates a database connexion
func New(database *gorm.DB) *Postgresql {
	return &Postgresql{database}
}
