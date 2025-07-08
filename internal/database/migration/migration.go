package migration

import (
	"embed"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
)

// Migrator handles database migrations
type Migrator struct {
	connString string
	enabled    bool
	migrations embed.FS
}

// NewMigrator creates a new migrator instance
func NewMigrator(enabled bool, connString string, migrations embed.FS) *Migrator {
	return &Migrator{
		enabled:    enabled,
		connString: connString,
		migrations: migrations,
	}
}

// RunMigrations runs all pending migrations
func (m *Migrator) RunMigrations() error {
	if !m.enabled {
		log.Info().Msg("Database migrations are disabled")
		return nil
	}

	log.Info().Msg("Running database migrations")

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
	source, err := iofs.New(m.migrations, "migrations")
	if err != nil {
		return nil, eris.Wrap(err, "cannot open migrations")
	}

	migrator, err := migrate.NewWithSourceInstance("iofs", source, m.connString)
	if err != nil {
		return nil, eris.Wrap(err, "failed to create migrator")
	}

	return migrator, nil
}

// MigrateDown rolls back all migrations
func (m *Migrator) MigrateDown() error {
	if !m.enabled {
		log.Info().Msg("Database migrations are disabled")
		return nil
	}

	log.Info().Msg("Rolling back all database migrations")

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
	if err := migrator.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return eris.Wrap(err, "failed to roll back migrations")
	}

	log.Info().Msg("Database migrations rolled back successfully")
	return nil
}
