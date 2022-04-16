package postgresql

import (
	"github.com/jinzhu/gorm"
)

type postgresql struct {
	*gorm.DB
}

// New creates a database connexion
func New(database *gorm.DB) *postgresql {
	return &postgresql{database}
}
