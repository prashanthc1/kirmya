# Makefile for Recession Recovery Workspace
# A modular monolith platform for connecting job seekers, recruiters, founders, freelancers, mentors during recession

.DEFAULT_GOAL := help
ifneq ($(OS),Windows_NT)
SHELL := /bin/bash
endif

# Variables
BIN := workspace-app
CMD := ./cmd/workspace-app
BUILD_DIR := bin
COVERAGE_DIR := coverage
GO := go
GOFLAGS := -v
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d %H:%M:%S' 2>/dev/null || echo "unknown")
LDFLAGS := -ldflags "-X 'main.Version=$(VERSION)' -X 'main.Commit=$(COMMIT)' -X 'main.BuildTime=$(BUILD_TIME)'"

# Check if air is installed and available in PATH
AIR := $(shell where air 2>nul || which air 2>/dev/null)

.PHONY: help build build-release build-dev run run-dev test test-verbose test-jobs test-ideas test-coverage \
        fmt lint vet tidy deps migrate clean clean-coverage clean-all coverage coverage-html coverage-report \
        coverage-gate docker-build docker-run docker-clean install-tools pre-commit check

# ============================================================================
# HELP
# ============================================================================

help:
	@echo "╔════════════════════════════════════════════════════════════════════════════╗"
	@echo "║     Recession Recovery Workspace - Makefile Commands                       ║"
	@echo "╚════════════════════════════════════════════════════════════════════════════╝"
	@echo ""
	@echo "BUILD TARGETS:"
	@echo "  make build              Build release binary with version info"
	@echo "  make build-dev          Build debug binary (faster, larger)"
	@echo "  make build-release      Build optimized release binary"
	@echo "  make build-linux        Build for Linux (x86_64)"
	@echo "  make build-macos        Build for macOS (x86_64)"
	@echo "  make build-windows      Build for Windows (x86_64)"
	@echo ""
	@echo "RUN TARGETS:"
	@echo "  make run                Run application (requires .env)"
	@echo "  make run-dev            Run with watch mode (requires CompileDaemon)"
	@echo "  make run-migrate        Run migrations then start app"
	@echo ""
	@echo "TEST TARGETS:"
	@echo "  make test               Run all unit tests"
	@echo "  make test-verbose       Run tests with verbose output"
	@echo "  make test-jobs          Run only jobs domain tests"
	@echo "  make test-ideas         Run only ideas domain tests"
	@echo "  make test-user          Run only user domain tests"
	@echo "  make test-coverage      Run tests and generate coverage report"
	@echo "  make coverage           Generate coverage report (text)"
	@echo "  make coverage-html      Generate HTML coverage report"
	@echo "  make coverage-report    Display coverage summary"
	@echo "  make coverage-gate      Fail if application+domain coverage < 70%"
	@echo ""
	@echo "CODE QUALITY TARGETS:"
	@echo "  make fmt                Format all Go code"
	@echo "  make vet                Run go vet static analysis"
	@echo "  make lint               Run golangci-lint (requires install)"
	@echo "  make tidy               Clean up go.mod and go.sum"
	@echo "  make deps               Verify dependencies"
	@echo ""
	@echo "DATABASE TARGETS:"
	@echo "  make migrate            Run database migrations"
	@echo "  make migrate-reset      Reset database (drop + recreate)"
	@echo "  make migrate-status     Show migration status"
	@echo ""
	@echo "DOCKER TARGETS:"
	@echo "  make docker-build       Build Docker image"
	@echo "  make docker-run         Run application in Docker"
	@echo "  make docker-clean       Remove Docker image"
	@echo ""
	@echo "UTILITY TARGETS:"
	@echo "  make install-tools      Install development tools (golangci-lint, etc)"
	@echo "  make pre-commit         Run pre-commit checks (fmt, vet, test)"
	@echo "  make check              Run all quality checks"
	@echo "  make version            Show build version info"
	@echo ""
	@echo "CLEANUP TARGETS:"
	@echo "  make clean              Remove build artifacts"
	@echo "  make clean-coverage     Remove coverage files"
	@echo "  make clean-all          Remove all generated files"
	@echo ""

# ============================================================================
# BUILD TARGETS
# ============================================================================

build: $(BUILD_DIR)
	@echo "🔨 Building $(BIN) ($(VERSION))..."
	@$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BIN) $(CMD)
	@echo "✓ Built: $(BUILD_DIR)/$(BIN)"
	@ls -lh $(BUILD_DIR)/$(BIN)

build-dev: $(BUILD_DIR)
	@echo "🔨 Building debug binary..."
	@$(GO) build -o $(BUILD_DIR)/$(BIN)-debug $(CMD)
	@echo "✓ Built: $(BUILD_DIR)/$(BIN)-debug"

