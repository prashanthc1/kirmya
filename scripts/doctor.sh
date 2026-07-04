#!/usr/bin/env bash
# doctor.sh — preflight environment check for Kirmya.
#
# Verifies the host has everything the Makefile / docker-compose targets need.
# The repo's Makefiles and Dockerfiles are correct; "it doesn't work" is almost
# always a missing tool or a version older than what go.mod / the images pin.
#
#   make doctor        # or:  bash scripts/doctor.sh
#
# Exit code 0 = all required checks passed, 1 = at least one required check failed.

set -uo pipefail

GREEN='\033[0;32m'; RED='\033[0;31m'; YELLOW='\033[1;33m'; BLUE='\033[0;34m'; NC='\033[0m'
fail=0

ok()   { printf "${GREEN}  ✓ %s${NC}\n" "$1"; }
bad()  { printf "${RED}  ✗ %s${NC}\n" "$1"; fail=1; }
warn() { printf "${YELLOW}  ! %s${NC}\n" "$1"; }
section() { printf "\n${BLUE}%s${NC}\n" "$1"; }

# Compare two dotted versions: returns 0 if $1 >= $2.
ver_ge() { [ "$(printf '%s\n%s\n' "$2" "$1" | sort -V | head -n1)" = "$2" ]; }

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

printf "${BLUE}Kirmya environment doctor${NC}\n"
printf "repo: %s\n" "$ROOT"

# ---- Go (backend) -----------------------------------------------------------
section "Go (backend — required for make dev/build/test)"
NEED_GO="1.26"
if command -v go >/dev/null 2>&1; then
  GOV="$(go version | grep -oE 'go[0-9]+\.[0-9]+(\.[0-9]+)?' | head -n1 | sed 's/go//')"
  if [ -n "$GOV" ] && ver_ge "$GOV" "$NEED_GO"; then
    ok "go $GOV (>= $NEED_GO required by go.mod)"
  else
    bad "go $GOV found, but go.mod needs >= $NEED_GO — upgrade Go (https://go.dev/dl/)"
  fi
else
  bad "go not found in PATH — install Go >= $NEED_GO (https://go.dev/dl/). Without it, every backend make target fails."
fi

# ---- Node / npm (frontend) --------------------------------------------------
section "Node.js + npm (frontend)"
NEED_NODE="22"
if command -v node >/dev/null 2>&1; then
  NODEV="$(node --version | sed 's/v//')"
  if ver_ge "$NODEV" "$NEED_NODE.0.0"; then ok "node $NODEV (>= $NEED_NODE)"; else bad "node $NODEV — Dockerfile/CI use Node $NEED_NODE; upgrade"; fi
else
  bad "node not found — install Node $NEED_NODE+ (https://nodejs.org/)"
fi
command -v npm >/dev/null 2>&1 && ok "npm $(npm --version)" || bad "npm not found"

# ---- Docker (compose stack) -------------------------------------------------
section "Docker (only needed for 'docker compose up' / make docker-*)"
if command -v docker >/dev/null 2>&1; then
  ok "docker $(docker --version | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | head -n1)"
  if docker info >/dev/null 2>&1; then ok "docker daemon is running"; else bad "docker is installed but the daemon is NOT running — start Docker Desktop / the docker service"; fi
  if docker compose version >/dev/null 2>&1; then ok "docker compose v2 ($(docker compose version --short 2>/dev/null))"; else bad "'docker compose' v2 plugin missing (the v1 'docker-compose' is not used by this repo)"; fi
else
  warn "docker not found — fine if you run natively with 'make dev'; required only for the compose stack"
fi

# ---- .env files -------------------------------------------------------------
section "Environment files"
if [ -f "$ROOT/backend/.env" ]; then ok "backend/.env present"; else warn "backend/.env missing — run 'make setup' (copies backend/.env.example). 'make dev' needs DB creds here."; fi
if [ -f "$ROOT/frontend/.env.local" ]; then ok "frontend/.env.local present"; else warn "frontend/.env.local missing (optional for local dev)"; fi

# ---- Ports the stack binds --------------------------------------------------
section "Ports (must be free for the stack to bind)"
check_port() {
  local p="$1" name="$2"
  if command -v lsof >/dev/null 2>&1; then
    if lsof -iTCP:"$p" -sTCP:LISTEN >/dev/null 2>&1; then warn "port $p ($name) is IN USE — stop the other process or it will collide"; else ok "port $p ($name) free"; fi
  else
    ok "port $p ($name) — install lsof to check (skipped)"
  fi
}
check_port 3000 frontend
check_port 8080 backend
check_port 5432 postgres
check_port 6379 redis

# ---- Summary ----------------------------------------------------------------
echo
if [ "$fail" -eq 0 ]; then
  printf "${GREEN}All required checks passed. Try: make setup && make dev${NC}\n"
else
  printf "${RED}Some required checks failed (see ✗ above). Fix those, then re-run: make doctor${NC}\n"
fi
exit "$fail"
