package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/go-lumen/lumen-api/server"
)

var (
	migrations []*gormigrate.Migration
)

// Migrator contains api server
type Migrator struct {
	api *server.API
}

// New initiates Migrator struct
func New(api *server.API) *Migrator {
	return &Migrator{
		api: api,
	}
}

// Migrate runs migrations
func (m *Migrator) Migrate() error {
	gm := gormigrate.New(m.api.PostgreDatabase, gormigrate.DefaultOptions, migrations)

	if err := gm.Migrate(); err != nil {
		return err
	}
	return nil
}

// MigrateTo executes all migrations that did not run yet up to
// the migration that matches `migrationID`.
func (m *Migrator) MigrateTo(migrationID string) error {
	gm := gormigrate.New(m.api.PostgreDatabase, gormigrate.DefaultOptions, migrations)

	if err := gm.MigrateTo(migrationID); err != nil {
		return err
	}
	return nil
}

// RollbackLast undo the last migration
func (m *Migrator) RollbackLast() error {
	gm := gormigrate.New(m.api.PostgreDatabase, gormigrate.DefaultOptions, migrations)
	if err := gm.RollbackLast(); err != nil {
		return err
	}
	return nil
}

// RollbackTo undoes migrations up to the given migration that matches the `migrationID`.
// Migration with the matching `migrationID` is not rolled back.
func (m *Migrator) RollbackTo(migrationID string) error {
	gm := gormigrate.New(m.api.PostgreDatabase, gormigrate.DefaultOptions, migrations)
	if err := gm.RollbackTo(migrationID); err != nil {
		return err
	}
	return nil
}
