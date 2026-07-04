# Kirmya — Agents & Scaling Strategy

> Status: Draft v1 · Last updated: 2026-06-15
> Companion to `10_MVP_ROADMAP.md`. Describes the specialized AI agents used to build
> Kirmya in parallel and how to grow the codebase without losing consistency.

## 1. Why agents

Kirmya is a modular monolith with ~12 bounded contexts on the backend and a
matching set of App Router features on the frontend. The work is highly parallelizable
*because* the module boundaries are clean — but only if every slice is built to the
same DDD layering, contract, and test bar. Specialized agents encode those conventions
once and apply them consistently, so multiple modules can advance at the same time
without architectural drift.

The agent definitions live in `.claude/agents/` and are invoked per task.

## 2. The agent roster

| Agent | Scope | Primary inputs | Definition |
|---|---|---|---|
| **backend-module-builder** | Go vertical slices (domain → application → infra → api → wiring) | An existing module to mirror; `06_API_CONTRACTS.md` | `.claude/agents/backend-module-builder.md` |
| **frontend-feature-builder** | Next.js App Router routes + components for a module | The backend contract; an existing comparable route | `.claude/agents/frontend-feature-builder.md` |
| **api-contract-guardian** | Keep docs ↔ Go handlers ↔ frontend client in sync | `06_API_CONTRACTS.md`, `api/routes.go`, `frontend/lib` | `.claude/agents/api-contract-guardian.md` |
| **test-coverage-engineer** | Unit/integration/e2e tests; coverage gating (≥70% app+domain) | `TESTING.md`, `vitest.config.ts`, `playwright.config.ts` | `.claude/agents/test-coverage-engineer.md` |
| **security-auth-reviewer** | OWASP review of auth/RBAC/CSRF/OAuth/data-exposure | `AUTHENTICATION.md`, `CSRF_SECURITY.md`, identity module | `.claude/agents/security-auth-reviewer.md` |
| **devops-release-engineer** | Docker, GitHub Actions, Helm, Prometheus/Grafana, releases | `10_MVP_ROADMAP.md` §3–5, `deploy/helm`, `ops/`, `MAKEFILE_GUIDE.md` | `.claude/agents/devops-release-engineer.md` |

These complement the installed **engineering** plugin skills (`system-design`,
`architecture`, `code-review`, `testing-strategy`, `debug`, `deploy-checklist`,
`documentation`, `tech-debt`), which the agents can invoke for cross-cutting work.

## 3. How the agents collaborate per feature

A typical end-to-end feature flows through the agents like a small assembly line:

1. **architecture** skill (optional) — record the design decision as an ADR when there
   is a real trade-off.
2. **backend-module-builder** — implement the Go slice + migration + service tests, wire
   `module.go`, update `06_API_CONTRACTS.md`.
3. **api-contract-guardian** — confirm the documented contract matches the new handlers
   before any frontend work starts.
4. **frontend-feature-builder** — build the route/components against the (now trusted)
   contract, with component + e2e tests.
5. **security-auth-reviewer** — required gate whenever the slice touches auth, RBAC,
   ownership checks, or PII exposure.
6. **test-coverage-engineer** — fill remaining gaps to the ≥70% bar; wire coverage into CI.
7. **devops-release-engineer** — extend compose/Helm/observability and the deploy
   checklist if the feature adds infra (new table, queue, search index, env var).

## 4. Mapping to the MVP phases

Phase numbers follow `10_MVP_ROADMAP.md`. Module maturity below reflects the current
backend file counts (identity is the most built-out; community and mentorship are the
thinnest).

| Phase | Modules | Lead agents |
|---|---|---|
| **1 — Identity & Profiles** (mostly done) | identity, profile | security-auth-reviewer, test-coverage-engineer, frontend-feature-builder |
| **2 — Resume, Jobs, Referrals** | resume, jobs, referrals, search | backend-module-builder, api-contract-guardian, frontend-feature-builder |
| **3 — Intelligence, Communities, Mentorship** | ai, community, mentorship | backend-module-builder (community/mentorship are thin), frontend-feature-builder |
| **4 — Messaging, Notifications, Admin, Hardening** | messaging, notifications, admin | devops-release-engineer, security-auth-reviewer, test-coverage-engineer |

## 5. Running agents in parallel safely

The module boundary is the unit of parallelism. To avoid collisions:

- **One agent owns one module slice at a time.** Two agents must not edit the same
  package concurrently. The shared seams — `module.go` registration, `platform`,
  migrations, and `06_API_CONTRACTS.md` — are serialization points: change them in a
  dedicated, short task rather than inside two parallel slices.
- **Contract-first.** The api-contract-guardian validates the contract before the
  frontend agent starts, so backend and frontend can then proceed independently.
- **Isolation for risky work.** Run a builder agent in a git worktree when the change
  is large or experimental, so `main` stays green; merge after build+tests pass.
- **CI is the real gate.** `ci.yml` (lint/test/build) must pass on every PR regardless
  of which agent produced the code — agents do not bypass it.

## 6. Definition of done (every agent slice)

- Code matches the DDD layering and existing module style.
- `06_API_CONTRACTS.md` updated for any endpoint change.
- Backend: `gofmt` clean, `go build ./...` and `go test ./internal/<module>/...` green;
  app+domain coverage ≥70%.
- Frontend: `npm run lint`, `tsc --noEmit`, `npm run build` clean; touched components
  tested; critical journeys have a Playwright e2e.
- Auth-touching slices reviewed by security-auth-reviewer.
- New infra reflected in compose/Helm/observability and a stated rollback path.

## 7. Next steps to "big project"

1. Add a top-level `CLAUDE.md` capturing the conventions above so every agent (and human)
   starts with the same context. (Not present yet.)
2. Stand up the contract-guardian as a CI check, not just an on-demand agent.
3. Knock out the thin modules (community, mentorship) to full vertical slices using
   backend-module-builder + frontend-feature-builder.
4. Promote messaging and AI toward extractable services per the roadmap's multi-region
   plan, once their contracts are stable.
