# Kirmya — Full Project Audit

**Date:** 2026-06-21
**Method:** Six parallel review agents (backend architecture, frontend, security/auth, test coverage, DevOps/CI, API-contract drift) run read-only over the repo, plus a git-state investigation. Findings consolidated here.
**Important caveat:** The Go toolchain was **not available** in the audit environment, so `go build`, `go vet`, `gofmt -l`, and `go test -cover` could **not be executed**. Backend and coverage findings are static estimates and must be confirmed by running the suite on a Go 1.26 machine. The frontend `tsc --noEmit` did run and passed.

---

## 1. Executive summary

Kirmya is a well-architected, security-conscious codebase that is further along on engineering discipline than most projects of its size. The backend DDD migration is essentially complete (all 13 bounded contexts on the target layout), the identity/auth module is genuinely strong, the API contract is in near-perfect sync across code, the human-readable doc, and the frontend client, and the infrastructure (compose, Helm, Caddy, CI) is unusually complete.

The problems are concentrated in a few specific places: **a known default secret shipped in the Helm chart**, **observability config that is referenced but missing**, **a rate limiter that can be bypassed**, **CI that does not actually enforce the project's own Definition of Done**, **documentation that has drifted from reality**, and **a git index that has been accidentally wiped from cache**. None of these are deep design flaws; all are fixable without rearchitecting.

### Health scoreboard

| Domain | Grade | One-line |
|---|---|---|
| Backend architecture | A− | Clean DDD across all modules; a few convention slips (nil-safe events, locking, one cross-module import). |
| Security / auth | A− | Excellent crypto & RBAC; one High (rate-limiter XFF) + hardening items. |
| API contract | A | Code ↔ doc ↔ frontend perfectly aligned; only `openapi.yaml` lags by 10 endpoints. |
| Frontend | B | Solid API client & auth; but every page is a client component and there are no error/loading boundaries. |
| Test coverage | B− | Every module has application tests; domain layer & frontend e2e have real gaps. (Bar unverified — toolchain absent.) |
| DevOps / CI | B− | Very complete infra, but a critical default secret, missing observability files, and CI that doesn't gate the DoD. |
| Repo hygiene | C | Git index wiped from cache; dead files; doc drift. |

---

## 2. Top priorities (do these first)

1. **Fix the git index** — `git reset` (mixed) to restore staging from HEAD. Non-destructive. (§9)
2. **Remove the default `jwtSecret` from Helm `values.yaml`** — a `helm install` with no overrides yields a *known* JWT signing key = full auth bypass. (§7, Critical-1)
3. **Gate the X-Forwarded-For rate limiter behind `TRUST_PROXY`** — today the only throttle on `/auth/login` password attempts can be bypassed by rotating a spoofed header. (§5, High)
4. **Add the missing Grafana provisioning/dashboards** (or drop the volume mounts) so the observability profile actually works. (§7, Critical-2)
5. **Make CI enforce the Definition of Done** — add `tsc --noEmit`, a `gofmt -l` check, and the ≥70% coverage gate; ensure branch protection marks these as required. (§7, High)
6. **Sync `openapi.yaml`** — add the 10 implemented endpoints it's missing, and extend the CI contract guard to diff against it. (§6)
7. **Fix the doc drift in `CLAUDE.md`** — the frontend `src/features/` layout and the `{ data }` success envelope it describes are both wrong. (§8 + new agent, §10)

---

## 3. Backend architecture & Go code health  — Grade A−

**Strengths.** All 13 bounded contexts use the full DDD layout (api/application/domain/infrastructure + `module.go`); no legacy `service/repository/handler` layouts remain. Domain imports nothing from other layers. `context.Context` is consistently the first argument. Sentinel errors map to HTTP correctly. Optimistic locking is implemented in profile, settings, resume, referrals, and identity. No stray `TODO`/`FIXME`/`panic` in production code (the `log.Fatalf` guards for missing prod secrets are intentional fail-fast).

**Issues.**

