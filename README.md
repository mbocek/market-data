# Market Data Service

A microservice for retrieving and processing market data for trading applications.

## Overview

This service provides real-time and historical market data for financial instruments. It's built with Go and containerized with Docker for easy deployment.

## Project Structure

```
market-data/
├── cmd/
│   └── market-data/     # Main application entry point
├── internal/            # Private application code
│   └── data/            # Data processing logic
├── pkg/                 # Public libraries
├── Dockerfile           # Container definition
├── docker-compose.yml   # Container orchestration
└── README.md            # This file
```

## Requirements

- Go 1.21 or higher
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

# Start services with docker-compose
make docker-compose-up

# Stop services with docker-compose
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

## API Endpoints

- `GET /` - Service status
- `GET /health` - Health check endpoint
- `GET /symbols` - Get all available market data symbols
- `GET /data/{symbol}` - Get market data for a specific symbol

## Environment Variables

- `PORT` - The port the service listens on (default: 8080)
