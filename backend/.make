# Task file for Recession Recovery Workspace.
# Use with: make -f .make <target>

.DEFAULT_GOAL := help
SHELL := cmd.exe
.SHELLFLAGS := /C

BIN := workspace-app
CMD := ./cmd/workspace-app
BUILD_DIR := bin
COVERAGE_DIR := coverage
GO := go
VERSION ?= dev
COMMIT ?= unknown
BUILD_TIME ?= unknown
LDFLAGS := -ldflags "-X 'main.Version=$(VERSION)' -X 'main.Commit=$(COMMIT)' -X 'main.BuildTime=$(BUILD_TIME)'"

.PHONY: help \
	build build-dev build-release build-linux build-macos build-windows \
	run run-dev run-migrate \
	test test-verbose test-auth test-common test-user test-profile test-jobs test-ideas test-coverage \
	coverage coverage-html coverage-report \
	fmt vet lint tidy deps check pre-commit \
	migrate migrate-status migrate-reset \
	docs init-env version \
	clean clean-coverage clean-all

help:
	@echo "Available targets:"
	@echo ""
	@echo "Build:"
	@echo "  make build             Build release binary"
	@echo "  make build-dev         Build debug binary"
	@echo "  make build-release     Build optimized release binary"
	@echo "  make build-linux       Build Linux amd64 binary"
	@echo "  make build-macos       Build macOS amd64 binary"
	@echo "  make build-windows     Build Windows amd64 binary"
	@echo ""
	@echo "Run:"
	@echo "  make run               Build and run the app"
	@echo "  make run-dev           Run with go run"
	@echo "  make run-migrate       Run migrations through app startup"
	@echo ""
	@echo "Test:"
	@echo "  make test              Run all tests"
	@echo "  make test-verbose      Run all tests with verbose output"
	@echo "  make test-auth         Run auth tests"
	@echo "  make test-common       Run common package tests"
	@echo "  make test-user         Run user domain tests"
	@echo "  make test-profile      Run profile domain tests"
	@echo "  make test-jobs         Run jobs domain tests"
	@echo "  make test-ideas        Run ideas domain tests"
	@echo "  make test-coverage     Run tests with coverage"
	@echo "  make coverage          Print coverage summary"
	@echo "  make coverage-html     Generate HTML coverage report"
	@echo "  make coverage-report   Show existing coverage total"
	@echo ""
	@echo "Quality:"
	@echo "  make fmt               Format Go code"
	@echo "  make vet               Run go vet"
	@echo "  make lint              Run golangci-lint if installed"
	@echo "  make tidy              Run go mod tidy"
	@echo "  make deps              Verify module dependencies"
	@echo "  make check             Run fmt, vet, tidy, and tests"
	@echo "  make pre-commit        Run commit readiness checks"
	@echo ""
	@echo "Database:"
	@echo "  make migrate           Run app startup migrations"
	@echo "  make migrate-status    Print migration status command"
	@echo "  make migrate-reset     Print database reset guidance"
	@echo ""
	@echo "Utility:"
	@echo "  make docs              Show docs URLs/files"
	@echo "  make init-env          Create .env from .env.example"
	@echo "  make version           Show build metadata"
	@echo "  make clean             Remove build artifacts"
	@echo "  make clean-coverage    Remove coverage artifacts"
	@echo "  make clean-all         Remove all generated artifacts"

$(BUILD_DIR):
	@if not exist "$(BUILD_DIR)" mkdir "$(BUILD_DIR)"

build: $(BUILD_DIR)
	@echo "Building $(BUILD_DIR)/$(BIN)..."
	@$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BIN) $(CMD)

build-dev: $(BUILD_DIR)
	@echo "Building debug binary..."
	@$(GO) build -o $(BUILD_DIR)/$(BIN)-debug $(CMD)

build-release: $(BUILD_DIR)
	@echo "Building optimized release binary..."
	@$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BIN)-release $(CMD)