- **High — Unguarded event publishing in `identity`.** `internal/identity/application/{service.go,auth_flows.go,oauth.go}` call `s.events.Publish(...)` directly with no nil guard, and `NewService` assigns `events` with no no-op default. This violates the project's own best-effort/nil-safe convention — and identity is the *gold standard* others copy. Fix: add `if s.events != nil` (mirror the `publish()` helper the other 8 modules already use) or default to a no-op publisher.
- **Medium — Missing migration `004`.** `migrations/` jumps `003_create_jobs_table.sql` → `005_create_profile_tables.sql`. With forward-only sequential application this is a permanent gap. Confirm it was intentionally squashed and document it, or backfill a no-op `004`.
- **Medium — No optimistic locking on `jobs` and `mentorship`.** Both have mutable aggregates (job edits, application-status transitions, booking state) but no `Version` field / version-checked update. Add it per convention, or explicitly note these are last-writer-wins.
- **Medium — Cross-module import violation.** `internal/jobs/infrastructure/aimatch/matcher.go` imports `ai/domain` *and* `ai/infrastructure/anthropic` (a concrete adapter of another module). Convention is to depend only on another module's `application.Service`. Define a jobs-owned port and inject the ai adapter in `main.go`.
- **Low — Dead `internal/domain/models.go`** (top-level `User`/`Job`/`Workspace` structs, imported nowhere). Delete it.
- **Low — `search` has no `domain/` folder** (api/application/infra only) — probably fine for a thin slice; note vs. documented layout.

**Build/vet/gofmt:** not run (toolchain absent). Nothing static suggests a compile break; run `cd backend && gofmt -l ./... && go build ./... && go vet ./...` to confirm.

---

## 4. Frontend / Next.js  — Grade B

**Strengths.** `lib/api.ts` is a strong, fully-typed single-source client: in-memory access token mirrored to localStorage, transparent refresh-and-retry exactly once with correct loop guards, `credentials: "include"` for the httpOnly refresh cookie, FormData-aware content type, SSE via `fetch` so the Bearer token rides the header. `lib/auth-context.tsx` is clean. The `messages` SSE page handles abort/reconnect with a ref to dodge stale closures. No `any`, no stray `console.*`. `tsc --noEmit` passes clean; lint passes with only intentional warnings.

**Issues.**

- **High — Every one of the 19 pages is `"use client"`**, including the public marketing landing `app/page.tsx`. This abandons the documented "Server Components by default" principle and ships unnecessary JS / loses SSR & SEO on public pages. Make static/marketing routes Server Components and isolate auth-dependent bits into small client children.
- **Medium — No App Router error/loading/not-found boundaries.** Zero `error.tsx`, `loading.tsx`, `not-found.tsx`, or `global-error.tsx` under `app/`; each page hand-rolls inline states. A thrown render error has no boundary. Add at least root `app/error.tsx` and `app/not-found.tsx`.
- **Medium — Auth gating is client-only and duplicated.** `components/RequireAuth.tsx` and an inline guard in `app/admin/page.tsx` both redirect via `useEffect` after hydration (protected content briefly mounts before redirect). Factor a `RequireRole`/`RequireAuth(role)` wrapper; consider `middleware.ts` for a server-side redirect. (Server-side RBAC *is* enforced — see §5 — so this is UX/defense-in-depth, not an authz hole.)
- **Medium — Docs-vs-reality:** `CLAUDE.md` documents a `frontend/src/features/<f>/` layout that **does not exist** — the app is flat `app/` + `components/` + `lib/`, with the whole API surface in one 742-line `lib/api.ts`. Update the doc or track the migration.
- **Low — Dual design systems.** `package.json` pulls in `@mui/material` + `@emotion/*` **and** Tailwind v4 + shadcn `components/ui` + lucide. Audited routes use only Tailwind/shadcn. Confirm MUI is actually used; if not, remove it to cut bundle size.
- **Low — three `react-hooks/set-state-in-effect` warnings** (login prefill, GlobalSearch, NotificationBell) — deliberately downgraded; idiomatic fetch-on-mount/debounce. Low priority.

