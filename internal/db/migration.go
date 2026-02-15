package db

import (
	"database/sql"
	"embed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var fs embed.FS

func setupMigrations(db *sql.DB) (*migrate.Migrate, error) {
	d, err := iofs.New(fs, "migrations")
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithInstance("iofs", d, "db", driver)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// UpgradeDB runs migrations all the way to the latest possible version
func UpgradeDB(db *sql.DB) (uint, error) {
	m, err := setupMigrations(db)
	if err != nil {
		return 0, err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return 0, err
	}
	v, _, err := m.Version()
	return v, err
}

// MigrateDBTo runs migrations to a specific version. To undo all migrations, pass target version 0
func MigrateDBTo(db *sql.DB, targetVersion uint) (version uint, err error) {
	m, err := setupMigrations(db)
	if err != nil && err != migrate.ErrNoChange {
		return
	}
	if targetVersion == 0 {
		err = m.Down()
	} else {
		err = m.Migrate(targetVersion)
	}

	if err != nil && err != migrate.ErrNoChange {
		return
	}

	return targetVersion, nil
}

// Version returns the current version of the db
func Version(db *sql.DB) (uint, bool, error) {
	m, err := setupMigrations(db)
	if err != nil && err != migrate.ErrNoChange {
		return 0, false, err
	}
	return m.Version()
}

// Force changes DB version and resets its dirty flag without running migrations.
func Force(db *sql.DB, targetVersion int) error {
	m, err := setupMigrations(db)
	if err != nil {
		return err
	}
	return m.Force(targetVersion)
}
