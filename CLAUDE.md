# CLAUDE.md — Kirmya

Shared context for every agent and human working in this repo. Read this first, then
the deeper docs in `docs/`. Keep it accurate: if you change a convention, update this file.

## 1. Project overview

Kirmya (a.k.a. the *Recession Recovery Workspace*) is a platform for connecting
job seekers, recruiters, founders, freelancers, mentors, and collaborators during
recession and transition periods. It is a monorepo with two projects:

- **`backend/`** — a Go (1.26, module `workspace-app`) **modular monolith** HTTP API.
  One binary, ~12 bounded contexts, Clean Architecture / DDD per module. Runs at
  `http://localhost:8080`, REST/JSON under `/api/v1`.
- **`frontend/`** — a **Next.js** (App Router, React 19, TypeScript, TailwindCSS,
  ShadCN UI) web app. Mobile-first. Runs at `http://localhost:3000`.

Backing services: PostgreSQL (pgx), Redis (cache), OpenSearch (search, with a
PostgreSQL `ILIKE` fallback), an in-process event bus (NATS-ready) + transactional
outbox, and OpenTelemetry/Prometheus observability.

## 2. Repo layout

```
backend/
├── cmd/workspace-app/main.go     # composition root: config, DB, wire modules, serve
├── internal/
│   ├── platform/                 # framework/wiring (config, database, router, server,
│   │                             #   middleware, eventbus, cache, search, observability)
│   ├── common/                   # shared kernel: errors, response, context, pagination, ids
│   └── <module>/                 # one bounded context per folder (DDD layout below)
├── migrations/                   # PostgreSQL forward-only NNN_name.sql + seed/
└── docs/openapi.yaml             # machine source of truth, served at /swagger-ui/
```

Each module follows the same DDD layering (gold standard: `internal/identity`):

```
internal/<module>/
├── module.go            # composition root: NewModule(deps) -> RegisterRoutes, wires repo→service→handler
├── domain/              # entities, value objects, sentinel errors (errors.go-style), ports.go, events.go
├── application/         # use-case Service + Deps; depends only on domain ports; unit-tested with fakes
├── infrastructure/      # adapters implementing domain ports (postgres/, crypto/, mailer/, oauth/, jwtauth/)
└── api/                 # delivery: routes.go, handlers.go, dto.go, middleware.go
```

Frontend (App Router):

```
frontend/
├── app/                 # routes = URLs; route groups: (marketing), (auth), (app), (admin)
├── src/
│   ├── components/      # ui/ (ShadCN primitives) + shared/ (composed components)
│   ├── features/<f>/    # one folder per domain (mirrors backend): api.ts, schemas.ts, hooks.ts, components/, types.ts
│   ├── lib/             # api/client.ts (auth header, CSRF, refresh-on-401), auth/, hooks/, utils/, config.ts
│   └── types/           # shared/generated TS types
└── tests/               # components/ (Vitest + RTL) and e2e/ (Playwright)
```

> Both repos are mid-migration. Some backend modules still use the older
> `service/repository/handler` layout and some frontend pages predate `features/`;
> migrate to the target layout *as you touch them*. Everything compiles under one binary/build.

## 3. DDD conventions (backend)

- **Dependency direction:** `api → application → domain ← infrastructure`. `domain`
  imports nothing from other layers or modules. Modules depend on each other **only**
  via another module's `application.Service` interface, injected in `main.go`.
- **Ports:** repositories and gateways are interfaces declared in `domain/ports.go`;
  `infrastructure/` provides the concrete adapters.
- **`context.Context` is the first argument** of every port and service method.
- **Sentinel errors** live in the domain package (e.g. `ErrUserNotFound`,
  `ErrEmailTaken`, `ErrOptimisticLock`) and are mapped to HTTP statuses in `api/`.
- **Event bus is best-effort:** construct services so `events` may be nil and publish
  defensively — `if s.events != nil { _ = s.events.Publish(ctx, eventType, aggregateID, payload) }`.
  Event types are PascalCase string consts (e.g. `MentorshipBooked`, `JobPosted`).
- **Optimistic locking:** aggregates carry a `Version`; `Update` is version-checked and
  returns the optimistic-lock sentinel (HTTP 409) on a stale write.
- **Migrations:** new tables get a forward-only `NNN_name.sql` file in
  `backend/migrations/` (applied in sequence by `platform/migrate`; there are no
  down files / automated rollback), matching existing naming/sequence.

