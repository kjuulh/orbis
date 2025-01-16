package persistence

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	migratepgx "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

const migrationSource = "migrations"

//go:embed migrations/*.sql
var migrations embed.FS

func Migrate() error {
	db, err := sql.Open("pgx", os.Getenv("ORBIS_POSTGRES_DB"))
	if err != nil {
		return fmt.Errorf("failed to establish connection to database: %w", err)
	}

	driver, err := migratepgx.WithInstance(db, &migratepgx.Config{})
	if err != nil {
		return fmt.Errorf("failed to install postgres driver for migrations: %w", err)
	}

	migrationSource, err := NewEmbedDriver(migrationSource, migrations)
	if err != nil {
		return fmt.Errorf("failed to setup embedded driver for migrations: %w", err)
	}

	migration, err := migrate.NewWithInstance("iofs", migrationSource, "postgres", driver)
	if err != nil {
		return err
	}

	if err := migration.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to migrate database: %w", err)
		}
	}

	return nil
}

func NewEmbedDriver(path string, files embed.FS) (source.Driver, error) {
	driver, err := iofs.New(files, migrationSource)
	if err != nil {
		return nil, err
	}

	return driver, nil
}
