# Running Kirmya with Docker

The whole stack — PostgreSQL, the Go API, and the Next.js frontend — runs with one command via [docker-compose.yml](../docker-compose.yml).

## Prerequisites
- Docker Engine 24+ with the Compose plugin (`docker compose`).

## Run

```bash
# from the repo root
docker compose up --build
```

Then open:
- **Frontend:** http://localhost:3000
- **API health:** http://localhost:8080/api/v1/health
- **Swagger UI:** http://localhost:8080/swagger-ui/

Migrations run automatically on API startup, and the six default communities are seeded. Register an account at `/register` to start.

## Demo data

`docker-compose.yml` sets `SEED_DEMO_DATA: "true"`, so a fresh database is populated with explorable demo content on first boot (idempotent — it does nothing if the demo admin already exists): demo users (job seekers, a referrer, mentors, a recruiter, and an admin), profiles + skills, jobs, mentor profiles, community posts/comments, and a sample referral.

Log in with any demo email and the password **`Password123!`**:

| Email | Role |
|---|---|
| `admin@demo.kirmya.io` | admin (see the **Admin** console + nav link) |
| `asha.rao@demo.kirmya.io` | job seeker |
| `rita.shah@demo.kirmya.io` | recruiter (posted the demo jobs) |
| `carla.mendes@demo.kirmya.io` | referrer |
| `deepa.nair@demo.kirmya.io` · `omar.farouk@demo.kirmya.io` | mentors |

Set `SEED_DEMO_DATA` to anything other than `true` (or remove it) to disable seeding.

## AI features (optional)
The AI module (resume reviewer, career coach, skill-gap) calls Claude. Provide a key before starting:

```bash
export ANTHROPIC_API_KEY=sk-ant-...   # PowerShell: $env:ANTHROPIC_API_KEY="sk-ant-..."
docker compose up --build
```

Without a key, AI endpoints return `503` and the rest of the app works normally.

## Services

| Service | Image / build | Port | Notes |
|---|---|---|---|
| `postgres` | `postgres:16-alpine` | 5432 | Data persisted in the `pgdata` volume; healthcheck-gated |
| `redis` | `redis:7-alpine` | 6379 | Cache-aside for profiles/jobs; in-memory only (persistence off). Optional — the API degrades to a no-op cache if absent |
| `opensearch` | `opensearchproject/opensearch:2` | 9200 | Full-text search over users/jobs/communities/skills. Optional — the API falls back to PostgreSQL `ILIKE` if absent |
| `nats` | `nats:2.10-alpine` (`-js`) | 4222, 8222 | JetStream event bus (cross-replica domain events). Optional — the API falls back to the in-process bus if absent |
| `backend` | `./backend/Dockerfile` (multi-stage Go → alpine) | 8080 | Reads `DATABASE_URL` + `REDIS_URL` + `OPENSEARCH_URL` + `NATS_URL` + `OTEL_EXPORTER_OTLP_ENDPOINT`; exposes Prometheus metrics at `/metrics`; uploads persisted in the `uploads` volume |
| `frontend` | `./frontend/Dockerfile` (Next standalone) | 3000 | Proxies `/api/v1/*` → `http://backend:8080` (baked via `API_PROXY_TARGET` build arg) |
| `jaeger` | `jaegertracing/all-in-one` | 16686 (UI), 4318 (OTLP) | Receives OpenTelemetry traces from the backend |
| `prometheus` | `prom/prometheus` | 9090 | Scrapes `backend:8080/metrics` every 10s |
| `grafana` | `grafana/grafana` | 3001 | Pre-provisioned Prometheus datasource + "Kirmya — API Overview" dashboard |

## Security middleware

Every response passes through a security pipeline (in `internal/platform/middleware`):
- **Security headers** — `X-Content-Type-Options`, `X-Frame-Options: DENY`, `Referrer-Policy`, `Permissions-Policy`, HSTS, and a path-aware Content-Security-Policy (locked down for the API, relaxed for the Swagger UI page).
- **Rate limiting** — per-client-IP token bucket; tune with `RATE_LIMIT_RPS` (default 50) and `RATE_LIMIT_BURST` (default 100). Health and metrics are never limited; over-limit returns `429` with `Retry-After`.
- **CSRF** — the refresh-token cookie is `SameSite=Strict` and all other endpoints use Bearer tokens, so CSRF is covered by default. An extra Origin allowlist check is available opt-in via `CSRF_VERIFY_ORIGIN=true` (requires `APP_URL` to exactly match the browser origin).

RBAC: posting/managing jobs and viewing applicants requires the `recruiter` role; the admin console requires `admin`.

