package tests

import (
	"context"
	"embed"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/market-data/internal/config"
	"github.com/market-data/internal/database/migration"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"testing"
)

const (
	defaultDBName     = "trading_db"
	defaultDBUser     = "postgres"
	defaultDBPassword = "your_password_here"
	defaultPort       = "5432"
	imageTag          = "timescale/timescaledb:latest-pg17"
)

type TimescaleDB struct {
	ctx       context.Context
	container *postgres.PostgresContainer
	dsn       string
	conn      *pgx.Conn
	testingT  *testing.T
}

func NewTimescaleDB(t *testing.T) *TimescaleDB {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		imageTag,
		postgres.WithDatabase(defaultDBName),
		postgres.WithUsername(defaultDBUser),
		postgres.WithPassword(defaultDBPassword),
		postgres.BasicWaitStrategies(),
	)
	failIfErr(t, err, "failed to start container")

	host, _ := pgContainer.Host(ctx)
	port, _ := pgContainer.MappedPort(ctx, defaultPort)
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", defaultDBUser, defaultDBPassword, host, port.Port(), defaultDBName)

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

func (db *TimescaleDB) ApplyMigrations(migrations embed.FS) {
	err := migration.NewMigrator(true, db.DSN(), migrations).RunMigrations()
	failIfErr(db.testingT, err, "failed to run migrations")
}

func (db *TimescaleDB) Terminate() {
	err := db.container.Terminate(db.ctx)
	failIfErr(db.testingT, err, "failed to terminate container")
}

func (db *TimescaleDB) DSN() string {
	return db.dsn
}

func (db *TimescaleDB) Host() string {
	host, err := db.container.Host(db.ctx)
	failIfErr(db.testingT, err, "failed to get container host")
	return host
}

func (db *TimescaleDB) Port() int {
	port, err := db.container.MappedPort(db.ctx, defaultPort)
	failIfErr(db.testingT, err, "failed to get container port")
	return port.Int()
}

func failIfErr(t *testing.T, err error, msg string) {
	if err != nil {
		t.Fatalf("%s: %v", msg, err)
	}
}

// DatabaseConfigForTimescale returns a DatabaseConfig using this container's dynamic host/port and project defaults.
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