## 4. Build / run / test commands

The top-level **`Makefile`** orchestrates both stacks (run `make help` for the full
menu). Common targets:

| Task | Make target | Underlying command |
|---|---|---|
| Setup deps + `.env` | `make setup` | `go mod download` / `npm ci` + copy `.env.example` |
| Run both (watch) | `make dev` | backend `go run` + frontend `npm run dev` |
| Run backend only | `make backend-run` | `cd backend && go run ./cmd/workspace-app/main.go` |
| Run frontend only | `make frontend-run` | `cd frontend && npm run dev` |
| All tests | `make test` | `make backend-test` + `make frontend-test` |
| Backend tests | `make backend-test` | `cd backend && make test` (→ `go test ./...`) |
| Backend coverage | `make backend-coverage` | `cd backend && make coverage-html` |
| Frontend tests | `make frontend-test` | `cd frontend && npm run test` |
| Build backend binary | `make build` | `cd backend && make build-release` |
| Build everything | `make build-all` | backend binary + `npm run build` |
| Format / lint / vet | `make fmt` / `make lint` / `make vet` | gofmt + prettier / golangci + eslint / go vet |
| All quality checks | `make check` | `fmt` + `lint` + `vet` + `test` |
| Migrations / seed | `make migrate` / `make seed` | `cd backend && make migrate` |

Direct commands (use these when verifying a single slice; the Go toolchain is required):

- **Backend:** `cd backend && gofmt -w ./... && go build ./... && go test ./internal/<module>/...`.
  Never finish on a failing build or test.
- **Frontend** (`cd frontend`): `npm run lint`, `npx tsc --noEmit`, `npm run build`,
  `npm run test` (Vitest), `npm run e2e` (Playwright; `npm run e2e:install` once for browsers).

API docs while running: OpenAPI at `http://localhost:8080/openapi.yaml`, Swagger UI at
`/swagger-ui/`, health at `/api/v1/health`.

## 5. Contract-first workflow

- `docs/06_API_CONTRACTS.md` is the **human-readable** contract (conventions, error
  envelope, per-module endpoint tables). Start here when designing or consuming an endpoint.
- `backend/docs/openapi.yaml` is the **machine source of truth**, served at `/swagger-ui/`.
- Any endpoint change updates **both** the docs table and the OpenAPI spec, and is
  flagged in your task summary. Backend handlers and the frontend `lib/api` client must
  stay in sync with them — the **api-contract-guardian** agent validates this before
  frontend work begins (and is on the roadmap to become a CI check).

Shared conventions worth memorizing: base path `/api/v1`; UUID string IDs; RFC 3339 UTC
timestamps; bearer access token + httpOnly refresh cookie (rotated, reuse-revokes-family);
Bearer auth is CSRF-immune (optional Origin check, off by default); no cursor
pagination — list endpoints return named arrays, the admin list uses `?limit=&offset=`;
success `{ "data": ... }` / error `{ "error": { code, message, details } }`.

## 6. The agent roster

Specialized agents in `.claude/agents/` encode these conventions so modules advance in
parallel without drift. One agent owns one module slice at a time; the shared seams
(`module.go`, `platform/`, migrations, `06_API_CONTRACTS.md`) are serialization points.

| Agent | Use it when… |
|---|---|
| **backend-module-builder** | Adding/extending a Go bounded context as a full vertical slice (domain → application → infra → api → `module.go` wiring + migration + service tests). |
| **frontend-feature-builder** | Building an App Router route + components for a backend module against the documented contract, with Vitest/Playwright coverage. |
| **api-contract-guardian** | Reviewing any endpoint change for drift across `06_API_CONTRACTS.md`, Go handlers, and the frontend client — run before merging cross-stack work. |
| **test-coverage-engineer** | Closing coverage gaps or adding unit/integration/e2e tests to the ≥70% (application+domain) bar; wiring coverage gating. |
| **security-auth-reviewer** | Any change touching auth, sessions, RBAC, OAuth, CSRF, password/token handling, or PII exposure. Required gate for such slices. |
| **devops-release-engineer** | Docker/compose, GitHub Actions, Helm, Prometheus/Grafana, and release/deploy checklists — i.e. when a feature adds infra. |

