package postgresql

import (
	"github.com/go-pg/pg"
)

type postgres struct {
	*pg.DB
}

// New creates a database connexion
func New(database *pg.DB) *postgres {
	return &postgres{database}
}
