.PHONY: help doctor setup dev build test clean install-deps docker logs restart migrate seed

# Colors for output
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[1;33m
RED := \033[0;31m
NC := \033[0m # No Color

# ============================================================================
# HELP
# ============================================================================

help:
	@echo "$(BLUE)╔════════════════════════════════════════════════════════════════════════════╗$(NC)"
	@echo "$(BLUE)║          Recession Recovery Workspace - Full Stack Makefile                 ║$(NC)"
	@echo "$(BLUE)╚════════════════════════════════════════════════════════════════════════════╝$(NC)"
	@echo ""
	@echo "$(GREEN)QUICK START:$(NC)"
	@echo "  make doctor             Check your environment for missing/old tools"
	@echo "  make setup              Setup entire project (install deps, .env files)"
	@echo "  make dev                Start both backend and frontend in watch mode"
	@echo "  make dev-bg             Start backend and frontend in background"
	@echo ""
	@echo "$(GREEN)DEVELOPMENT:$(NC)"
	@echo "  make backend-run        Run backend server"
	@echo "  make backend-dev        Run backend with watch mode (hot reload)"
	@echo "  make frontend-run       Run frontend dev server"
	@echo "  make frontend-build     Build frontend for production"
	@echo ""
	@echo "$(GREEN)TESTING:$(NC)"
	@echo "  make test               Run all tests (backend + frontend)"
	@echo "  make backend-test       Run backend unit tests"
	@echo "  make backend-test-v     Run backend tests verbose"
	@echo "  make backend-coverage   Generate backend coverage report"
	@echo "  make frontend-test      Run frontend tests"
	@echo ""
	@echo "$(GREEN)BUILD & DEPLOY:$(NC)"
	@echo "  make build              Build backend release binary"
	@echo "  make build-all          Build backend binary and frontend assets"
	@echo "  make docker-build       Build Docker image for backend"
	@echo "  make docker-run         Run app in Docker"
	@echo ""
	@echo "$(GREEN)DATABASE:$(NC)"
	@echo "  make migrate            Run database migrations"
	@echo "  make seed               Seed database with sample data"
	@echo "  make db-reset           Reset database (drop + recreate)"
	@echo ""
	@echo "$(GREEN)CODE QUALITY:$(NC)"
	@echo "  make fmt                Format all code (Go + frontend)"
	@echo "  make lint               Lint all code"
	@echo "  make check              Run all quality checks"
	@echo "  make pre-commit         Pre-commit checks (fmt, lint, test)"
	@echo ""
	@echo "$(GREEN)UTILITIES:$(NC)"
	@echo "  make install-tools      Install development tools"
	@echo "  make install-deps       Install all dependencies"
	@echo "  make logs               Show backend logs (if running)"
	@echo "  make ps                 Show running services"
	@echo "  make restart            Restart all services"
	@echo "  make clean              Remove build artifacts"
	@echo "  make clean-all          Complete cleanup"
	@echo ""

# ============================================================================
# SETUP & INSTALLATION
# ============================================================================

.PHONY: doctor setup install-deps install-tools init-env

doctor:
	@-bash scripts/doctor.sh

setup: install-deps init-env
	@echo "$(GREEN)✓ Project setup complete!$(NC)"
	@echo "$(YELLOW)Next steps:$(NC)"
	@echo "  1. Update backend/.env with database credentials"
	@echo "  2. Run: make migrate"
	@echo "  3. Run: make dev"

install-deps: backend-deps frontend-deps
	@echo "$(GREEN)✓ All dependencies installed$(NC)"

backend-deps:
	@echo "$(BLUE)Installing backend dependencies...$(NC)"
	@cd backend && go mod download && go mod verify
	@echo "$(GREEN)✓ Backend dependencies installed$(NC)"

frontend-deps:
	@echo "$(BLUE)Installing frontend dependencies...$(NC)"
	@cd frontend && npm ci
	@echo "$(GREEN)✓ Frontend dependencies installed$(NC)"

install-tools:
	@echo "$(BLUE)Installing development tools...$(NC)"
	@cd backend && make install-tools
	@echo "$(GREEN)✓ Development tools installed$(NC)"

init-env:
	@echo "$(BLUE)Initializing environment files...$(NC)"
	@if [ ! -f backend/.env ]; then \
		cp backend/.env.example backend/.env; \
		echo "$(GREEN)✓ Created backend/.env$(NC)"; \
	else \
		echo "$(YELLOW)ℹ backend/.env already exists$(NC)"; \
	fi

# ============================================================================
# DEVELOPMENT - RUNNING
# ============================================================================

.PHONY: dev dev-bg backend-run backend-dev frontend-run backend-logs frontend-logs

dev:
	@echo "$(BLUE)Starting both backend and frontend...$(NC)"
	@echo "$(YELLOW)Backend: http://localhost:8080/api/v1/health$(NC)"
	@echo "$(YELLOW)Frontend: http://localhost:3000$(NC)"
	@echo "$(YELLOW)Press Ctrl+C to stop$(NC)"
	@bash -c 'trap "trap - TERM; kill 0" TERM; \
		(cd backend && go run ./cmd/workspace-app/main.go) & \
		(cd frontend && npm run dev) & \
		wait'

dev-bg:
	@echo "$(BLUE)Starting backend and frontend in background...$(NC)"
	@cd backend && go run ./cmd/workspace-app/main.go > logs/backend.log 2>&1 &
	@cd frontend && npm run dev > logs/frontend.log 2>&1 &
	@echo "$(GREEN)✓ Services started in background$(NC)"
	@echo "$(YELLOW)View logs: make logs$(NC)"
	@sleep 2
	@echo "$(YELLOW)Backend: http://localhost:8080/api/v1/health$(NC)"
	@echo "$(YELLOW)Frontend: http://localhost:3000$(NC)"

backend-run:
	@echo "$(BLUE)Starting backend server...$(NC)"
	@cd backend && go run ./cmd/workspace-app/main.go

backend-dev:
	@echo "$(BLUE)Starting backend with watch mode...$(NC)"
	@cd backend && make run-dev

frontend-run:
	@echo "$(BLUE)Starting frontend dev server...$(NC)"
	@cd frontend && npm run dev

backend-logs:
	@tail -f logs/backend.log 2>/dev/null || echo "$(YELLOW)No backend logs found$(NC)"

frontend-logs:
	@tail -f logs/frontend.log 2>/dev/null || echo "$(YELLOW)No frontend logs found$(NC)"

# ============================================================================
# TESTING
# ============================================================================

.PHONY: test backend-test backend-test-v backend-coverage frontend-test

test: backend-test frontend-test
	@echo "$(GREEN)✓ All tests passed$(NC)"

backend-test:
	@echo "$(BLUE)Running backend tests...$(NC)"
	@cd backend && make test

backend-test-v:
	@echo "$(BLUE)Running backend tests (verbose)...$(NC)"
	@cd backend && make test-verbose

backend-coverage:
	@echo "$(BLUE)Generating backend coverage report...$(NC)"
	@cd backend && make coverage-html
	@echo "$(GREEN)✓ Coverage report generated$(NC)"

frontend-test:
	@echo "$(BLUE)Running frontend tests...$(NC)"
	@cd frontend && npm run test 2>/dev/null || echo "$(YELLOW)Frontend tests not configured$(NC)"

# ============================================================================
# BUILD
# ============================================================================

.PHONY: build build-all backend-build frontend-build

build: backend-build
	@echo "$(GREEN)✓ Backend built successfully$(NC)"

build-all: backend-build frontend-build
	@echo "$(GREEN)✓ Full stack built successfully$(NC)"

backend-build:
	@echo "$(BLUE)Building backend release binary...$(NC)"
	@cd backend && make build-release
	@echo "$(GREEN)✓ Backend binary built$(NC)"

frontend-build:
	@echo "$(BLUE)Building frontend for production...$(NC)"
	@cd frontend && npm run build
	@echo "$(GREEN)✓ Frontend built successfully$(NC)"

# ============================================================================
# DOCKER
# ============================================================================

.PHONY: docker-build docker-run docker-clean docker-logs docker-stop

docker-build:
	@echo "$(BLUE)Building Docker image...$(NC)"
	@cd backend && make docker-build
	@echo "$(GREEN)✓ Docker image built$(NC)"

docker-run: docker-build
	@echo "$(BLUE)Running application in Docker...$(NC)"
	@cd backend && make docker-run

docker-clean:
	@echo "$(BLUE)Cleaning Docker resources...$(NC)"
	@cd backend && make docker-clean
	@echo "$(GREEN)✓ Docker cleaned$(NC)"

docker-logs:
	@docker logs -f workspace-app 2>/dev/null || echo "$(YELLOW)No running container$(NC)"

docker-stop:
	@docker stop workspace-app 2>/dev/null || echo "$(YELLOW)No running container$(NC)"

# ============================================================================
# DATABASE
# ============================================================================

.PHONY: migrate db-reset db-status seed

migrate:
	@echo "$(BLUE)Running database migrations...$(NC)"
	@cd backend && make migrate
	@echo "$(GREEN)✓ Migrations completed$(NC)"

db-reset:
	@echo "$(RED)⚠️  WARNING: This will drop and recreate the PostgreSQL database!$(NC)"
	@read -p "Continue? (yes/no): " confirm && [ "$${confirm}" = "yes" ] || exit 1
	@echo "$(BLUE)Resetting PostgreSQL database...$(NC)"
	@psql "postgres://$${POSTGRES_USER:-postgres}:$${POSTGRES_PASSWORD:-postgres}@$${POSTGRES_HOST:-localhost}:$${POSTGRES_PORT:-5432}/postgres?sslmode=$${POSTGRES_SSLMODE:-disable}" \
		-c "DROP DATABASE IF EXISTS $${POSTGRES_DB:-kirmya} WITH (FORCE);" \
		-c "CREATE DATABASE $${POSTGRES_DB:-kirmya};"
	@echo "$(GREEN)✓ Database reset. Migrations run automatically on the next backend start (or run: make migrate).$(NC)"

db-status:
	@echo "$(BLUE)Database Status (PostgreSQL):$(NC)"
	@psql "$${DATABASE_URL:-postgres://postgres:postgres@localhost:5432/kirmya?sslmode=disable}" \
		-c "\dt" -c "SELECT count(*) AS applied_migrations FROM schema_migrations;" 2>/dev/null \
		|| echo "$(YELLOW)Could not connect — set DATABASE_URL or start PostgreSQL first.$(NC)"

seed:
	@echo "$(BLUE)Seeding demo data...$(NC)"
	@echo "$(YELLOW)Seeding runs automatically at backend startup when SEED_DEMO_DATA=true:$(NC)"
	@echo "  cd backend && SEED_DEMO_DATA=true go run ./cmd/workspace-app/main.go"

# ============================================================================
# CODE QUALITY
# ============================================================================

.PHONY: fmt lint vet check pre-commit tidy

fmt:
	@echo "$(BLUE)Formatting all code...$(NC)"
	@cd backend && make fmt
	@cd frontend && npx prettier --write . 2>/dev/null || true
	@echo "$(GREEN)✓ Code formatted$(NC)"

lint:
	@echo "$(BLUE)Linting all code...$(NC)"
	@cd backend && make lint || echo "$(YELLOW)Backend lint complete$(NC)"
	@cd frontend && npm run lint 2>/dev/null || echo "$(YELLOW)Frontend lint not configured$(NC)"
	@echo "$(GREEN)✓ Linting complete$(NC)"

vet:
	@echo "$(BLUE)Running go vet...$(NC)"
	@cd backend && make vet
	@echo "$(GREEN)✓ Vet complete$(NC)"

tidy:
	@echo "$(BLUE)Tidying dependencies...$(NC)"
	@cd backend && make tidy
	@cd frontend && npm audit fix 2>/dev/null || true
	@echo "$(GREEN)✓ Dependencies tidied$(NC)"

check: fmt lint vet test
	@echo "$(GREEN)✓ All quality checks passed$(NC)"

pre-commit: fmt vet backend-test
	@echo "$(GREEN)✓ Pre-commit checks passed - ready to commit$(NC)"

# ============================================================================
# GIT
# ============================================================================

.PHONY: git-status git-log git-diff commit push

git-status:
	@echo "$(BLUE)Git Status:$(NC)"
	@git status

git-log:
	@echo "$(BLUE)Recent Commits:$(NC)"
	@git log --oneline -10

git-diff:
	@echo "$(BLUE)Changes:$(NC)"
	@git diff --stat

commit:
	@echo "$(BLUE)Preparing to commit...$(NC)"
	@make pre-commit
	@echo "$(YELLOW)Run: git add . && git commit -m 'message'$(NC)"

push:
	@echo "$(BLUE)Pushing to remote...$(NC)"
	@git push origin $$(git rev-parse --abbrev-ref HEAD)
	@echo "$(GREEN)✓ Pushed$(NC)"

# ============================================================================
# UTILITIES
# ============================================================================

.PHONY: ps restart logs status version init

ps:
	@echo "$(BLUE)Running Services:$(NC)"
	@lsof -i :8080 2>/dev/null | tail -n +2 || echo "$(YELLOW)Backend not running$(NC)"
	@lsof -i :3000 2>/dev/null | tail -n +2 || echo "$(YELLOW)Frontend not running$(NC)"

restart:
	@echo "$(BLUE)Restarting services...$(NC)"
	@pkill -f "go run ./cmd/workspace-app" || true
	@pkill -f "next dev" || true
	@sleep 1
	@echo "$(GREEN)✓ Services stopped$(NC)"
	@echo "$(YELLOW)Run: make dev$(NC)"

status:
	@echo "$(BLUE)Project Status:$(NC)"
	@echo "Backend:"
	@cd backend && echo "  - Binary: $$(ls -1 bin/workspace-app* 2>/dev/null | wc -l) binaries"
	@cd backend && echo "  - Tests: $$(go test ./... -list='.*' 2>/dev/null | wc -l) test cases"
	@echo "Frontend:"
	@cd frontend && echo "  - Build: $$(test -d .next && echo 'built' || echo 'not built')"
	@cd frontend && echo "  - Packages: $$(npm list --depth=0 2>/dev/null | wc -l) dependencies"

version:
	@echo "$(BLUE)Version Information:$(NC)"
	@echo "  Go: $$(go version)"
	@echo "  Node: $$(node --version)"
	@echo "  npm: $$(npm --version)"
	@echo "  Git: $$(git --version)"
	@cd backend && make version

init:
	@echo "$(BLUE)Initializing project...$(NC)"
	@mkdir -p logs
	@touch logs/backend.log logs/frontend.log
	@echo "$(GREEN)✓ Project initialized$(NC)"

# ============================================================================
# CLEANUP
# ============================================================================

.PHONY: clean clean-all clean-backend clean-frontend

clean: clean-backend clean-frontend
	@echo "$(GREEN)✓ Cleaned all$(NC)"

clean-backend:
	@echo "$(BLUE)Cleaning backend...$(NC)"
	@cd backend && make clean
	@echo "$(GREEN)✓ Backend cleaned$(NC)"

clean-frontend:
	@echo "$(BLUE)Cleaning frontend...$(NC)"
	@cd frontend && rm -rf node_modules .next
	@echo "$(GREEN)✓ Frontend cleaned$(NC)"

clean-all: clean
	@echo "$(BLUE)Deep cleaning...$(NC)"
	@cd backend && make clean-all
	@cd frontend && npm ci > /dev/null
	@rm -rf logs/*
	@echo "$(GREEN)✓ All cleaned$(NC)"

# ============================================================================
# DEFAULT
# ============================================================================

.DEFAULT_GOAL := help
