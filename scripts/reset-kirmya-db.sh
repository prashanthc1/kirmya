#!/usr/bin/env bash
#
# reset-kirmya-db.sh — finish the Career Bridge -> Kirmya rebrand at runtime.
#
# What it does:
#   1. Repairs a corrupt git index (if present) so `git status` works again.
#   2. Creates a fresh `kirmya` database.
#   3. Drops the old `career_bridge` database.
# After it runs, start the backend (`make backend-run`) — migrations re-apply
# at boot and create the correct schema, which also fixes the signup 500.
#
# Usage (from Git Bash / WSL, at the repo root):
#   bash scripts/reset-kirmya-db.sh
#
# Connection settings are read from the environment / backend/.env defaults.
set -euo pipefail

# Resolve repo root (this script lives in <root>/scripts).
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

# Load DB settings from backend/.env if present (without clobbering existing env).
if [ -f backend/.env ]; then
  set -a
  # shellcheck disable=SC1091
  . backend/.env
  set +a
fi

PGHOST="${POSTGRES_HOST:-127.0.0.1}"
PGPORT="${POSTGRES_PORT:-5432}"
PGUSER="${POSTGRES_USER:-postgres}"
PGPASSWORD="${POSTGRES_PASSWORD:-postgres}"
PGSSL="${POSTGRES_SSLMODE:-disable}"
export PGPASSWORD
ADMIN_URL="postgres://${PGUSER}@${PGHOST}:${PGPORT}/postgres?sslmode=${PGSSL}"

echo "==> 1/3  Repairing git index (if needed)"
rm -f .git/index.lock .git/index 2>/dev/null || true
git reset -q 2>/dev/null || true
echo "    git index rebuilt; run 'git status' afterwards."

echo "==> 2/3  Recreating database 'kirmya' on ${PGHOST}:${PGPORT}"
if ! command -v psql >/dev/null 2>&1; then
  echo "    psql not found on PATH."
  echo "    If you run Postgres via docker compose, use instead:"
  echo "        docker compose down -v && docker compose up -d"
  exit 1
fi
psql "$ADMIN_URL" -v ON_ERROR_STOP=1 -c "DROP DATABASE IF EXISTS kirmya WITH (FORCE);"
psql "$ADMIN_URL" -v ON_ERROR_STOP=1 -c "CREATE DATABASE kirmya;"

echo "==> 3/3  Dropping old database 'career_bridge'"
psql "$ADMIN_URL" -v ON_ERROR_STOP=1 -c "DROP DATABASE IF EXISTS career_bridge WITH (FORCE);"

echo
echo "Done. 'kirmya' is ready and 'career_bridge' is gone."
echo "Next: start the backend so migrations apply ->  make backend-run"
