# Market-Data Module Guidelines

## Project Overview
The market-data module is a Go-based microservice for retrieving and processing market data for trading applications. It provides real-time and historical market data for financial instruments, using TimescaleDB (a PostgreSQL extension optimized for time-series data) for data storage.

## Project Structure
The market-data project follows a clean architecture approach with the following structure:

```
market-data/
├── cmd/
│   └── market-data/     # Main application entry point
├── config/              # Configuration files
├── db/                  # Database initialization scripts
├── internal/            # Private application code
│   ├── config/          # Configuration management
│   ├── database/        # Database connection and utilities
│   ├── domain/          # Domain models and business logic
│   │   └── market/      # Market data domain
│   └── interfaces/      # Interface adapters
│       └── api/         # API controllers
├── Dockerfile           # Container definition
├── docker-compose.yml   # Container orchestration
├── Makefile             # Build automation
└── README.md            # Project documentation
```

## Development Guidelines

### Code Organization
- Follow the clean architecture principles
- Keep domain logic separate from infrastructure concerns
- Use dependency injection for better testability

### Error Handling
- Use the Eris package for error handling
- Wrap errors with context information
- Log errors appropriately using Zerolog

### Logging
- Use structured logging with Zerolog
- Include relevant context in log entries
- Use appropriate log levels (debug, info, warn, error)

### Configuration
- Use Viper for configuration management
- Support configuration via files, environment variables, and flags
- Validate configuration at startup
- Configuration includes server settings, database connection parameters, and logging options

## API Endpoints
The service provides the following API endpoints:

- `GET /` - Service status
- `GET /health` - Health check endpoint
- `GET /symbols/:symbol` - Get market data for a specific symbol

## Requirements
- Go 1.24 or higher
- Docker
- Docker Compose

## Building and Running
The project includes a Makefile with common commands:

```bash
# Build the binary
make build

# Run the application
make run

# Clean build files
make clean

# Run tests
make test

# Build Docker image
make docker-build

# Run Docker container
make docker-run

# Start services with docker compose
make docker-compose-up

# Stop services with docker compose
make docker-compose-down
```

## Testing
Currently, the project doesn't have comprehensive automated tests. When implementing tests:
- Use the Testify package for assertions and mocks
- Write unit tests for individual components
- Write integration tests for the complete flow
- Use dependency injection to facilitate testing

## Contribution Guidelines
When contributing to this project:

1. **Understand the Architecture**:
   - Recognize the clean architecture approach in the project

2. **Follow Code Style**:
   - Adhere to Go best practices

3. **Testing**:
   - If implementing or modifying tests, use the appropriate testing frameworks
   - Run tests when available to verify changes

4. **Building**:
   - Use the provided Makefile commands

5. **Error Handling**:
   - Use Eris for error handling

6. **Logging**:
   - Use structured logging with Zerolog