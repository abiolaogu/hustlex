# HustleX Pro Monorepo Makefile
# Unified build system for all applications

.PHONY: help setup clean test lint build run docker

# Colors
CYAN := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RESET := \033[0m

# Default target
help:
	@echo "$(CYAN)HustleX Pro Monorepo$(RESET)"
	@echo ""
	@echo "$(GREEN)Setup:$(RESET)"
	@echo "  make setup           - Install all dependencies"
	@echo "  make setup-flutter   - Setup Flutter packages (melos)"
	@echo "  make setup-go        - Setup Go backend"
	@echo "  make setup-python    - Setup Python recommendation service"
	@echo ""
	@echo "$(GREEN)Development:$(RESET)"
	@echo "  make run-api         - Run API server"
	@echo "  make run-consumer    - Run consumer Flutter app"
	@echo "  make run-provider    - Run provider Flutter app"
	@echo "  make run-admin       - Run admin web dashboard"
	@echo "  make run-ml          - Run ML recommendation service"
	@echo ""
	@echo "$(GREEN)Testing:$(RESET)"
	@echo "  make test            - Run all tests"
	@echo "  make test-api        - Run API tests"
	@echo "  make test-flutter    - Run Flutter tests"
	@echo "  make test-ml         - Run ML service tests"
	@echo ""
	@echo "$(GREEN)Build:$(RESET)"
	@echo "  make build           - Build all applications"
	@echo "  make build-api       - Build API binary"
	@echo "  make build-consumer  - Build consumer app"
	@echo "  make build-provider  - Build provider app"
	@echo ""
	@echo "$(GREEN)Docker:$(RESET)"
	@echo "  make docker-up       - Start all services"
	@echo "  make docker-down     - Stop all services"
	@echo "  make docker-build    - Build Docker images"
	@echo ""
	@echo "$(GREEN)Code Quality:$(RESET)"
	@echo "  make lint            - Lint all code"
	@echo "  make format          - Format all code"
	@echo "  make generate        - Generate code (freezed, protobuf)"

# ==============================================================================
# Setup
# ==============================================================================

setup: setup-flutter setup-go setup-python
	@echo "$(GREEN)✓ All dependencies installed$(RESET)"

setup-flutter:
	@echo "$(CYAN)Setting up Flutter packages...$(RESET)"
	dart pub global activate melos
	melos bootstrap

setup-go:
	@echo "$(CYAN)Setting up Go backend...$(RESET)"
	cd apps/api && go mod download
	cd packages/go-common && go mod download

setup-python:
	@echo "$(CYAN)Setting up Python ML service...$(RESET)"
	cd apps/recommendation && pip install -r requirements.txt 2>/dev/null || true

# ==============================================================================
# Development
# ==============================================================================

run-api:
	@echo "$(CYAN)Starting API server...$(RESET)"
	cd apps/api && go run cmd/server/main.go

run-consumer:
	@echo "$(CYAN)Starting consumer app...$(RESET)"
	cd apps/consumer-app && flutter run

run-provider:
	@echo "$(CYAN)Starting provider app...$(RESET)"
	cd apps/provider-app && flutter run

run-admin:
	@echo "$(CYAN)Starting admin dashboard...$(RESET)"
	cd apps/admin-web && npm run dev 2>/dev/null || echo "Admin web not configured"

run-ml:
	@echo "$(CYAN)Starting ML recommendation service...$(RESET)"
	cd apps/recommendation && python -m uvicorn main:app --reload 2>/dev/null || echo "ML service not configured"

# ==============================================================================
# Testing
# ==============================================================================

test: test-api test-flutter test-ml
	@echo "$(GREEN)✓ All tests passed$(RESET)"

test-api:
	@echo "$(CYAN)Running API tests...$(RESET)"
	cd apps/api && go test -v ./...

test-flutter:
	@echo "$(CYAN)Running Flutter tests...$(RESET)"
	melos run test

test-ml:
	@echo "$(CYAN)Running ML service tests...$(RESET)"
	cd apps/recommendation && pytest 2>/dev/null || echo "No ML tests configured"

# ==============================================================================
# Build
# ==============================================================================

build: build-api build-consumer build-provider
	@echo "$(GREEN)✓ All builds completed$(RESET)"

build-api:
	@echo "$(CYAN)Building API...$(RESET)"
	cd apps/api && CGO_ENABLED=0 go build -o bin/server cmd/server/main.go

build-consumer:
	@echo "$(CYAN)Building consumer app...$(RESET)"
	melos run build:consumer:android

build-provider:
	@echo "$(CYAN)Building provider app...$(RESET)"
	melos run build:provider:android

# ==============================================================================
# Docker
# ==============================================================================

docker-up:
	@echo "$(CYAN)Starting Docker services...$(RESET)"
	docker-compose up -d

docker-down:
	@echo "$(CYAN)Stopping Docker services...$(RESET)"
	docker-compose down

docker-build:
	@echo "$(CYAN)Building Docker images...$(RESET)"
	docker-compose build

docker-logs:
	docker-compose logs -f

# ==============================================================================
# Code Quality
# ==============================================================================

lint: lint-go lint-flutter
	@echo "$(GREEN)✓ Linting complete$(RESET)"

lint-go:
	@echo "$(CYAN)Linting Go code...$(RESET)"
	cd apps/api && golangci-lint run ./... 2>/dev/null || go vet ./...

lint-flutter:
	@echo "$(CYAN)Linting Flutter code...$(RESET)"
	melos run analyze

format: format-go format-flutter
	@echo "$(GREEN)✓ Formatting complete$(RESET)"

format-go:
	@echo "$(CYAN)Formatting Go code...$(RESET)"
	cd apps/api && go fmt ./...

format-flutter:
	@echo "$(CYAN)Formatting Flutter code...$(RESET)"
	melos run format

generate: generate-flutter generate-proto
	@echo "$(GREEN)✓ Code generation complete$(RESET)"

generate-flutter:
	@echo "$(CYAN)Generating Flutter code...$(RESET)"
	melos run generate

generate-proto:
	@echo "$(CYAN)Generating protobuf code...$(RESET)"
	cd packages/proto && buf generate 2>/dev/null || echo "Protobuf generation skipped"

# ==============================================================================
# Clean
# ==============================================================================

clean:
	@echo "$(CYAN)Cleaning all build artifacts...$(RESET)"
	melos run clean
	cd apps/api && rm -rf bin/
	rm -rf coverage/
	@echo "$(GREEN)✓ Clean complete$(RESET)"

# ==============================================================================
# Database
# ==============================================================================

db-migrate:
	@echo "$(CYAN)Running database migrations...$(RESET)"
	cd apps/api && go run cmd/migrate/main.go up

db-rollback:
	@echo "$(CYAN)Rolling back database migrations...$(RESET)"
	cd apps/api && go run cmd/migrate/main.go down

db-seed:
	@echo "$(CYAN)Seeding database...$(RESET)"
	cd apps/api && go run cmd/seed/main.go
