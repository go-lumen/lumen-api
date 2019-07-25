package postgresql

import (
	"github.com/go-pg/pg"
)

type postgre struct {
	*pg.DB
}

// New creates a database connexion
func New(database *pg.DB) *postgre {
	return &postgre{database}
}
