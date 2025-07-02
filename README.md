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
- Testing with Testify
- Docker and Docker Compose support

## Project Structure

```
market-data/
├── cmd/
│   └── market-data/     # Main application entry point
├── config/              # Configuration files
├── db/                  # Database initialization scripts
├── internal/            # Private application code
│   ├── config/          # Configuration management
│   ├── controller/      # HTTP controllers
│   ├── data/            # Data models and repository
│   └── database/        # Database connection and utilities
├── pkg/                 # Public libraries
│   └── marketservice/   # Market data service
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

The service uses TimescaleDB, a PostgreSQL extension optimized for time-series data. When running with Docker Compose, the database is automatically set up with the schema defined in `db/init.sql`.

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
