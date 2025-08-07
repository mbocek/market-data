# Variables
BINARY_NAME=market-data
DOCKER_IMAGE=market-data
DOCKER_TAG=latest

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOTMP=$(GOBASE)/tmp
GOFILES=$(wildcard *.go)

# Linter variables
GOLANGCI_VERSION=v2.2.1
# Test variables
GOTESTSUM_VERSION=v1.12.3

# Make is verbose in Linux. Make it silent.
#MAKEFLAGS += --silent

## build: Build the binary
build:
	@echo "Building..."
	mkdir -p $(GOTMP)
	go build -o $(GOTMP)/$(BINARY_NAME) ./cmd/market-data

## run: Run the application
run: build
	@echo "Running..."
	$(GOTMP)/$(BINARY_NAME)

## clean: Clean build files
clean:
	@echo "Cleaning..."
	go clean
	rm -rf $(GOTMP)

## test: Run tests
test:
	@echo "Testing..."
	$(GOBIN)/gotestsum --format testname ./...

## docker-build: Build docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

## docker-run: Run docker container
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

## docker-compose-up: Start services with docker compose
docker-compose-up:
	@echo "Starting services with docker compose..."
	docker compose up -d

## docker-compose-down: Stop services with docker compose
docker-compose-down:
	@echo "Stopping services with docker compose..."
	docker compose down

## migrate-create: Create a new migration file
migrate-create:
	@echo "Creating migration file..."
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir db/migrations -seq $$name

## install-tools: Install golangci-lint, gotestsum
install-tools:
	@echo "Installing golangci-lint $(GOLANGCI_VERSION)..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) $(GOLANGCI_VERSION)
	@echo "Installing gotestsum $(GOTESTSUM_VERSION)..."
	@mkdir -p $(GOBIN)
	@if [ ! -f $(GOBIN)/gotestsum ]; then \
		GOBIN=$(GOBIN) go install gotest.tools/gotestsum@$(GOTESTSUM_VERSION); \
	fi


## lint: Run golangci-lint
lint:
	@echo "Running golangci-lint..."
	@$(GOBIN)/golangci-lint version
	@$(GOBIN)/golangci-lint run ./...

## lint-fix: Run golangci-lint with auto-fix
lint-fix:
	@echo "Running golangci-lint with auto-fix..."
	@$(GOBIN)/golangci-lint run --no-config --disable-all \
		--enable=gofmt \
		--enable=goimports \
		--enable=misspell \
		--enable=whitespace \
		--fix ./...

## help: Display available commands
help:
	@echo "Available commands:"
	@grep -E '^##' Makefile | sed -e 's/## //g'

.PHONY: build run clean test docker-build docker-run docker-compose-up docker-compose-down migrate-create migrate-up migrate-down migrate-reset lint lint-fix lint-install help
