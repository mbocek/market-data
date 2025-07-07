# Variables
BINARY_NAME=market-data
DOCKER_IMAGE=market-data
DOCKER_TAG=latest

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOTMP=$(GOBASE)/tmp
GOFILES=$(wildcard *.go)

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
	go test -v ./...

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

## migrate-up: Run all pending migrations
migrate-up:
	@echo "Running migrations..."
	migrate -path db/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" up

## migrate-down: Roll back the most recent migration
migrate-down:
	@echo "Rolling back migration..."
	migrate -path db/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" down 1

## migrate-reset: Roll back all migrations
migrate-reset:
	@echo "Rolling back all migrations..."
	migrate -path db/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" down

## help: Display available commands
help:
	@echo "Available commands:"
	@grep -E '^##' Makefile | sed -e 's/## //g'

# Database connection variables (override these with environment variables if needed)
DB_USER ?= postgres
DB_PASSWORD ?= your_password_here
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_NAME ?= trading_db
DB_SSLMODE ?= disable

.PHONY: build run clean test docker-build docker-run docker-compose-up docker-compose-down migrate-create migrate-up migrate-down migrate-reset help