build-release: $(BUILD_DIR)
	@echo "🔨 Building optimized release binary..."
	@$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BIN)-release $(CMD)
	@echo "✓ Built: $(BUILD_DIR)/$(BIN)-release"

build-linux: $(BUILD_DIR)
	@echo "🔨 Building for Linux x86_64..."
	@GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BIN)-linux-amd64 $(CMD)
	@echo "✓ Built: $(BUILD_DIR)/$(BIN)-linux-amd64"

build-macos: $(BUILD_DIR)
	@echo "🔨 Building for macOS..."
	@GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BIN)-macos-amd64 $(CMD)
	@echo "✓ Built: $(BUILD_DIR)/$(BIN)-macos-amd64"

build-windows: $(BUILD_DIR)
	@echo "🔨 Building for Windows..."
	@GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BIN)-windows-amd64.exe $(CMD)
	@echo "✓ Built: $(BUILD_DIR)/$(BIN)-windows-amd64.exe"

$(BUILD_DIR):
	@mkdir -p $(BUILD_DIR)

# ============================================================================
# RUN TARGETS
# ============================================================================

run: build
	@echo "🚀 Running $(BIN)..."
	@./$(BUILD_DIR)/$(BIN)

run-dev:
	@echo "🚀 Running in development mode (watch enabled via Air)..."
ifneq ($(AIR),)
	air -c .air.toml
else
	@echo "⚠️  Air not installed or not in PATH. Running without watch mode..."
	@echo "💡 To install Air: go install github.com/air-verse/air@latest"
	$(GO) run $(CMD)
endif

run-migrate: migrate run

# ============================================================================
# TEST TARGETS
# ============================================================================

test:
	@echo "🧪 Running all tests..."
	@$(GO) test ./... $(GOFLAGS)
	@echo "✓ Tests passed"

test-verbose:
	@echo "🧪 Running tests with verbose output..."
	@$(GO) test ./... -v

test-jobs:
	@echo "🧪 Running jobs domain tests..."
	@$(GO) test ./internal/jobs/... -v

test-ideas:
	@echo "🧪 Running ideas domain tests..."
	@$(GO) test ./internal/ideas/... -v

test-user:
	@echo "🧪 Running user domain tests..."
	@$(GO) test ./internal/user/... -v

test-coverage:
	@mkdir -p $(COVERAGE_DIR)
	@echo "🧪 Running tests with coverage..."
	@$(GO) test -coverprofile=$(COVERAGE_DIR)/coverage.out -covermode=atomic ./...
	@$(GO) tool cover -func=$(COVERAGE_DIR)/coverage.out | tail -20
	@echo "✓ Coverage data saved to $(COVERAGE_DIR)/coverage.out"

coverage:
	@mkdir -p $(COVERAGE_DIR)
	@echo "📊 Generating coverage report..."
	@$(GO) test -coverprofile=$(COVERAGE_DIR)/coverage.out ./... > /dev/null 2>&1
	@$(GO) tool cover -func=$(COVERAGE_DIR)/coverage.out | tail -30

coverage-html:
	@mkdir -p $(COVERAGE_DIR)
	@echo "📊 Generating HTML coverage report..."
	@$(GO) test -coverprofile=$(COVERAGE_DIR)/coverage.out ./... > /dev/null 2>&1
	@$(GO) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "✓ Generated: $(COVERAGE_DIR)/coverage.html"
	@echo "💡 Open in browser: file://$(shell pwd)/$(COVERAGE_DIR)/coverage.html"

coverage-report:
	@echo "📊 Coverage Summary:"
	@if [ -f "$(COVERAGE_DIR)/coverage.out" ]; then \
		$(GO) tool cover -func=$(COVERAGE_DIR)/coverage.out | grep total; \
	else \
		echo "⚠️  No coverage data. Run 'make test-coverage' first"; \
	fi

coverage-gate:
	@mkdir -p $(COVERAGE_DIR)
	@echo "🚦 Enforcing application+domain coverage gate (>= $${COVERAGE_THRESHOLD:-70}%)..."
	@COVERAGE_DIR=$(COVERAGE_DIR) bash ../scripts/coverage-gate.sh

# ============================================================================
# CODE QUALITY TARGETS
# ============================================================================

fmt:
	@echo "🎨 Formatting code..."
	@$(GO) fmt ./...
	@echo "✓ Code formatted"

vet:
	@echo "🔍 Running go vet..."
	@$(GO) vet ./...
	@echo "✓ No issues found"

