# Market Data Service

A microservice for retrieving and processing market data for trading applications.

## Overview

This service provides real-time and historical market data for financial instruments. It's built with Go and containerized with Docker for easy deployment. The service uses TimescaleDB (a PostgreSQL extension optimized for time-series data) for data storage, making it ideal for storing and querying financial market data.

## Features

- RESTful API using Gin web framework
- Time-series data storage with TimescaleDB
- Structured logging with Zerolog
- Configuration management with Viper
- Error handling with Eris
- Database integration with pgx
- Database migrations with golang-migrate
- Code quality enforcement with golangci-lint
- Testing with Testify
- Docker and Docker Compose support

## Project Structure

```
market-data/
├── cmd/
│   └── market-data/     # Main application entry point
├── config/              # Configuration files
├── db/                  # Database scripts
│   └── migrations/      # Database migration files
├── internal/            # Private application code
│   ├── config/          # Configuration management
│   ├── database/        # Database connection and utilities
│   │   └── migration/   # Database migration functionality
│   ├── domain/          # Domain models and business logic
│   │   └── market/      # Market data domain
│   └── interfaces/      # Interface adapters
│       └── api/         # API controllers
├── tmp/                 # Build artifacts (gitignored)
├── Dockerfile           # Container definition
├── docker-compose.yml   # Container orchestration
├── Makefile             # Build automation
└── README.md            # This file
```

## Requirements

- Go 1.24 or higher
- Docker
- Docker Compose

## Building and Running

### Using Make

The project includes a Makefile with common commands:

```bash
# Display available commands
make help

# Build the binary
make build

# Run the application
make run

# Run tests
make test

# Clean build files
make clean

# Install golangci-lint
make lint-install

# Run linter
make lint

# Run linter with auto-fix
make lint-fix

# Build Docker image
make docker-build

# Run Docker container
make docker-run

# Start services with docker compose
make docker-compose-up

# Stop services with docker compose
make docker-compose-down
```

### Local Development

```bash
# Build and run locally
go run cmd/market-data/main.go
```

### Using Docker

```bash
# Build and run with Docker
docker build -t market-data .
docker run -p 8080:8080 market-data
```

### Using Docker Compose

```bash
# Start the service
docker compose up -d

# Stop the service
docker compose down
```

## Database

The service uses TimescaleDB, a PostgreSQL extension optimized for time-series data. The database schema is managed through migrations using the golang-migrate library.

### Database Migrations

The service uses database migrations to manage schema changes. Migrations are stored in the `db/migrations` directory and are automatically applied when the service starts. Each migration consists of two files:

- `<version>_<name>.up.sql` - SQL to apply the migration
- `<version>_<name>.down.sql` - SQL to roll back the migration

The service includes the following migration commands in the Makefile:

```bash
# Create a new migration file
make migrate-create

# Run all pending migrations
make migrate-up

# Roll back the most recent migration
make migrate-down

# Roll back all migrations
make migrate-reset
```

### Database Schema

The database includes the following tables:

- `symbols` - Stores information about financial instruments
- `stock_prices` - Stores time-series price data (converted to a TimescaleDB hypertable)
- `price_fetch_logs` - Logs data fetch operations

### Connecting to the Database

When running with Docker Compose, you can connect to the database using:

```bash
# Connect to the database
docker exec -it timescaledb psql -U postgres -d your_database_name
```

## Linting

The project uses [golangci-lint](https://golangci-lint.run/) for code quality enforcement. The linter is configured in the `.golangci.yml` file at the root of the project.

### Linter Configuration

The linter is configured to run the following checks:
- gofmt - Check whether code was gofmt-ed
- goimports - Check import statements are formatted according to the goimports tool
- misspell - Find commonly misspelled English words in comments
- whitespace - Tool for detection of leading and trailing whitespace

Additional linters are available in the configuration file but are currently disabled due to compatibility issues with the project structure. These can be enabled as needed in the future.

### Running the Linter

The linter can be run using the following make commands:

```bash
# Install golangci-lint
make lint-install

# Run the linter
make lint

# Run the linter with auto-fix for issues that can be automatically fixed
make lint-fix
```

## API Endpoints

- `GET /` - Service status
- `GET /health` - Health check endpoint
- `GET /symbols` - Get all available market data symbols
- `GET /data/{symbol}` - Get market data for a specific symbol

## Configuration

The application uses Viper for configuration management. Configuration can be provided through:

1. Configuration file (`config/config.yaml`)
2. Environment variables
3. Command-line flags

### Configuration File

The default configuration file is located at `config/config.yaml`. It contains settings for:

```yaml
# Server configuration
server:
  port: 8080
  host: "0.0.0.0"

# Database configuration
database:
  host: "timescaledb"
  port: 5432
  user: "postgres"
  password: "your_password_here"
  dbname: "your_database_name"
  sslmode: "disable"
  max_connections: 10
  connection_timeout: 5 # seconds

# Logging configuration
logging:
  level: "info" # debug, info, warn, error
  format: "json" # json, console
  output: "stdout" # stdout, file
  file_path: "logs/market-data.log" # only used if output is file

# Migrations configuration
migrations:
  enabled: true
  path: "db/migrations"
```

### Environment Variables

Environment variables can be used to override configuration settings:

- `PORT` - The port the service listens on (default: 8080)
- `DATABASE_HOST` - Database host (default: "timescaledb")
- `DATABASE_PORT` - Database port (default: 5432)
- `DATABASE_USER` - Database user (default: "postgres")
- `DATABASE_PASSWORD` - Database password
- `DATABASE_NAME` - Database name
- `LOGGING_LEVEL` - Logging level (default: "info")
- `MIGRATIONS_ENABLED` - Enable or disable migrations (default: true)
- `MIGRATIONS_PATH` - Path to migration files (default: "db/migrations")
