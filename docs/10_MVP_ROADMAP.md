# Kirmya — MVP Roadmap & Production Strategy

> Status: Draft v1 · Last updated: 2026-06-14

## 1. Phased Roadmap

### Phase 0 — Foundation (re-platform) ✦ in progress
- PostgreSQL platform layer (config, pgx pool, migration runner) replacing MySQL.
- Middleware pipeline (request-id, recover, security headers, CORS, rate-limit, CSRF, auth, RBAC, OTel).
- In-process event bus + transactional outbox.
- Docs (this set). **Deliverable:** boots, health check, migrations run on Postgres.

### Phase 1 — Identity & Profiles
- Identity module (DDD layout): register, login, logout, email verification, password reset, refresh-token rotation, Google + LinkedIn OAuth, MFA-ready (TOTP), RBAC, Argon2id hashing, audit logging.
- Profile module: profile + experience/education/certifications/skills/languages/portfolio.
- Frontend: auth flows + dashboard shell + profile pages.
- **Deliverable:** a user can sign up, verify, log in (incl. OAuth), and complete a profile.

### Phase 2 — Resume, Jobs, Referrals
- Resume upload/parse/score/versions + improvement suggestions.
- Jobs: post/search/save/apply/track + AI matching (basic).
- Referral marketplace with full state machine.
- Search (OpenSearch) for jobs + users.
- **Deliverable:** the core recovery loop — upload resume, find jobs, request referrals, apply.

### Phase 3 — Career Intelligence, Communities, Mentorship
- AI skill-gap engine + career paths + salary/market insights + learning paths.
- Communities: posts/comments/reactions/polls/tags/moderation.
- Mentorship: profiles/booking/reviews.
- **Deliverable:** guidance + connection layers live.

### Phase 4 — Messaging, Notifications, Admin, Hardening
- Real-time messaging (WS) + notifications.
- Admin console (users, moderation, reports, analytics).
- Security review (OWASP), load testing, observability dashboards.
- **Deliverable:** production-ready MVP.

## 2. Testing Strategy
- **Backend:** unit (domain/application), repository/integration (dockerized Postgres), API (`httptest`). Target ≥70% on application+domain.
- **Frontend:** component (Vitest + RTL), E2E (Playwright) for critical journeys.
- **CI gate:** lint + unit + build must pass on every PR (GitHub Actions).

## 3. DevOps & Infrastructure
- **Dockerfiles:** multi-stage for backend (distroless final) and frontend (standalone Next build).
- **Docker Compose (local):** api, frontend, postgres, redis, opensearch, (later) nats, otel-collector, grafana, prometheus.
- **GitHub Actions:** `ci.yml` (lint/test/build), `docker.yml` (build+push images), `deploy.yml` (env promotion).
- **Kubernetes + Helm:** chart with deployments (api, frontend), services, ingress, HPA, configmaps/secrets, postgres (managed/operator), redis, opensearch. Probes: `/api/v1/health` (liveness/readiness).
- **Environments:** dev → staging → prod; blue/green or rolling with readiness gates.

## 4. Production Deployment Strategy
- Stateless API behind ingress + HPA (CPU/RPS). Managed Postgres (primary + read replicas), managed Redis, managed OpenSearch.
- Secrets via cloud secret manager / sealed-secrets; never in images.
- Migrations run as a pre-deploy Job (gated, backward-compatible / expand-contract).
- Multi-region readiness: stateless API per region, primary Postgres + cross-region replicas, Redis per region, global LB. Messaging/AI extractable to dedicated services first.
- Rollback: keep previous image + `down` migrations only when safe; prefer forward fixes. Document rollback triggers per release.

## 5. Observability
- **Tracing:** OpenTelemetry SDK → collector → backend (Tempo/Jaeger).
- **Metrics:** Prometheus (RED metrics per route, DB pool, event bus lag, cache hit rate).
- **Logs:** structured JSON with trace/span IDs.
- **Dashboards/alerts:** Grafana — latency p95/p99, error rate, saturation; alert on SLO burn.

## 6. Scalability Path to 10M users / 100M messages
- Read replicas + Redis cache-aside for hot reads (profiles/jobs).
- Partition high-volume tables (messages/notifications/audit) by time.
- Extract Messaging → dedicated service + queue when WS/throughput demands.
- Move event bus from in-process → NATS JetStream for durability + cross-service.
- OpenSearch scaled independently for search load.
