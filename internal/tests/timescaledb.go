package tests

import (
	"context"
	"embed"
	"fmt"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/market-data/internal/config"
	"github.com/market-data/internal/database/migration"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	defaultDBName     = "trading_db"                        // Name for the test database
	defaultDBUser     = "postgres"                          // Default DB user
	defaultDBPassword = "your_password_here"                // Default DB password used in tests
	defaultPort       = "5432"                              // Database port inside container
	imageTag          = "timescale/timescaledb:latest-pg17" // Docker image tag for TimescaleDB
)

// TimescaleDB encapsulates the lifecycle of a temporary TimescaleDB Docker container
// and exposes helpers for integration testing.
type TimescaleDB struct {
	ctx       context.Context             // Context for managing container and DB resources
	container *postgres.PostgresContainer // The running database container
	dsn       string                      // Connection string for the database
	conn      *pgx.Conn                   // Active PostgreSQL connection
	testingT  *testing.T                  // Test instance to report errors
}

// NewTimescaleDB spins up a new TimescaleDB container for testing purposes,
// establishes a connection, and prepares the helper struct.
// Errors cause fatal test failures.
func NewTimescaleDB(t *testing.T) *TimescaleDB {
	ctx := context.Background()

	// Start a new TimescaleDB container with specified credentials and settings
	pgContainer, err := postgres.Run(ctx,
		imageTag,
		postgres.WithDatabase(defaultDBName),
		postgres.WithUsername(defaultDBUser),
		postgres.WithPassword(defaultDBPassword),
		postgres.BasicWaitStrategies(),
	)
	failIfErr(t, err, "failed to start container")

	// Retrieve the mapped host and port from the running container
	host, err := pgContainer.Host(ctx)
	failIfErr(t, err, "cannot get container host")
	port, err := pgContainer.MappedPort(ctx, defaultPort)
	failIfErr(t, err, "cannot get mapped port for container")
	// Construct the DSN string for connecting to the database
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", defaultDBUser, defaultDBPassword, host, port.Port(), defaultDBName)

	// Open a connection to verify the database is ready
	conn, err := pgx.Connect(ctx, dsn)
	failIfErr(t, err, "failed to connect to db")

	return &TimescaleDB{
		ctx:       ctx,
		container: pgContainer,
		dsn:       dsn,
		conn:      conn,
		testingT:  t,
	}
}

// ApplyMigrations applies database migrations to the temporary test database using the provided migrations embed.FS.
func (db *TimescaleDB) ApplyMigrations(migrations embed.FS) {
	err := migration.NewMigrator(true, db.DSN(), migrations).RunMigrations()
	failIfErr(db.testingT, err, "failed to run migrations")
}

// ApplyMigrationsWitTableName applies database migrations using the specified table name to track migration history.
func (db *TimescaleDB) ApplyMigrationsWitTableName(migrations embed.FS, tableName string) {
	dsn := fmt.Sprintf("%s&x-migrations-table=%s", db.DSN(), tableName)
	err := migration.NewMigrator(true, dsn, migrations).RunMigrations()
	failIfErr(db.testingT, err, "failed to run migrations")
}

// Terminate shuts down and removes the database container.
// Use this to clean up resources after test execution.
func (db *TimescaleDB) Terminate() {
	err := db.container.Terminate(db.ctx)
	failIfErr(db.testingT, err, "failed to terminate container")
}

// DSN returns the PostgreSQL connection string for the test container.
func (db *TimescaleDB) DSN() string {
	return db.dsn
}

// Host returns the dynamically-assigned host of the running container.
func (db *TimescaleDB) Host() string {
	host, err := db.container.Host(db.ctx)
	failIfErr(db.testingT, err, "failed to get container host")
	return host
}

// Port returns the dynamically-assigned database port.
func (db *TimescaleDB) Port() int {
	port, err := db.container.MappedPort(db.ctx, defaultPort)
	failIfErr(db.testingT, err, "failed to get container port")
	return port.Int()
}

// failIfErr is a helper to report fatal test errors in a consistent manner.
func failIfErr(t *testing.T, err error, msg string) {
	if err != nil {
		t.Fatalf("%s: %v", msg, err)
	}
}

// DatabaseConfigForTimescale constructs a config.DatabaseConfig for the running TimescaleDB instance,
// resetting host/port dynamically and using default credentials.
func (db *TimescaleDB) DatabaseConfigForTimescale() *config.DatabaseConfig {
	return &config.DatabaseConfig{
		Host:              db.Host(),
		Port:              db.Port(),
		User:              defaultDBUser,
		Password:          defaultDBPassword,
		DBName:            defaultDBName,
		SSLMode:           "disable",
		MaxConnections:    3,
		ConnectionTimeout: 5,
	}
}
