# Deploying Kirmya to a Hostinger VPS (Docker)

This targets a **Hostinger VPS** (or any single Docker host). It uses
`docker-compose.prod.yml`, which runs PostgreSQL + the Go API + the Next.js frontend
behind a **Caddy** reverse proxy that handles automatic HTTPS. Redis/OpenSearch/NATS
are omitted — the app degrades gracefully without them — and observability is optional.

> Hostinger **shared/web hosting** cannot run this (no Docker, no long-running Go/Node
> processes). You need the **VPS** product (or Hostinger's Docker/VPS template).

## Prerequisites

- A Hostinger VPS with Docker + the Compose plugin. On a fresh Ubuntu VPS:
  ```bash
  curl -fsSL https://get.docker.com | sh
  ```
- A domain's **A record** pointing at the VPS public IP.
- Ports **80** and **443** open in the VPS firewall (Caddy needs both for ACME + TLS).

## 1. Get the code and configure secrets

```bash
git clone <your-repo-url> kirmya && cd kirmya
cp .env.prod.example .env
nano .env        # fill in DOMAIN, ACME_EMAIL, POSTGRES_PASSWORD, JWT_SECRET, MFA_ENC_KEY, SMTP_*
```

Generate the two required keys:

```bash
openssl rand -hex 32   # JWT_SECRET
openssl rand -hex 32   # MFA_ENC_KEY
```

`JWT_SECRET` and `MFA_ENC_KEY` are **mandatory in production** — the backend refuses
to start without them. Leave `EMAIL_VERIFICATION_REQUIRED=true` only if SMTP is set;
otherwise set it `false` until you configure a mailer (production fails closed on the
log-only mailer).

## 2. Launch

```bash
docker compose -f docker-compose.prod.yml up -d --build
```

What happens: Postgres starts and becomes healthy → backend connects, **runs
migrations automatically**, and starts on the internal network → frontend builds with
`API_PROXY_TARGET=http://backend:8080` and starts → Caddy obtains a Let's Encrypt
certificate for `DOMAIN` and serves the site on 443.

Check status and logs:

```bash
docker compose -f docker-compose.prod.yml ps
docker compose -f docker-compose.prod.yml logs -f caddy backend
```

Visit `https://<DOMAIN>`. The backend health endpoint is reachable internally; from the
host you can verify with `docker compose -f docker-compose.prod.yml exec backend wget -qO- http://127.0.0.1:8080/api/v1/health`.

## 3. Updates / redeploys

```bash
git pull
docker compose -f docker-compose.prod.yml up -d --build
```

Because the frontend bakes the API proxy target at build time, always rebuild
(`--build`) rather than just restarting. Migrations are forward-only and re-applied
idempotently on each backend boot.

## 4. Optional: observability

```bash
docker compose -f docker-compose.prod.yml --profile observability up -d
```

Starts Prometheus + Grafana (set `GRAFANA_ADMIN_PASSWORD` in `.env` first). Neither is
published publicly by default — reach them via an SSH tunnel, e.g.
`ssh -L 3000:localhost:3000 user@vps` after adding a temporary port mapping, or put
them behind Caddy with auth.

## Operational notes

- **Backups:** the database lives in the `pgdata` volume and uploads in `uploads`.
  Back them up regularly, e.g.
  `docker compose -f docker-compose.prod.yml exec postgres pg_dump -U $POSTGRES_USER $POSTGRES_DB > backup.sql`.
- **Only Caddy is exposed.** Postgres, backend, and frontend have no published ports;
  they talk over the private compose network. Don't add `ports:` to them in prod.
- **`restart: unless-stopped`** on every long-running service means the stack comes
  back automatically after a reboot or crash.
- **Secrets** come from `.env` (gitignored). Rotate `JWT_SECRET`/`MFA_ENC_KEY` only
  with care — rotating `JWT_SECRET` invalidates existing sessions; rotating
  `MFA_ENC_KEY` invalidates stored TOTP secrets (users must re-enrol MFA).
- **`TRUST_PROXY=true`** is set so the app reads the real client IP from Caddy's
  `X-Forwarded-For` for audit logs.
