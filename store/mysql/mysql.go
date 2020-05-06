package mysql

import (
	"database/sql"
)

// Mysql holds db
type Mysql struct {
	*sql.DB
}

// New creates a database connexion
func New(database *sql.DB) *Mysql {
	return &Mysql{database}
}