lint:
	@echo "🔍 Running golangci-lint..."
	@if command -v golangci-lint &> /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "⚠️  golangci-lint not installed. Install with: make install-tools"; \
	fi

tidy:
	@echo "🧹 Tidying dependencies..."
	@$(GO) mod tidy
	@echo "✓ Dependencies tidied"

deps:
	@echo "🔗 Verifying dependencies..."
	@$(GO) mod verify
	@echo "✓ All dependencies verified"

check: fmt vet tidy test
	@echo "✅ All quality checks passed"

pre-commit: fmt vet test
	@echo "✅ Pre-commit checks passed - ready to commit"

# ============================================================================
# DATABASE TARGETS
# ============================================================================

migrate:
	@echo "🔄 Running migrations..."
	@$(GO) run $(CMD)
	@echo "✓ Migrations completed"

migrate-reset:
	@echo "⚠️  WARNING: This will drop and recreate the PostgreSQL database!"
	@read -p "Continue? (yes/no) " confirm && [ "$${confirm}" = "yes" ] || exit 1
	@echo "🔄 Resetting database..."
	@psql "postgres://$${POSTGRES_USER:-postgres}:$${POSTGRES_PASSWORD:-postgres}@$${POSTGRES_HOST:-localhost}:$${POSTGRES_PORT:-5432}/postgres?sslmode=$${POSTGRES_SSLMODE:-disable}" \
		-c "DROP DATABASE IF EXISTS $${POSTGRES_DB:-kirmya} WITH (FORCE);" \
		-c "CREATE DATABASE $${POSTGRES_DB:-kirmya};"
	@echo "✓ Database reset. Run 'make migrate' (or start the server) to re-apply migrations."

migrate-status:
	@echo "📋 Migration Status:"
	@echo "Migrations are tracked in the 'schema_migrations' table."
	@psql "$${DATABASE_URL:-postgres://postgres:postgres@localhost:5432/kirmya?sslmode=disable}" \
		-c "SELECT * FROM schema_migrations ORDER BY 1;" 2>/dev/null \
		|| echo "Could not connect — set DATABASE_URL or start PostgreSQL first."

# ============================================================================
# DOCKER TARGETS
# ============================================================================

docker-build: build-release
	@echo "🐳 Building Docker image..."
	@docker build -t workspace-app:$(VERSION) -t workspace-app:latest .
	@echo "✓ Docker image built"
	@docker images | grep workspace-app

docker-run: docker-build
	@echo "🐳 Running Docker container..."
	@docker run -p 8080:8080 \
		-e DATABASE_URL=postgres://postgres:postgres@host.docker.internal:5432/kirmya?sslmode=disable \
		-e JWT_SECRET=dev-secret-change-me-in-production \
		-e APP_URL=http://localhost:3000 \
		workspace-app:latest

docker-clean:
	@echo "🐳 Cleaning Docker..."
	@docker rmi -f workspace-app:latest workspace-app:$(VERSION) 2>/dev/null || true
	@echo "✓ Docker cleaned"

# ============================================================================
# UTILITY TARGETS
# ============================================================================

install-tools:
	@echo "🛠️  Installing development tools..."
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@$(GO) install github.com/githubnemo/CompileDaemon@latest
	@echo "✓ Tools installed"
	@echo "  - golangci-lint (linting)"
	@echo "  - CompileDaemon (live reload)"

version:
	@echo "📦 Build Information:"
	@echo "  Version:    $(VERSION)"
	@echo "  Commit:     $(COMMIT)"
	@echo "  Build Time: $(BUILD_TIME)"
	@echo "  Go Version: $$($(GO) version | awk '{print $$3}')"

docs:
	@echo "📚 Generating documentation..."
	@echo "API Documentation: http://localhost:8080/swagger-ui/"
	@echo "Testing Guide: cat TESTING.md"
	@echo "Code Structure: tree -d -L 2 internal/"

init-env:
	@if [ ! -f .env ]; then \
		echo "📝 Creating .env file from .env.example..."; \
		cp .env.example .env; \
		echo "✓ Created .env (update with your settings)"; \
	else \
		echo "ℹ️  .env already exists"; \
	fi

# ============================================================================
# CLEANUP TARGETS
# ============================================================================

clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "✓ Build directory cleaned"

clean-coverage:
	@echo "🧹 Cleaning coverage files..."
	@rm -rf $(COVERAGE_DIR)
	@echo "✓ Coverage files cleaned"

clean-all: clean clean-coverage
	@echo "🧹 Deep cleaning..."
	@rm -f coverage.out coverage.html
	@$(GO) clean
	@echo "✓ All artifacts cleaned"
