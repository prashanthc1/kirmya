# Makefile Guide - Recession Recovery Workspace

Complete guide for all Makefile commands across the project.

## Table of Contents

1. [Quick Start](#quick-start)
2. [Root Makefile](#root-makefile)
3. [Backend Makefile](#backend-makefile)
4. [Frontend Makefile](#frontend-makefile)
5. [Common Workflows](#common-workflows)
6. [Troubleshooting](#troubleshooting)

## Quick Start

```bash
# Initial setup
make setup              # Install deps, create .env files
make migrate            # Run database migrations
make dev                # Start both backend and frontend

# Or start individually
make backend-run        # Backend only
make frontend-run       # Frontend only
```

## Root Makefile

**Location:** `./Makefile`

The root Makefile orchestrates the entire project. Use this for full-stack operations.

### Setup & Installation

```bash
make setup              # Complete project setup
make install-deps       # Install backend and frontend dependencies
make install-tools      # Install development tools (golangci-lint, etc)
make init-env           # Create .env files from examples
```

### Development

```bash
# Start everything
make dev                # Start backend + frontend (foreground, watch mode)
make dev-bg             # Start backend + frontend (background)

# Start individually
make backend-run        # Start Go backend
make backend-dev        # Start backend with hot reload (requires CompileDaemon)
make frontend-run       # Start Next.js frontend

# View logs
make backend-logs       # Tail backend logs
make frontend-logs      # Tail frontend logs
```

### Testing

```bash
make test               # Run all tests (backend + frontend)
make backend-test       # Run backend tests
make backend-test-v     # Run backend tests (verbose)
make backend-coverage   # Generate coverage report (HTML)
make frontend-test      # Run frontend tests
```

### Build & Deploy

```bash
make build              # Build backend release binary
make build-all          # Build backend + frontend
make docker-build       # Build Docker image
make docker-run         # Run in Docker
```

### Database

```bash
make migrate            # Run migrations
make db-reset           # Reset database (with confirmation)
make db-status          # Show database status
make seed               # Seed with sample data (not yet implemented)
```

### Code Quality

```bash
make fmt                # Format all code (Go + frontend)
make lint               # Lint all code
make vet                # Run go vet
make tidy               # Tidy dependencies
make check              # Run all quality checks
make pre-commit         # Pre-commit checks (fmt, vet, test)
```

### Git

```bash
make git-status         # Show git status
make git-log            # Show recent commits
make git-diff           # Show changes
make commit             # Prepare commit (runs pre-commit checks)
make push               # Push to remote
```

### Utilities

```bash
make ps                 # Show running services
make restart            # Restart all services
make status             # Show project status
make version            # Show version information
make clean              # Remove build artifacts
make clean-all          # Complete cleanup
```

## Backend Makefile

**Location:** `./backend/Makefile`

### Build Targets

```bash
make build              # Build release binary with version info
make build-dev          # Build debug binary (faster, larger)
make build-release      # Build optimized release binary
make build-linux        # Build for Linux (x86_64)
make build-macos        # Build for macOS (x86_64)
make build-windows      # Build for Windows (x86_64)
```

### Run Targets

```bash
make run                # Run application (requires .env)
make run-dev            # Run with watch mode (requires CompileDaemon)
make run-migrate        # Run migrations then start app
```

### Test Targets

```bash
make test               # Run all unit tests
make test-verbose       # Run tests with verbose output
make test-jobs          # Run only jobs domain tests
make test-ideas         # Run only ideas domain tests
make test-user          # Run only user domain tests
make test-coverage      # Run tests and generate coverage report
make coverage           # Generate coverage report (text)
make coverage-html      # Generate HTML coverage report
make coverage-report    # Display coverage summary
```

### Code Quality

```bash
make fmt                # Format all Go code
make vet                # Run go vet static analysis
make lint               # Run golangci-lint
make tidy               # Clean up go.mod and go.sum
make deps               # Verify dependencies
make check              # Run all quality checks (fmt, vet, tidy, test)
make pre-commit         # Run pre-commit checks (fmt, vet, test)
```

### Database

```bash
make migrate            # Run database migrations
make migrate-reset      # Reset database (drop + recreate)
make migrate-status     # Show migration status
```

### Docker

```bash
make docker-build       # Build Docker image
make docker-run         # Run application in Docker
make docker-clean       # Remove Docker image
```

### Utilities

```bash
make install-tools      # Install development tools
make version            # Show build version info
make docs               # Show documentation links
make init-env           # Create .env file from .env.example
make clean              # Remove build artifacts
make clean-coverage     # Remove coverage files
make clean-all          # Remove all generated files
```

## Frontend Makefile

**Location:** `./frontend/Makefile`

### Development

```bash
make dev                # Run development server (port 3000)
make dev-turbo          # Run with Turbopack (faster, experimental)
```

### Build

```bash
make build              # Build for production
make build-analyze      # Build and analyze bundle size
make start              # Run production build
make serve              # Serve production build locally
```

### Testing

```bash
make test               # Run test suite
make test-watch         # Run tests in watch mode
make test-coverage      # Run tests with coverage
make test-debug         # Run tests in debug mode
```

### Code Quality

```bash
make lint               # Run ESLint
make lint-fix           # Fix linting issues
make fmt                # Format code with Prettier
make fmt-check          # Check formatting
make type-check         # Run TypeScript compiler
make check              # Run all quality checks
```

### Dependencies

```bash
make install            # Install dependencies (npm ci)
make update             # Update dependencies
make audit              # Check for vulnerabilities
make audit-fix          # Fix vulnerabilities
make deps               # List dependencies
```

### Utilities

```bash
make version            # Show versions (Node, npm, Next.js)
make info               # Show project information
make env-check          # Check environment setup
make clean              # Remove .next build directory
make clean-all          # Remove all generated files (node_modules, .next, etc)
make clean-cache        # Clear Next.js cache
```

### Pre-commit & Pre-push

```bash
make pre-commit         # Run pre-commit checks (lint, type-check, fmt)
make pre-push           # Run pre-push checks (build, lint, type-check)
```

## Common Workflows

### First Time Setup

```bash
# 1. Clone and navigate
cd My\ Site

# 2. Install everything
make setup

# 3. Update .env with database credentials
vi backend/.env

# 4. Run migrations
make migrate

# 5. Start development
make dev
```

### Daily Development

```bash
# Start everything
make dev

# In another terminal, run tests
make test

# Format and lint before committing
make fmt lint pre-commit
```

### Code Review & Quality

```bash
# Run comprehensive checks
make check              # All quality checks
make backend-coverage   # See coverage reports
make pre-commit         # Final check before commit
```

### Building for Production

```bash
# Backend
cd backend
make build-release      # Creates bin/workspace-app-release

# Frontend
cd frontend
make build              # Creates .next/

# Or from root
make build-all
```

### Docker Deployment

```bash
# Build and run Docker container
make docker-build
make docker-run

# Or directly
cd backend
make docker-run
```

### Database Operations

```bash
# Run migrations
make migrate

# Reset database (careful!)
make db-reset

# Check status
make db-status
```

## Troubleshooting

### Port Already in Use

```bash
# Check what's using the ports
make ps

# Kill existing processes
make restart

# Then restart
make dev
```

### Dependencies Issues

```bash
# Reinstall dependencies
cd backend && go mod tidy
cd frontend && make clean-all && make install
```

### Build Failures

```bash
# Clean and rebuild
make clean-all
make install-deps
make build
```

### Test Failures

```bash
# Run verbose tests
make backend-test-v

# Check coverage
make backend-coverage

# Run specific domain tests
make test-jobs
```

### Database Connection Issues

```bash
# Check environment
make env-check

# Verify .env file
cat backend/.env

# Test MySQL connection
mysql -u root -p -h localhost
```

### Hot Reload Not Working

```bash
# Backend hot reload requires CompileDaemon
make install-tools

# Then use
make backend-dev
```

## Best Practices

1. **Run `make pre-commit` before committing code**
   ```bash
   make pre-commit  # Runs fmt, vet, test
   ```

2. **Use `make check` for comprehensive quality checks**
   ```bash
   make check       # Runs fmt, lint, vet, test
   ```

3. **Keep .env files updated**
   ```bash
   make init-env    # Creates from .env.example
   ```

4. **Regular cleanup to avoid disk space issues**
   ```bash
   make clean       # Clean build artifacts
   make clean-all   # Deep cleanup (runs dependencies again)
   ```

5. **Check logs when things go wrong**
   ```bash
   make backend-logs
   make frontend-logs
   ```

## Tips & Tricks

### Run specific backend tests
```bash
cd backend
make test-jobs          # Jobs tests only
make test-ideas         # Ideas tests only
```

### Generate coverage reports
```bash
make backend-coverage   # Generates HTML report at backend/coverage/coverage.html
```

### Parallel development
```bash
# Terminal 1: Backend
make backend-dev

# Terminal 2: Frontend
make frontend-run

# Terminal 3: Tests
make backend-test-v
```

### Dry run format check
```bash
cd frontend
make fmt-check          # Shows what would be formatted
make fmt                # Actually format
```

### View all make commands
```bash
make help               # Root level
cd backend && make help # Backend
cd frontend && make help # Frontend
```

---

For questions or issues, check the relevant README:
- Backend: `./backend/README.md`
- Frontend: `./frontend/README.md`
- Main: `./README.md`
