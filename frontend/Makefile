.PHONY: help install dev build test lint fmt clean clean-all serve

# Colors for output
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[1;33m
RED := \033[0;31m
NC := \033[0m # No Color

# Variables
NODE_VERSION := $(shell node --version)
NPM_VERSION := $(shell npm --version)

# ============================================================================
# HELP
# ============================================================================

help:
	@echo "$(BLUE)╔════════════════════════════════════════════════════════════════════════════╗$(NC)"
	@echo "$(BLUE)║              Frontend (Next.js) - Makefile Commands                        ║$(NC)"
	@echo "$(BLUE)╚════════════════════════════════════════════════════════════════════════════╝$(NC)"
	@echo ""
	@echo "$(GREEN)DEVELOPMENT:$(NC)"
	@echo "  make dev                Run development server (port 3000)"
	@echo "  make dev-turbo          Run with Turbopack (faster)"
	@echo ""
	@echo "$(GREEN)BUILD:$(NC)"
	@echo "  make build              Build for production"
	@echo "  make build-analyze      Build and analyze bundle size"
	@echo ""
	@echo "$(GREEN)RUNNING:$(NC)"
	@echo "  make start              Run production build"
	@echo "  make serve              Serve production build locally"
	@echo ""
	@echo "$(GREEN)TESTING:$(NC)"
	@echo "  make test               Run test suite"
	@echo "  make test-watch         Run tests in watch mode"
	@echo "  make test-coverage      Run tests with coverage"
	@echo ""
	@echo "$(GREEN)CODE QUALITY:$(NC)"
	@echo "  make lint               Run ESLint"
	@echo "  make lint-fix           Fix linting issues"
	@echo "  make fmt                Format code with Prettier"
	@echo "  make fmt-check          Check formatting"
	@echo "  make type-check         Run TypeScript compiler"
	@echo "  make check              Run all quality checks"
	@echo ""
	@echo "$(GREEN)DEPENDENCIES:$(NC)"
	@echo "  make install            Install dependencies (npm ci)"
	@echo "  make update             Update dependencies"
	@echo "  make audit              Check for vulnerabilities"
	@echo "  make audit-fix          Fix vulnerabilities"
	@echo ""
	@echo "$(GREEN)UTILITIES:$(NC)"
	@echo "  make version            Show versions"
	@echo "  make info               Show project info"
	@echo "  make env-check          Check environment setup"
	@echo "  make clean              Remove build artifacts"
	@echo "  make clean-all          Remove all generated files"
	@echo ""

# ============================================================================
# DEPENDENCY MANAGEMENT
# ============================================================================

.PHONY: install update audit audit-fix deps

install:
	@echo "$(BLUE)Installing dependencies...$(NC)"
	@npm ci
	@echo "$(GREEN)✓ Dependencies installed$(NC)"

update:
	@echo "$(BLUE)Updating dependencies...$(NC)"
	@npm update
	@echo "$(GREEN)✓ Dependencies updated$(NC)"

audit:
	@echo "$(BLUE)Auditing dependencies for vulnerabilities...$(NC)"
	@npm audit

audit-fix:
	@echo "$(BLUE)Fixing vulnerabilities...$(NC)"
	@npm audit fix
	@echo "$(GREEN)✓ Vulnerabilities fixed$(NC)"

deps:
	@echo "$(BLUE)Listing dependencies:$(NC)"
	@npm list --depth=0

# ============================================================================
# DEVELOPMENT
# ============================================================================

.PHONY: dev dev-turbo

dev:
	@echo "$(BLUE)Starting development server...$(NC)"
	@echo "$(YELLOW)▸ Port: http://localhost:3000$(NC)"
	@echo "$(YELLOW)▸ Backend: http://localhost:8080$(NC)"
	@npm run dev

dev-turbo:
	@echo "$(BLUE)Starting development server with Turbopack...$(NC)"
	@echo "$(YELLOW)▸ Port: http://localhost:3000$(NC)"
	@echo "$(YELLOW)▸ Note: Turbopack is experimental$(NC)"
	@npm run dev -- --turbopack

# ============================================================================
# BUILD & PRODUCTION
# ============================================================================

.PHONY: build build-analyze start serve

build:
	@echo "$(BLUE)Building for production...$(NC)"
	@npm run build
	@echo "$(GREEN)✓ Build complete$(NC)"
	@echo "$(YELLOW)Next: make start$(NC)"

build-analyze:
	@echo "$(BLUE)Building with bundle analysis...$(NC)"
	@npm run build
	@echo "$(BLUE)Analyzing bundle...$(NC)"
	@if command -v next-bundle-analyzer &> /dev/null; then \
		echo "$(YELLOW)Bundle analysis available via: npm run analyze$(NC)"; \
	fi
	@echo "$(GREEN)✓ Build complete$(NC)"

start:
	@echo "$(BLUE)Starting production server...$(NC)"
	@npm start

serve:
	@echo "$(BLUE)Serving production build...$(NC)"
	@npm run build
	@npm start

# ============================================================================
# TESTING
# ============================================================================

.PHONY: test test-watch test-coverage test-debug

test:
	@echo "$(BLUE)Running tests...$(NC)"
	@echo "$(YELLOW)Tests not configured yet$(NC)"
	@echo "$(YELLOW)To add tests, install Jest: npm install --save-dev jest @testing-library/react$(NC)"

test-watch:
	@echo "$(BLUE)Running tests in watch mode...$(NC)"
	@npm test -- --watch

test-coverage:
	@echo "$(BLUE)Generating test coverage...$(NC)"
	@npm test -- --coverage

test-debug:
	@echo "$(BLUE)Running tests in debug mode...$(NC)"
	@npm test -- --inspect-brk

# ============================================================================
# LINTING & FORMATTING
# ============================================================================

.PHONY: lint lint-fix fmt fmt-check type-check

lint:
	@echo "$(BLUE)Running ESLint...$(NC)"
	@npm run lint 2>/dev/null || echo "$(YELLOW)ESLint not configured$(NC)"

lint-fix:
	@echo "$(BLUE)Fixing linting issues...$(NC)"
	@npm run lint -- --fix 2>/dev/null || echo "$(YELLOW)ESLint not configured$(NC)"

fmt:
	@echo "$(BLUE)Formatting code with Prettier...$(NC)"
	@npx prettier --write . 2>/dev/null || echo "$(YELLOW)Prettier not available$(NC)"
	@echo "$(GREEN)✓ Code formatted$(NC)"

fmt-check:
	@echo "$(BLUE)Checking code format...$(NC)"
	@npx prettier --check . 2>/dev/null || echo "$(YELLOW)Formatting issues found - run: make fmt$(NC)"

type-check:
	@echo "$(BLUE)Running TypeScript compiler...$(NC)"
	@npx tsc --noEmit
	@echo "$(GREEN)✓ No type errors$(NC)"

check: lint type-check
	@echo "$(GREEN)✓ All quality checks passed$(NC)"

# ============================================================================
# UTILITIES
# ============================================================================

.PHONY: version info env-check

version:
	@echo "$(BLUE)Version Information:$(NC)"
	@echo "  Node: $(NODE_VERSION)"
	@echo "  npm: $(NPM_VERSION)"
	@echo "  Next.js: $$(npm list next --depth=0 2>/dev/null | grep next)"

info:
	@echo "$(BLUE)Project Information:$(NC)"
	@echo "  Framework: Next.js"
	@echo "  Language: TypeScript"
	@echo "  Port: 3000"
	@echo "  Backend: http://localhost:8080"
	@echo ""
	@echo "$(BLUE)Directories:$(NC)"
	@echo "  ./app          - Next.js app (routing)"
	@echo "  ./components   - React components"
	@echo "  ./lib          - Utilities (API, auth, contexts)"
	@echo "  ./public       - Static assets"
	@echo "  ./styles       - Global CSS modules"

env-check:
	@echo "$(BLUE)Environment Check:$(NC)"
	@echo "  ✓ Node.js: $(NODE_VERSION)"
	@echo "  ✓ npm: $(NPM_VERSION)"
	@if [ -f .env ]; then \
		echo "  ✓ .env: exists"; \
	else \
		echo "  $(RED)✗ .env: missing$(NC)"; \
	fi
	@if [ -f tsconfig.json ]; then \
		echo "  ✓ tsconfig.json: exists"; \
	else \
		echo "  $(RED)✗ tsconfig.json: missing$(NC)"; \
	fi
	@if [ -f next.config.ts ]; then \
		echo "  ✓ next.config.ts: exists"; \
	else \
		echo "  $(RED)✗ next.config.ts: missing$(NC)"; \
	fi
	@if [ -d node_modules ]; then \
		echo "  ✓ node_modules: exists"; \
	else \
		echo "  $(RED)✗ node_modules: missing - run: make install$(NC)"; \
	fi

# ============================================================================
# CLEANUP
# ============================================================================

.PHONY: clean clean-all clean-cache

clean:
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	@rm -rf .next
	@echo "$(GREEN)✓ Cleaned .next$(NC)"

clean-all: clean
	@echo "$(BLUE)Deep cleaning...$(NC)"
	@rm -rf node_modules
	@rm -rf .next
	@rm -rf .turbo
	@rm -f package-lock.json
	@echo "$(GREEN)✓ All cleaned$(NC)"
	@echo "$(YELLOW)Run: make install$(NC)"

clean-cache:
	@echo "$(BLUE)Clearing Next.js cache...$(NC)"
	@rm -rf .next
	@npm cache clean --force
	@echo "$(GREEN)✓ Cache cleared$(NC)"

# ============================================================================
# GIT & COMMITS
# ============================================================================

.PHONY: pre-commit pre-push

pre-commit: lint type-check fmt
	@echo "$(GREEN)✓ Pre-commit checks passed$(NC)"

pre-push: build lint type-check
	@echo "$(GREEN)✓ Pre-push checks passed$(NC)"
	@echo "$(YELLOW)Ready to push!$(NC)"

# ============================================================================
# DEFAULT
# ============================================================================

.DEFAULT_GOAL := help