build-linux: $(BUILD_DIR)
	@echo "Building Linux amd64 binary..."
	@set GOOS=linux&& set GOARCH=amd64&& $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BIN)-linux-amd64 $(CMD)

build-macos: $(BUILD_DIR)
	@echo "Building macOS amd64 binary..."
	@set GOOS=darwin&& set GOARCH=amd64&& $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BIN)-macos-amd64 $(CMD)

build-windows: $(BUILD_DIR)
	@echo "Building Windows amd64 binary..."
	@set GOOS=windows&& set GOARCH=amd64&& $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BIN)-windows-amd64.exe $(CMD)

run: build
	@echo "Running $(BUILD_DIR)/$(BIN)..."
	@./$(BUILD_DIR)/$(BIN)

run-dev:
	@$(GO) run $(CMD)

run-migrate: migrate

test:
	@$(GO) test ./...

test-verbose:
	@$(GO) test -v ./...

test-auth:
	@$(GO) test -v ./internal/auth

test-common:
	@$(GO) test -v ./internal/common

test-user:
	@$(GO) test -v ./internal/user/...

test-profile:
	@$(GO) test -v ./internal/profile/...

test-jobs:
	@$(GO) test -v ./internal/jobs/...

test-ideas:
	@$(GO) test -v ./internal/ideas/...

test-coverage:
	@if not exist "$(COVERAGE_DIR)" mkdir "$(COVERAGE_DIR)"
	@$(GO) test -coverprofile=$(COVERAGE_DIR)/coverage.out -covermode=atomic ./...
	@$(GO) tool cover -func=$(COVERAGE_DIR)/coverage.out

coverage:
	@if not exist "$(COVERAGE_DIR)" mkdir "$(COVERAGE_DIR)"
	@$(GO) test -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	@$(GO) tool cover -func=$(COVERAGE_DIR)/coverage.out

coverage-html:
	@if not exist "$(COVERAGE_DIR)" mkdir "$(COVERAGE_DIR)"
	@$(GO) test -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	@$(GO) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "Generated $(COVERAGE_DIR)/coverage.html"

coverage-report:
	@if exist "$(COVERAGE_DIR)\coverage.out" ($(GO) tool cover -func=$(COVERAGE_DIR)/coverage.out | findstr total) else (echo No coverage file found. Run: make test-coverage)

fmt:
	@$(GO) fmt ./...

vet:
	@$(GO) vet ./...

lint:
	@where golangci-lint >NUL 2>NUL && golangci-lint run ./... || echo golangci-lint is not installed

tidy:
	@$(GO) mod tidy

deps:
	@$(GO) mod verify

check: fmt vet tidy test

pre-commit: fmt vet test

migrate:
	@$(GO) run $(CMD)

migrate-status:
	@echo "Migrations are tracked in the migrations table."
	@echo "mysql -u root -p my_site -e \"SELECT * FROM migrations ORDER BY executed_at;\""

migrate-reset:
	@echo "Reset manually if needed:"
	@echo "DROP DATABASE IF EXISTS my_site;"
	@echo "CREATE DATABASE my_site CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
	@echo "Then run: make -f .make migrate"

docs:
	@echo "Swagger UI: http://localhost:8080/swagger-ui/"
	@echo "OpenAPI:    http://localhost:8080/openapi.yaml"
	@echo "Testing:    TESTING.md"

init-env:
	@if not exist ".env" (copy ".env.example" ".env" >NUL && echo Created .env from .env.example) else (echo .env already exists)

version:
	@echo "Version:    $(VERSION)"
	@echo "Commit:     $(COMMIT)"
	@echo "Build time: $(BUILD_TIME)"
	@$(GO) version

clean:
	@if exist "$(BUILD_DIR)" rmdir /S /Q "$(BUILD_DIR)"

clean-coverage:
	@if exist "$(COVERAGE_DIR)" rmdir /S /Q "$(COVERAGE_DIR)"
	@if exist "coverage.out" del /Q "coverage.out"
	@if exist "coverage.html" del /Q "coverage.html"

clean-all: clean clean-coverage
	@$(GO) clean
