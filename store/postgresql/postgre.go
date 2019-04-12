package postgresql

import (
	"database/sql"
)

type postgre struct {
	*sql.DB
}

// New creates a database connexion
func New(database *sql.DB) *postgre {
	return &postgre{database}
}
