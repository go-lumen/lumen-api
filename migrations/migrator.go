package migrations

import (
	"github.com/go-lumen/lumen-api/server"
	"github.com/go-lumen/lumen-api/utils"
	"gopkg.in/gormigrate.v1"
)

var (
	migrations []*gormigrate.Migration
)

// Migrator holds server API
type Migrator struct {
	api *server.API
}

// New gene
func New(api *server.API) *Migrator {
	return &Migrator{
		api: api,
	}
}

// Migrate runs a migration
func (m *Migrator) Migrate() error {
	gm := gormigrate.New(m.api.PostgreDatabase, gormigrate.DefaultOptions, migrations)

	if err := gm.Migrate(); err != nil {
		return err
	}
	utils.Log(nil, "info", "migration OK")

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
