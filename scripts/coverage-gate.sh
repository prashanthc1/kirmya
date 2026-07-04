#!/usr/bin/env bash
# coverage-gate.sh
#
# Enforces the Definition-of-Done coverage bar (CLAUDE.md §7): the combined
# statement coverage of every module's application/ and domain/ package must
# be >= THRESHOLD (default 70%). Anything else in the tree (api/, infrastructure/,
# platform/, cmd/) is intentionally excluded -- the bar is on the use-case and
# domain logic, which is unit-tested with fakes and must not regress.
#
# Usage (run from backend/):
#   scripts/coverage-gate.sh            # threshold 70
#   COVERAGE_THRESHOLD=80 scripts/coverage-gate.sh
#
# Exit code 0 when coverage >= threshold, 1 when below, 2 on a hard error.
set -euo pipefail

THRESHOLD="${COVERAGE_THRESHOLD:-70}"
COVERAGE_DIR="${COVERAGE_DIR:-coverage}"
PROFILE="${COVERAGE_DIR}/app-domain.out"

mkdir -p "${COVERAGE_DIR}"

# The packages we measure: application + domain only. Enumerate every package
# and keep those whose import path contains an /application or /domain segment.
mapfile -t PKGS < <(go list ./internal/... 2>/dev/null | grep -E '/(application|domain)(/|$)' | sort -u)

if [ "${#PKGS[@]}" -eq 0 ]; then
  echo "coverage-gate: found no application/domain packages to measure" >&2
  exit 2
fi

echo "coverage-gate: measuring ${#PKGS[@]} application/domain packages (threshold ${THRESHOLD}%)"

# -coverpkg restricts instrumentation to the measured set so cross-package calls
# (e.g. api -> application) do not dilute the percentage. We still run only the
# measured packages' tests.
COVERPKG=$(IFS=,; echo "${PKGS[*]}")
go test -covermode=atomic -coverpkg="${COVERPKG}" -coverprofile="${PROFILE}" "${PKGS[@]}"

TOTAL=$(go tool cover -func="${PROFILE}" | awk '/^total:/ {gsub(/%/,"",$3); print $3}')
if [ -z "${TOTAL}" ]; then
  echo "coverage-gate: could not parse total coverage from ${PROFILE}" >&2
  exit 2
fi

echo "coverage-gate: application+domain coverage = ${TOTAL}% (threshold ${THRESHOLD}%)"

# Float comparison via awk (bash cannot compare decimals).
if awk -v t="${TOTAL}" -v thr="${THRESHOLD}" 'BEGIN { exit !(t + 0 < thr + 0) }'; then
  echo "::error::coverage ${TOTAL}% is below the required ${THRESHOLD}% (application+domain)"
  exit 1
fi

echo "coverage-gate: PASS"
