package storage

import (
	"ct/db/migrations"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
)

// DoMigrateDb performs DB migrations.
func DoMigrateDb(dbUrl string) error {
	resources := bindata.Resource(migrations.AssetNames(),
		func(name string) ([]byte, error) {
			return migrations.Asset(name)
		})

	migrationData, err := bindata.WithInstance(resources)
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("go-bindata", migrationData, dbUrl)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