**tsc/lint:** `tsc --noEmit` PASS (0 diagnostics). Scoped eslint PASS (0 errors / 3 intentional warnings). `npm run build` not run (time budget).

---

## 5. Security & auth  — Grade A−

**Strengths.** Argon2id at OWASP baseline (64 MiB, t=3, p=2) with constant-time verify. JWT alg-confusion blocked (non-HMAC rejected; HS256 secret ≥32 bytes enforced, process aborts in prod if unset). Refresh rotation with reuse-revokes-family; refresh tokens stored as SHA-256 hashes, never raw; password reset/change revoke all sessions. OAuth login-CSRF defended via cookie-bound, constant-time, single-use `state`. **RBAC enforced server-side** — all `/admin/*` routes wrapped by `RequireRole(RoleAdmin)`, self-registration role escalation blocked by an allowlist. DTO hygiene good (no password hash exposed; directory omits email). Refresh cookie is HttpOnly + SameSite=Strict + path-scoped + Secure. Parameterized SQL throughout.

**Issues.**

- **High — Rate limiter unconditionally trusts `X-Forwarded-For`.** `platform/middleware/ratelimit.go` reads the client IP from XFF with no `TRUST_PROXY` gate — inconsistent with the identity handler's `clientIP`, which gates it correctly. This global limiter is the *only* throttle on `/auth/login` password attempts (the attempt-limiter only covers the MFA/TOTP step). An attacker rotating a spoofed XFF per request bypasses it entirely. Fix: gate XFF behind `TRUST_PROXY`, fall back to `RemoteAddr`, and consider a per-account failed-password limiter.
- **Medium — No CORS policy.** The middleware chain has no CORS layer; the model depends on same-origin via the proxy. Add an explicit allowlist CORS keyed on `APP_URL` (never `*` with credentials), or document the same-origin requirement.
- **Medium — OAuth: no PKCE, ID token not validated.** Plain auth-code exchange; userinfo email trusted without checking `email_verified`/`iss`/`aud`. Add PKCE (both providers support it) and validate the OIDC ID token before auto-provisioning accounts.
- **Medium — Access-token revocation lag** (documented). Stateless JWT means a suspended user stays valid until expiry (≤15 min). Durable fix: a `token_version`/revocation check.
- **Low — TOTP same-window replay** (no persisted last-used counter); **Low — `csrf_token` cookie is unenforced scaffolding**; **Low — duplicated `clientIP` helpers** (the root cause of the High finding) — consolidate into `common`.

**Committed secrets:** **None committed.** `backend/.env` and `frontend/.env.local` exist on disk and `backend/.env` holds real-looking LinkedIn OAuth + SMTP + JWT values, but neither file is git-tracked and they never appear in history; `.gitignore` excludes them. Flag (operational, not a repo exposure): rotate those live dev credentials before any production use and never copy them into a tracked file or CI artifact.

---

## 6. API contract drift  — Grade A

**Alignment:** Go routes (120) ↔ `docs/06_API_CONTRACTS.md` (120) = **perfect**. Go routes ↔ `frontend/lib/api.ts` (89 calls) = **perfect** (no orphan calls). Envelope shapes consistent across code/frontend/openapi. The one real drift is `openapi.yaml`.

- **High — 10 implemented endpoints missing from `backend/docs/openapi.yaml`** (the machine source of truth), so Swagger UI / codegen consumers can't see them. They exist in Go routes *and* in the markdown doc. Add them to openapi:
  `GET /communities/{slug}/tags`, `GET /communities/{slug}/reports`, `DELETE /communities/{slug}/posts/{id}`, `POST /posts/{id}/polls`, `POST /posts/{id}/report`, `GET /polls/{id}`, `POST /polls/{id}/vote`, `POST /mentorship/availability`, `GET /mentorship/mentors/{id}/reviews`, `GET /mentorship/mentors/{id}/availability`.
