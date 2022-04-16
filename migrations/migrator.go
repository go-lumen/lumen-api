package migrations

import (
	"github.com/adrien3d/stokelp-poc/server"
	"gopkg.in/gormigrate.v1"
)

var (
	migrations []*gormigrate.Migration
)

type migrator struct {
	api *server.API
}

func New(api *server.API) *migrator {
	return &migrator{
		api: api,
	}
}

// Run migration
func (m *migrator) Migrate() error {
	gm := gormigrate.New(m.api.PostgreDatabase, gormigrate.DefaultOptions, migrations)

	if err := gm.Migrate(); err != nil {
		return err
	}
	return nil
}

// MigrateTo executes all migrations that did not run yet up to
// the migration that matches `migrationID`.
func (m *migrator) MigrateTo(migrationID string) error {
	gm := gormigrate.New(m.api.PostgreDatabase, gormigrate.DefaultOptions, migrations)

	if err := gm.MigrateTo(migrationID); err != nil {
		return err
	}
	return nil
}

// RollbackLast undo the last migration
func (m *migrator) RollbackLast() error {
	gm := gormigrate.New(m.api.PostgreDatabase, gormigrate.DefaultOptions, migrations)
	if err := gm.RollbackLast(); err != nil {
		return err
	}
	return nil
}

// RollbackTo undoes migrations up to the given migration that matches the `migrationID`.
// Migration with the matching `migrationID` is not rolled back.
func (m *migrator) RollbackTo(migrationID string) error {
	gm := gormigrate.New(m.api.PostgreDatabase, gormigrate.DefaultOptions, migrations)
	if err := gm.RollbackTo(migrationID); err != nil {
		return err
	}
	return nil
}