## Real-time notifications (SSE)

`GET /api/v1/notifications/stream` is a Server-Sent Events endpoint that pushes
notifications to the authenticated user in real time (the bell updates instantly
when a referral, message, or mentorship event fires). The frontend connects with
`fetch` (so the Bearer token rides the Authorization header), auto-reconnects on
drop, and a slow background poll reconciles state.

**Multi-instance:** when `NATS_URL` is set, the hubs broadcast each event over a
core-NATS fanout subject (`sse.notifications` / `sse.messages`), so whichever
backend instance holds a user's SSE connection delivers it — not only the
instance that produced the event. Without NATS they fall back to single-instance
in-memory delivery.

Live chat uses the same pattern (and the same NATS fanout): `GET /api/v1/conversations/stream`
multiplexes **message**, **typing**, and **read** events across the user's
conversations. The Messages page appends new messages instantly (deduped by id),
shows a "typing…" indicator (throttled `POST /conversations/{id}/typing`), and
shows "Seen" read receipts when the other participant reads the thread.

## Frontend unit / component tests (Vitest)

Component and unit tests run in jsdom with Vitest + React Testing Library — no
backend needed:

```bash
cd frontend
npm test            # vitest run (CI mode)
npm run test:watch  # watch mode
```

They cover the API client (envelope handling + transparent refresh-on-401), the
auth context (login/logout/MFA), and the login form. CI runs them in the
`frontend` job.

## End-to-end tests (Playwright)

A lean stack (`docker-compose.e2e.yml` — Postgres + API + frontend, demo data
seeded, heavy deps omitted so the app runs on its graceful fallbacks) backs the
Playwright suite:

```bash
docker compose -f docker-compose.e2e.yml up -d --build
cd frontend
npm ci && npx playwright install --with-deps chromium
npm run e2e          # runs e2e/*.spec.ts against http://localhost:3000
docker compose -f docker-compose.e2e.yml down -v
```

The suite covers the core journeys: register/login/logout, the admin console
(and non-admin redirect), search, jobs, and referrals — all asserting against the
seeded demo content. CI runs it as the `e2e` job in `.github/workflows/ci.yml`.

## Repository / integration tests

Repository-layer tests run against a real PostgreSQL and are gated behind the
`integration` build tag (so the normal `go test ./...` stays fast and DB-free).
They skip unless `TEST_DATABASE_URL` is set:

```bash
docker run -d --name cb-itest -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=kirmya -p 5546:5432 postgres:16-alpine
cd backend
TEST_DATABASE_URL='postgres://postgres:postgres@localhost:5546/kirmya?sslmode=disable' \
  go test -tags=integration -p 1 ./...
```

`-p 1` serializes the package test binaries so they don't race while applying
migrations. CI runs these in the `integration` job (`.github/workflows/ci.yml`)
against a Postgres service container.

## Observability

Once the stack is up, open:
- **Grafana dashboards:** http://localhost:3001 (anonymous viewer; admin/admin to edit) — RED metrics per route, latency p95, in-flight requests, DB pool, cache hit rate.
- **Prometheus:** http://localhost:9090 — raw metrics and queries (`/metrics` on the API is at http://localhost:8080/metrics).
- **Jaeger traces:** http://localhost:16686 — per-request spans (select service `kirmya`).

Tracing is gated by `OTEL_EXPORTER_OTLP_ENDPOINT`; unset it to disable. Metrics are always exposed at `/metrics`.

## How the proxy / cookies work
The browser only ever talks to the frontend origin (`:3000`). Next's `rewrites` forward `/api/v1/*` to the backend service, so the auth **refresh cookie stays same-origin** — no CORS, no cross-site cookie issues. In local non-Docker dev the same rewrite defaults to `http://localhost:8080`.

## Common commands

```bash
docker compose up --build        # build + start
docker compose up -d             # start detached
docker compose logs -f backend   # tail API logs (verification email links print here in dev)
docker compose down              # stop
docker compose down -v           # stop and wipe the database + uploads volumes
```

## Configuration (backend env)
Set via `docker-compose.yml` (or override in your shell). Key vars: `DATABASE_URL`, `JWT_SECRET` (change for production), `JWT_ACCESS_TTL`, `JWT_REFRESH_TTL`, `APP_ENV` (`production` enables Secure cookies), `APP_URL`, `RESUME_UPLOAD_DIR`, `ANTHROPIC_API_KEY`, and OAuth `GOOGLE_*` / `LINKEDIN_*`. See [backend/.env.example](../backend/.env.example) for the full list.
