package migration

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/market-data/internal/config"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
)

// Migrator handles database migrations
type Migrator struct {
	config     *config.MigrationsConfig
	connString string
}

// NewMigrator creates a new migrator instance
func NewMigrator(migrationsConfig *config.MigrationsConfig, connString string) *Migrator {
	return &Migrator{
		config:     migrationsConfig,
		connString: connString,
	}
}

// RunMigrations runs all pending migrations
func (m *Migrator) RunMigrations() error {
	if !m.config.Enabled {
		log.Info().Msg("Database migrations are disabled")
		return nil
	}

	log.Info().Str("path", m.config.Path).Msg("Running database migrations")

	// Create a new migrate instance
	migrator, err := m.createMigrator()
	if err != nil {
		return eris.Wrap(err, "failed to create migrator")
	}
	defer func() {
		sourceErr, dbErr := migrator.Close()
		if sourceErr != nil {
			log.Error().Err(sourceErr).Msg("Error closing migration source")
		}
		if dbErr != nil {
			log.Error().Err(dbErr).Msg("Error closing migration database")
		}
	}()

	// Run migrations
	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return eris.Wrap(err, "failed to run migrations")
	}

	log.Info().Msg("Database migrations completed successfully")
	return nil
}

// createMigrator creates a new migrate instance
func (m *Migrator) createMigrator() (*migrate.Migrate, error) {
	sourceURL := fmt.Sprintf("file://%s", m.config.Path)
	migrator, err := migrate.New(sourceURL, m.connString)
	if err != nil {
		return nil, eris.Wrap(err, "failed to create migrator")
	}

	return migrator, nil
}

// MigrateDown rolls back all migrations
func (m *Migrator) MigrateDown() error {
	if !m.config.Enabled {
		log.Info().Msg("Database migrations are disabled")
		return nil
	}

	log.Info().Str("path", m.config.Path).Msg("Rolling back all database migrations")

	// Create a new migrate instance
	migrator, err := m.createMigrator()
	if err != nil {
		return eris.Wrap(err, "failed to create migrator")
	}
	defer func() {
		sourceErr, dbErr := migrator.Close()
		if sourceErr != nil {
			log.Error().Err(sourceErr).Msg("Error closing migration source")
		}
		if dbErr != nil {
			log.Error().Err(dbErr).Msg("Error closing migration database")
		}
	}()

	// Run migrations down
	if err := migrator.Down(); err != nil && err != migrate.ErrNoChange {
		return eris.Wrap(err, "failed to roll back migrations")
	}

	log.Info().Msg("Database migrations rolled back successfully")
	return nil
}
