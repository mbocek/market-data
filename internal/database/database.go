package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"

	"github.com/market-data/internal/config"
)

// DB represents a database connection pool
type DB struct {
	Pool *pgxpool.Pool
}

// NewWithConfig creates a new database connection pool using the provided configuration
func NewWithConfig(dbConfig *config.DatabaseConfig) (*DB, error) {
	// Create connection string
	connString := dbConfig.GetConnectionString()

	// Configure connection pool
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, eris.Wrap(err, "failed to parse database connection string")
	}

	// Set max connections
	if dbConfig.MaxConnections > 0 {
		config.MaxConns = int32(dbConfig.MaxConnections)
	}

	// Set connection timeout
	if dbConfig.ConnectionTimeout > 0 {
		config.ConnConfig.ConnectTimeout = dbConfig.GetConnectionTimeout()
	}

	// Create connection pool
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, eris.Wrap(err, "failed to create database connection pool")
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		return nil, eris.Wrap(err, "failed to ping database")
	}

	log.Info().
		Str("host", dbConfig.Host).
		Int("port", dbConfig.Port).
		Str("user", dbConfig.User).
		Str("dbname", dbConfig.DBName).
		Msg("Connected to database")

	return &DB{Pool: pool}, nil
}


// Close closes the database connection pool
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
		log.Info().Msg("Database connection pool closed")
	}
}

// Ping checks if the database is reachable
func (db *DB) Ping(ctx context.Context) error {
	if err := db.Pool.Ping(ctx); err != nil {
		return eris.Wrap(err, "failed to ping database")
	}
	return nil
}

// ExecContext executes a query without returning any rows
func (db *DB) ExecContext(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return db.Pool.Exec(ctx, sql, args...)
}

// QueryContext executes a query that returns rows
func (db *DB) QueryContext(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return db.Pool.Query(ctx, sql, args...)
}

// QueryRowContext executes a query that is expected to return at most one row
func (db *DB) QueryRowContext(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return db.Pool.QueryRow(ctx, sql, args...)
}
