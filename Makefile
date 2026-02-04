.PHONY: help build run test clean docker-build docker-run keys migrate

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

keys: ## Generate RSA keys for JWT
	@echo "Generating RSA keys..."
	@chmod +x scripts/generate_keys.sh
	@./scripts/generate_keys.sh

build: ## Build the application
	@echo "Building application..."
	@go build -o bin/auth-service cmd/main.go
	@echo "Build complete: bin/auth-service"

run: ## Run the application
	@echo "Running application..."
	@go run cmd/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t auth-service:latest .

docker-run: ## Run with Docker Compose
	@echo "Starting services with Docker Compose..."
	@docker-compose up -d

docker-stop: ## Stop Docker Compose services
	@echo "Stopping services..."
	@docker-compose down

docker-logs: ## Show Docker Compose logs
	@docker-compose logs -f

migrate: ## Run database migrations
	@echo "Running migrations..."
	@go run cmd/main.go migrate

dev: ## Run with hot reload using Air (Linux/Mac)
	@echo "Starting development server with hot reload..."
	@air

dev-windows: ## Run with hot reload using Air (Windows)
	@echo "Starting development server with hot reload..."
	@air -c .air.windows.toml

install-tools: ## Install development tools (Air for hot reload, Atlas for migrations)
	@echo "Installing development tools..."
	@go install github.com/cosmtrek/air@latest
	@go install ariga.io/atlas/cmd/atlas@latest
	@echo "✅ Air installed successfully. Use 'make dev' to run with hot reload."
	@echo "✅ Atlas installed successfully. Use 'make atlas-help' for migration commands."
	@echo "See docs/DEVELOPMENT.md for more information."

db-create: ## Create database
	@echo "Creating database..."
	@psql -U postgres -c "CREATE DATABASE auth_service;"

db-drop: ## Drop database
	@echo "Dropping database..."
	@psql -U postgres -c "DROP DATABASE IF EXISTS auth_service;"

db-reset: db-drop db-create migrate ## Reset database

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

all: clean deps keys build ## Build everything from scratch

# Atlas Migration Commands
atlas-help: ## Show Atlas migration commands
	@echo "Atlas Migration Commands:"
	@echo "  make atlas-generate    - Generate new migration from schema changes"
	@echo "  make atlas-apply       - Apply pending migrations"
	@echo "  make atlas-status      - Check migration status"
	@echo "  make atlas-validate    - Validate migration files"
	@echo "  make atlas-inspect     - Inspect current database schema"
	@echo "  make atlas-diff        - Show diff between schema and database"
	@echo "  make atlas-clean       - Clean migration history (dev only)"
	@echo ""
	@echo "Environment: Use ATLAS_ENV=prod for production (default: local)"

atlas-generate: ## Generate migration from schema changes
	@echo "Generating migration from schema changes..."
	@atlas migrate diff --env $${ATLAS_ENV:-local}
	@echo "✅ Migration generated in migrations/ directory"

atlas-apply: ## Apply pending migrations
	@echo "Applying migrations..."
	@atlas migrate apply --env $${ATLAS_ENV:-local}

atlas-status: ## Check migration status
	@atlas migrate status --env $${ATLAS_ENV:-local}

atlas-validate: ## Validate migration files
	@echo "Validating migrations..."
	@atlas migrate validate --env $${ATLAS_ENV:-local}

atlas-inspect: ## Inspect current database schema
	@atlas schema inspect --env $${ATLAS_ENV:-local}

atlas-diff: ## Show diff between schema and database
	@atlas schema diff --from file://internal/domain --to env://local

atlas-hash: ## Rehash migration directory
	@echo "Rehashing migration directory..."
	@atlas migrate hash --env $${ATLAS_ENV:-local}

atlas-clean: ## Clean migration history (DEV ONLY - destructive!)
	@echo "⚠️  This will drop all tables and migration history!"
	@read -p "Are you sure? Type 'yes' to continue: " confirm; \
	if [ "$$confirm" = "yes" ]; then \
		atlas schema clean --env local --auto-approve; \
		echo "✅ Database cleaned"; \
	else \
		echo "Cancelled"; \
	fi

atlas-new: ## Create a new empty migration file
	@read -p "Enter migration name: " name; \
	atlas migrate new --env local $$name