- **Low — error-envelope `details` field** is declared in the frontend type and in `CLAUDE.md`/contract prose, but the backend `AppError` and openapi `ErrorEnvelope` are `{ code, message }` only. Either implement `details` or stop documenting it.
- **CI guard gap:** `scripts/check-api-contract.mjs` validates only Go `routes.go` ↔ the markdown doc (reports "OK — in sync"). It **never reads `openapi.yaml` or `frontend/lib/api.ts`**, which is exactly why the 10-endpoint openapi drift sails past CI. Extend it to diff openapi paths as a third set (and optionally the frontend client as a fourth).

---

## 7. DevOps / CI / infrastructure  — Grade B−

**Strengths.** Multi-stage Dockerfiles, static `CGO_ENABLED=0` Go binary with `-trimpath -ldflags="-s -w"`, non-root users, Next standalone output, good layer caching, both `.dockerignore` present. Compose dependency ordering uses `service_healthy`. Prod compose is hardened (`${VAR:?error}` fail-fast, no published app ports, observability behind a profile). Helm chart has probes, HPA, secrets in a `Secret`, Prometheus annotations. Caddyfile ships HSTS, nosniff, referrer-policy, auto-TLS. `check-api-contract.mjs` is well-engineered and correctly wired into `contract-check.yml`.

**Issues.**

