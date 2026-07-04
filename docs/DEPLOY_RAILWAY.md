# Deploying Kirmya to Railway

Kirmya runs on Railway as **two services from this one repo** plus a **Railway
PostgreSQL** plugin. Redis, OpenSearch, and NATS are optional — the app degrades
gracefully without them (no-op cache, DB-fallback search, in-process event bus), so
you do **not** need them to go live. Migrations run automatically on backend startup,
so there is no separate migration step.

## Overview

| Railway service | Source | Root directory | Builds from |
|---|---|---|---|
| `backend` | this repo | `backend` | `backend/Dockerfile` + `backend/railway.json` |
| `frontend` | this repo | `frontend` | `frontend/Dockerfile` + `frontend/railway.json` |
| `Postgres` | Railway plugin | — | managed |

Each `railway.json` pins the Dockerfile builder, the health-check path, and an
on-failure restart policy, so the only manual setup is the root directory and the
environment variables below.

## 1. Create the project and database

1. New Project → Deploy from GitHub repo → select this repo.
2. Add a **PostgreSQL** database (New → Database → PostgreSQL). Railway exposes a
   `DATABASE_URL` reference variable you will wire into the backend.

## 2. Backend service

- **Settings → Root Directory:** `backend` (so `railway.json` and the Dockerfile are found).
- **Variables:**

| Variable | Value | Notes |
|---|---|---|
| `DATABASE_URL` | `${{Postgres.DATABASE_URL}}` | Reference the Postgres plugin. Railway's URL already includes SSL. |
| `JWT_SECRET` | _(generate: `openssl rand -hex 32`)_ | **Required.** Never reuse the dev value. |
| `MFA_ENC_KEY` | _(generate: `openssl rand -hex 32`)_ | **Required in production** — the app aborts on boot without it (falls back to `JWT_SECRET` only in dev). |
| `TRUST_PROXY` | `true` | Railway terminates TLS in front of the app; needed for correct client IPs in audit logs. |
| `APP_ENV` | `production` | Enables Secure cookies + fail-closed startup checks. |
| `APP_URL` | `https://<your-frontend-domain>` | Public URL of the frontend; used in emails/links and CORS. |
| `EMAIL_VERIFICATION_REQUIRED` | `true` | Set `false` only if you have no mailer configured yet. |
| `SEED_DEMO_DATA` | `false` | Leave off in production. |
| `JWT_ACCESS_TTL` | `900` | Optional (default shown). |
| `JWT_REFRESH_TTL` | `2592000` | Optional. |
| `ANTHROPIC_API_KEY` | _(your key)_ | Optional — enables AI features. |
| SMTP vars | see `backend/.env.example` | Required if `EMAIL_VERIFICATION_REQUIRED=true`. |

> Do **not** set `PORT` — Railway injects it and the server already honors `$PORT`.

Railway reads the health-check path (`/api/v1/health`) from `backend/railway.json`.

## 3. Frontend service

- **Settings → Root Directory:** `frontend`.
- **Variable — set as a _build_ variable** (this is consumed at build time; see note):

| Variable | Value |
|---|---|
| `API_PROXY_TARGET` | `http://${{backend.RAILWAY_PRIVATE_DOMAIN}}:8080` |

Or, if you prefer the public URL: `https://<your-backend-domain>`.

> **Why build-time:** `next.config.ts` bakes the `/api/v1` rewrite destination into
> the route manifest at build, so `API_PROXY_TARGET` must be present **when the image
> builds**, not just at runtime. Railway passes service variables to the Docker build
> as `ARG`s (the Dockerfile already declares `ARG API_PROXY_TARGET`), so setting it as
> a normal service variable works — but you must **redeploy/rebuild** the frontend if
> the backend's URL ever changes.

## 4. Networking & domains

- Generate a public domain for **frontend** (Settings → Networking → Generate Domain).
- The backend does not need a public domain if the frontend reaches it over the
  private network (`*.railway.internal`). Give it one only if you want the API
  publicly reachable; if you do, set `API_PROXY_TARGET` to that public URL instead.
- Set the backend's `APP_URL` to the frontend's public domain.

## 5. Verify

- Backend: `https://<backend-domain>/api/v1/health` (if public) or the service's
  Railway health check should be green.
- Frontend: open the public domain; the app should load and API calls under
  `/api/v1/...` should succeed (they are proxied to the backend).

## Notes / gotchas

- **Private port:** Railway's private networking targets the port your service
  listens on. The backend listens on `$PORT`; the `:8080` in `API_PROXY_TARGET`
  above assumes the backend image's default. If you override `PORT` on the backend,
  match it in `API_PROXY_TARGET`.
- **No down-migrations:** migrations are forward-only and run on boot; a bad deploy
  rolls back the *image*, not the schema.
- **Optional infra:** to add Redis later, add the Railway Redis plugin and set
  `REDIS_URL`; the cache layer picks it up automatically.
