package tests

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/market-data/internal/config"
	"github.com/market-data/internal/database/migration"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"testing"
)

type TimescaleDB struct {
	ctx  context.Context
	c    *postgres.PostgresContainer
	dsn  string
	conn *pgx.Conn
	t    *testing.T
}

func NewTimescaleDB(t *testing.T) *TimescaleDB {
	ctx := context.Background()

	// Start a TimescaleDB container using the Postgres module
	pgContainer, err := postgres.Run(ctx,
		"timescale/timescaledb:latest-pg17",
		postgres.WithDatabase("trading_db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("your_password_here"),
	)
	if err != nil {
		t.Fatalf("failed to start container: %v", err)
	}

	host, _ := pgContainer.Host(ctx)
	port, _ := pgContainer.MappedPort(ctx, "5432")
	dsn := fmt.Sprintf("postgres://postgres:your_password_here@%s:%s/trading_db?sslmode=disable", host, port.Port())

	// Connect and perform a basic operation
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}

	return &TimescaleDB{
		ctx:  ctx,
		c:    pgContainer,
		dsn:  dsn,
		conn: conn,
		t:    t,
	}
}

func (t *TimescaleDB) ApplyMigrations(config *config.MigrationsConfig) {
	err := migration.NewMigrator(config, t.dsn).RunMigrations()
	if err != nil {
		t.t.Fatalf("failed to run migrations: %v", err)
	}
}

func (t *TimescaleDB) Terminate() {
	err := t.c.Terminate(t.ctx)
	if err != nil {
		t.t.Fatalf("failed to terminate container: %v", err)
	}
}