These complement the installed **engineering** plugin skills (`system-design`,
`architecture`, `code-review`, `testing-strategy`, `debug`, `deploy-checklist`,
`documentation`, `tech-debt`), which agents can invoke for cross-cutting work. A typical
feature flows: architecture (optional ADR) → backend-module-builder → api-contract-guardian
→ frontend-feature-builder → security-auth-reviewer → test-coverage-engineer →
devops-release-engineer. See `docs/12_AGENTS_AND_SCALING.md` for the full strategy.

## 7. Definition of Done (every slice)

- Code matches the DDD layering and the existing module's style.
- `docs/06_API_CONTRACTS.md` (and `openapi.yaml`) updated for any endpoint change.
- Backend: `gofmt` clean; `go build ./...` and `go test ./internal/<module>/...` green;
  application+domain coverage ≥70%.
- Frontend: `npm run lint`, `npx tsc --noEmit`, `npm run build` clean; touched components
  tested; critical journeys have a Playwright e2e.
- Auth-touching slices reviewed by **security-auth-reviewer**.
- New infra reflected in compose/Helm/observability with a stated rollback path.
- CI (`ci.yml`: lint/test/build) passes on the PR — agents do not bypass it.

## 8. Where to look first when extending a module

1. **`backend/internal/identity`** — the **gold standard**. Fully migrated DDD slice:
   `domain/ports.go` (interfaces), `domain/events.go`, `application/service.go` +
   `application/*_test.go` with `fakes_test.go`, `infrastructure/{postgres,crypto,jwtauth,
   mailer,oauth}/`, `api/{routes,handlers,dto,middleware}.go`, and `module.go` wiring.
   Mirror its structure, error handling, and test style.
2. **`backend/internal/mentorship`** — the **minimal reference** for a thin but correct
   slice (domain, application + tests, postgres repo, api, `module.go`) when identity is
   more than you need.
3. Match the closest existing module before writing; keep changes scoped to that module
   and report the files changed plus the build/test result.

## 9. gstack skills

The following `/gstack` skills are available in this project. Invoke them via the Skill
tool (e.g. `Skill("gstack", "/office-hours")`).

| Skill | Purpose |
|---|---|
| `/office-hours` | Open-ended Q&A and architectural discussion |
| `/plan-ceo-review` | Review a plan from a CEO / product-strategy lens |
| `/plan-eng-review` | Review a plan from an engineering-lead lens |
| `/plan-design-review` | Review a plan from a design / UX lens |
| `/review` | Code review of the current diff or a named PR |
| `/ship` | Pre-ship checklist: quality, security, docs, observability |
| `/qa` | QA pass — test plan generation and exploratory testing |
| `/investigate` | Root-cause investigation for a bug or incident |
| `/cso` | Chief Strategy Officer: market, competitive, and roadmap framing |
| `/autoplan` | Auto-generate an implementation plan from a task description |
| `/retro` | Retrospective facilitation for a completed sprint or feature |
| `/document-release` | Draft release notes and changelog entries |
| `/design-shotgun` | Rapid parallel design exploration (multiple concepts at once) |
| `/design-html` | Generate an HTML/CSS prototype from a design description |
| `/careful` | Enable extra-cautious mode before touching risky code |
| `/guard` | Add guard rails / validation to a code path |
| `/freeze` | Mark a module or file as frozen (no edits without review) |
| `/unfreeze` | Remove a freeze marker |
| `/learn` | Explain a concept or codebase area in depth |
| `/gstack-upgrade` | Upgrade the gstack skill suite to the latest version |

**Web browsing rule:** always use `/browse` from gstack for any web browsing task.
Never call `mcp__claude-in-chrome__*` tools directly.

## Skill routing

When the user's request matches an available skill, invoke it via the Skill tool. When in doubt, invoke the skill.

Key routing rules:
- Product ideas/brainstorming → invoke /office-hours
- Strategy/scope → invoke /plan-ceo-review
- Architecture → invoke /plan-eng-review
- Design system/plan review → invoke /design-consultation or /plan-design-review
- Full review pipeline → invoke /autoplan
- Bugs/errors → invoke /investigate
- QA/testing site behavior → invoke /qa or /qa-only
- Code review/diff check → invoke /review
- Visual polish → invoke /design-review
- Ship/deploy/PR → invoke /ship or /land-and-deploy
- Save progress → invoke /context-save
- Resume context → invoke /context-restore
- Author a backlog-ready spec/issue → invoke /spec
