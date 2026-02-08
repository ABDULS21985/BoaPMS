# ============================================================
# PMS API — Makefile
# ============================================================

APP_NAME    := pms-api
MAIN        := ./cmd/api
BUILD_DIR   := bin
GO          := go
GOFLAGS     := -trimpath -ldflags="-s -w"
DOCKER_IMG  := enterprise-pms/pms-api

# Database migration tool (golang-migrate)
MIGRATE     := migrate
PG_DSN      ?= "postgres://postgres:postgres@localhost:5432/pms_db?sslmode=disable"
MIGRATIONS  := ./migrations

.PHONY: all build run test lint fmt vet clean \
        migrate-up migrate-down migrate-create \
        swagger-gen docker-build docker-up docker-down \
        dev deps

## —— Default ——————————————————————————————————————————————
all: build

## —— Build ———————————————————————————————————————————————
build: ## Build the API binary
	@echo "==> Building $(APP_NAME)..."
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN)

## —— Run —————————————————————————————————————————————————
run: build ## Build and run the API
	@echo "==> Running $(APP_NAME)..."
	APP_ENV=development $(BUILD_DIR)/$(APP_NAME)

## —— Development —————————————————————————————————————————
dev: ## Run with hot-reload via air
	@command -v air >/dev/null 2>&1 || { echo "Install air: go install github.com/air-verse/air@latest"; exit 1; }
	APP_ENV=development air

## —— Test ————————————————————————————————————————————————
test: ## Run all tests
	$(GO) test -v -race -count=1 ./...

test-cover: ## Run tests with coverage report
	$(GO) test -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## —— Code Quality ————————————————————————————————————————
lint: ## Run golangci-lint
	@command -v golangci-lint >/dev/null 2>&1 || { echo "Install golangci-lint: https://golangci-lint.run/usage/install/"; exit 1; }
	golangci-lint run ./...

fmt: ## Format code
	$(GO) fmt ./...
	goimports -w .

vet: ## Run go vet
	$(GO) vet ./...

## —— Database Migrations —————————————————————————————————
migrate-up: ## Run all pending migrations
	$(MIGRATE) -path $(MIGRATIONS) -database $(PG_DSN) up

migrate-down: ## Rollback the last migration
	$(MIGRATE) -path $(MIGRATIONS) -database $(PG_DSN) down 1

migrate-create: ## Create a new migration (usage: make migrate-create name=create_users)
	@test -n "$(name)" || { echo "Usage: make migrate-create name=<migration_name>"; exit 1; }
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS) -seq $(name)

migrate-status: ## Show migration status
	$(MIGRATE) -path $(MIGRATIONS) -database $(PG_DSN) version

## —— Swagger —————————————————————————————————————————————
swagger-gen: ## Generate Swagger/OpenAPI docs
	@command -v swag >/dev/null 2>&1 || { echo "Install swag: go install github.com/swaggo/swag/cmd/swag@latest"; exit 1; }
	swag init -g $(MAIN)/main.go -o ./docs

## —— Docker ——————————————————————————————————————————————
docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMG):latest .

docker-up: ## Start all services via docker-compose
	docker compose up -d

docker-down: ## Stop all services via docker-compose
	docker compose down

docker-logs: ## Tail docker-compose logs
	docker compose logs -f

## —— Dependencies ————————————————————————————————————————
deps: ## Download and tidy Go module dependencies
	$(GO) mod download
	$(GO) mod tidy

## —— Clean ———————————————————————————————————————————————
clean: ## Remove build artifacts
	rm -rf $(BUILD_DIR) coverage.out coverage.html docs/

## —— Help ————————————————————————————————————————————————
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