- **Critical — Default `jwtSecret` in Helm `values.yaml`.** Ships `jwtSecret: "change-me-in-production-with-a-32-byte-minimum-secret"` (a usable 32-byte string) plus a default postgres password `postgres`. A no-override install = known signing key = auth bypass. Default to empty and `{{ required ... }}`-fail on render; document external secret management.
- **Critical — Grafana config referenced but missing.** `docker-compose.yml` and `docker-compose.prod.yml` bind-mount `./ops/grafana/provisioning` and `./ops/grafana/dashboards`, but only `ops/prometheus/prometheus.yml` exists. Grafana starts with no datasource/dashboards. Add the files or remove the mounts.
- **High — CI doesn't enforce the Definition of Done.** `ci.yml` runs Go vet/build/test + frontend lint/test/build, but **no `tsc --noEmit`, no `gofmt -l` check, no ≥70% coverage gate** (the `make backend-coverage` target is never invoked). And required-status enforcement lives in branch protection, which isn't visible in-repo. A PR can merge with type errors, unformatted Go, and no coverage.
- **High — `railway.json` documented but missing.** `docs/DEPLOY_RAILWAY.md` tells Railway to read `backend/railway.json` and `frontend/railway.json` (builder + `/api/v1/health` healthcheck). Neither exists; a Railway deploy falls back to defaults with no healthcheck. Add both or fix the doc.
- **Medium —** Floating base-image/action tags (no digest pinning; `air@latest` in dev stage); Helm dependency tiers (postgres/redis/opensearch/nats) have **no resource limits and only readiness, no liveness probes** (OpenSearch can OOM the node; a hung postgres won't restart); postgres password passed as plain env, not from the Secret; `sslmode=disable` everywhere; Caddyfile missing CSP / `X-Frame-Options` / Permissions-Policy.
- **Low —** the `docker` CI job builds with the dev override merged in (validates dev stages, not prod runtime images); heavy e2e/docker jobs aren't `needs:`-gated on unit tests; no Prometheus alerting; `DISABLE_SECURITY_PLUGIN=true` on the Helm OpenSearch with no NetworkPolicy.

**Referenced-but-missing files:** `ops/grafana/provisioning/`, `ops/grafana/dashboards/`, `backend/railway.json`, `frontend/railway.json`. (Note: `.env.prod.example` is present at repo root; Redis/OpenSearch/NATS absence from prod/e2e compose is intentional graceful degradation.)

---

## 8. Test coverage  — Grade B−  *(bar unverified — Go toolchain absent)*

**Strengths.** Every one of the 13 business modules has an application-layer service test following the gold-standard table-driven + `fakes_test.go` pattern. identity is the best-tested slice (10 files across application/domain/infra). referrals, resume, profile, ai, settings, notifications, mentorship, jobs all have solid application coverage. Postgres integration tests exist for several repos. Frontend has the API client and auth context under test.

**Gaps.**

- **Critical —** admin (RBAC / privileged mutations) has only ~4 application test funcs + 1 domain test — thin for the highest-privilege module. community/domain (membership, roles, moderation) has 3 source files and **zero** domain tests. No frontend e2e for register / forgot-password / reset-password / verify-email — security-sensitive recovery flows.
- **High —** 8 of 12 domain packages have **zero** unit tests (ai, community, jobs, mentorship, messaging, notifications, profile, settings). Since the ≥70% bar counts application **+ domain**, the domain half is likely failing for these. Frontend e2e missing for core mutation journeys (mentorship booking, messages send, communities join/post, resume upload, profile, settings).
- **Medium —** messaging lightly covered; search OpenSearch-vs-ILIKE fallback branch worth confirming; frontend component coverage sparse (only login + settings).

**Coverage numbers:** none obtained. Run `cd backend && go test ./... -cover` on Go 1.26 to verify the bar.

---

## 9. Repo hygiene / git state  — Grade C

- **High — Git index wiped from cache.** `master` HEAD is healthy (59 commits), but `git status` shows **368 staged deletions** of files that still exist on disk (and re-appear as untracked). This is the signature of a `git rm -r --cached .` / index reset: the working tree is intact, only the staging area is wrong. **Fix (non-destructive):** `git reset` (mixed reset to HEAD) restores the index and clears the bogus deletions; `git status` should then show only genuine changes. Review with `git status` afterward before committing.
- **Low — Dead/clutter files:** `internal/domain/models.go` (unused); `.claude/agents/principal-ai-architect.md` (self-marked "Superseded… Safe to delete"); `.claude/agents/CLAUDE-FABLE-5.md` (a stray model system prompt with no agent frontmatter — not a real agent definition).
- **Low — Documentation drift** (recurring across this audit): `CLAUDE.md` mis-describes the frontend layout (`src/features/`) and the success envelope (`{ data }` vs the real `{ success, data, meta }`); `openapi.yaml` lags code; deploy docs reference missing `railway.json`. This pattern is the basis for the new agent in §10.

---

## 10. Agent roster — gap assessment & new agent

**Existing roster** (`.claude/agents/`): backend-module-builder, frontend-feature-builder, api-contract-guardian, security-auth-reviewer, test-coverage-engineer, devops-release-engineer, principal-ai-ux-architect, sales-psychology-coach (+ two dead files noted in §9).

**Gap.** The single most common finding across all six audits was **documentation/spec drift from reality** — wrong layout and envelope in `CLAUDE.md`, 10 endpoints missing from `openapi.yaml`, and deploy docs referencing files that don't exist. No existing agent owns this: api-contract-guardian deliberately checks only `routes.go` ↔ the markdown contract (not openapi, not the docs corpus, not referenced-file existence), and the builder agents create features, not documentation truth.

**New agent created:** `.claude/agents/docs-sync-steward.md` — a read-only-plus-doc-edits agent that keeps `CLAUDE.md`, `docs/`, `backend/docs/openapi.yaml`, and inline references honest against the actual code: detects layout/convention claims that no longer match the tree, syncs openapi with implemented routes, and flags any doc that references a file which doesn't exist. It complements (does not overlap) api-contract-guardian, which remains the code↔contract route gate.

---

## 11. Appendix — what could not be verified

- Backend `go build` / `go vet` / `gofmt` / `go test -cover` — Go toolchain absent in the audit environment.
- `npm run build` and a full-tree eslint run — time budget (scoped runs passed; `tsc --noEmit` passed clean).
- Branch-protection / required-status configuration — lives in GitHub settings, not in-repo.
- Live behavior of the rate limiter, OAuth flows, and Helm rendering — assessed by code reading only.
