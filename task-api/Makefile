.PHONY: help build run docker-up docker-down test lint clean install-deps swagger

# Variables
APP_NAME = task-management-api
MAIN_PATH = ./cmd/main.go
BINARY_PATH = ./bin/$(APP_NAME)

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

install-deps: ## Install Go dependencies
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

build: ## Build the application
	@echo "Building $(APP_NAME)..."
	mkdir -p bin
	go build -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "Build complete: $(BINARY_PATH)"

run: ## Run the application locally
	@echo "Running $(APP_NAME)..."
	go run $(MAIN_PATH)

dev: ## Run with hot reload (requires air)
	@which air > /dev/null || go install github.com/cosmtrek/air@latest
	air

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	go clean

test: ## Run tests
	@echo "Running tests..."
	go test -v -cover ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint: ## Run linters
	@which golangci-lint > /dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run ./...

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

swagger: ## Generate Swagger documentation
	@which swag > /dev/null || go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Generating Swagger docs..."
	swag init
	@echo "Swagger docs generated"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(APP_NAME):latest .
	@echo "Docker image built: $(APP_NAME):latest"

docker-up: ## Start containers with docker-compose
	@echo "Starting containers..."
	docker-compose up -d
	@echo "Containers started. API available at http://localhost:8080"

docker-down: ## Stop docker-compose containers
	@echo "Stopping containers..."
	docker-compose down

docker-logs: ## Show docker-compose logs
	docker-compose logs -f api

docker-clean: ## Remove all Docker artifacts
	@echo "Cleaning Docker artifacts..."
	docker-compose down -v
	docker rmi $(APP_NAME):latest

db-reset: ## Reset database (clears all data)
	@echo "Resetting database..."
	docker-compose exec postgres psql -U postgres -d postgres -c "DROP DATABASE IF EXISTS task_management_db;"
	docker-compose exec postgres psql -U postgres -d postgres -c "CREATE DATABASE task_management_db;"
	@echo "Database reset complete"

db-shell: ## Open PostgreSQL shell
	docker-compose exec postgres psql -U postgres -d task_management_db

# Install development tools
install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Tools installed"

# Quick start
quick-start: install-deps docker-build docker-up ## Quick start: install deps, build, and run
	@echo "Quick start complete! API available at http://localhost:8080"

all: clean install-deps fmt vet lint test build ## Run all checks and build
	@echo "All checks passed!"
