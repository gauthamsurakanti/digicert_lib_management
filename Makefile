# Variables
APP_NAME=library-management
DOCKER_IMAGE=$(APP_NAME):latest
DOCKER_COMPOSE_FILE=docker-compose.yml

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

.PHONY: help build run test clean docker-build docker-up docker-down docker-logs migrate-up migrate-down

# Default target
help: ## Show this help message
	@echo "$(BLUE)Library Management API - Available Commands:$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

# Development commands
build: ## Build the application
	@echo "$(YELLOW)Building application...$(NC)"
	go build -o bin/$(APP_NAME) ./cmd/api

run: ## Run the application locally
	@echo "$(YELLOW)Running application...$(NC)"
	go run ./cmd/api

test: ## Run tests
	@echo "$(YELLOW)Running tests...$(NC)"
	go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "$(YELLOW)Running tests with coverage...$(NC)"
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

clean: ## Clean build artifacts
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	rm -f bin/$(APP_NAME)
	rm -f coverage.out coverage.html

# Docker commands
docker-build: ## Build Docker image
	@echo "$(YELLOW)Building Docker image...$(NC)"
	docker build -t $(DOCKER_IMAGE) .

docker-up: ## Start all services with Docker Compose
	@echo "$(YELLOW)Starting services with Docker Compose...$(NC)"
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d
	@echo "$(GREEN)Services started! API available at http://localhost:8080$(NC)"

docker-down: ## Stop all services
	@echo "$(YELLOW)Stopping services...$(NC)"
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

docker-logs: ## View logs from all services
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

docker-restart: ## Restart all services
	@echo "$(YELLOW)Restarting services...$(NC)"
	docker-compose -f $(DOCKER_COMPOSE_FILE) restart

docker-rebuild: ## Rebuild and restart services
	@echo "$(YELLOW)Rebuilding and restarting services...$(NC)"
	docker-compose -f $(DOCKER_COMPOSE_FILE) down
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d --build

# Database commands
migrate-up: ## Run database migrations up
	@echo "$(YELLOW)Running migrations up...$(NC)"
	migrate -path migrations -database "postgres://library_user:library_pass@localhost:5432/library_db?sslmode=disable" up

migrate-down: ## Run database migrations down
	@echo "$(YELLOW)Running migrations down...$(NC)"
	migrate -path migrations -database "postgres://library_user:library_pass@localhost:5432/library_db?sslmode=disable" down

migrate-create: ## Create a new migration (usage: make migrate-create name=migration_name)
	@echo "$(YELLOW)Creating migration: $(name)$(NC)"
	migrate create -ext sql -dir migrations $(name)

# Development helpers
dev-setup: ## Set up development environment
	@echo "$(YELLOW)Setting up development environment...$(NC)"
	cp .env.example .env
	@echo "$(GREEN)Environment file created. Please edit .env with your settings.$(NC)"

fmt: ## Format Go code
	@echo "$(YELLOW)Formatting code...$(NC)"
	go fmt ./...

lint: ## Run linter
	@echo "$(YELLOW)Running linter...$(NC)"
	golangci-lint run

deps: ## Download dependencies
	@echo "$(YELLOW)Downloading dependencies...$(NC)"
	go mod download
	go mod tidy

# API testing commands
test-api: ## Test API endpoints (requires running server)
	@echo "$(YELLOW)Testing API endpoints...$(NC)"
	@echo "$(BLUE)Health check:$(NC)"
	curl -s http://localhost:8080/health | jq .
	@echo "\n$(BLUE)Get all books:$(NC)"
	curl -s http://localhost:8080/api/v1/books | jq .
	@echo "\n$(BLUE)Get book by ID:$(NC)"
	curl -s http://localhost:8080/api/v1/books/1 | jq .

create-book: ## Create a test book
	@echo "$(YELLOW)Creating test book...$(NC)"
	curl -X POST http://localhost:8080/api/v1/books \
		-H "Content-Type: application/json" \
		-d '{"title":"Test Book","author":"Test Author","isbn":"978-1234567890","publisher":"Test Publisher","publish_year":2024,"genre":"Test","pages":100,"description":"A test book"}' | jq .

# Production commands
deploy: docker-build docker-up ## Build and deploy the application

status: ## Show status of all services
	@echo "$(BLUE)Service Status:$(NC)"
	docker-compose -f $(DOCKER_COMPOSE_FILE) ps