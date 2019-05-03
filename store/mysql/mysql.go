package mysql

import (
	"database/sql"
)

type mysql struct {
	*sql.DB
}

// New creates a database connexion
func New(database *sql.DB) *mysql {
	return &mysql{database}
}
