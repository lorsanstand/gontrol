package migrator

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

type Migrator struct {
	srcDriver source.Driver
}

// New создает новый экземпляр Migrator
func New(sqlFiles embed.FS, dirName string) (*Migrator, error) {
	d, err := iofs.New(sqlFiles, dirName)
	if err != nil {
		return nil, fmt.Errorf("create driver: %w", err)
	}

	return &Migrator{srcDriver: d}, nil
}

func (m *Migrator) Apply(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("create db instance: %w", err)
	}

	migrator, err := migrate.NewWithInstance("migration_embeded_sql_files", m.srcDriver, "postgres", driver)
	if err != nil {
		return fmt.Errorf("create migration: %v", err)
	}
	defer migrator.Close()

	err = migrator.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("apply migration: %w", err)
	}

	return nil
}
