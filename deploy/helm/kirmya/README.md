# Kirmya — Helm chart

Deploys the Kirmya API (Go modular monolith) and Next.js frontend, with
optional in-cluster PostgreSQL, Redis, and OpenSearch.

## Quick start

```bash
# from repo root — build & push images first (or use a registry that has them)
helm upgrade --install cb deploy/helm/kirmya \
  --set secrets.jwtSecret=$(openssl rand -hex 32) \
  --set backend.image.repository=YOUR_REGISTRY/kirmya-backend \
  --set frontend.image.repository=YOUR_REGISTRY/kirmya-frontend
```

Then add an `/etc/hosts` entry for `kirmya.local` → your ingress IP and
open <http://kirmya.local>.

## How requests route

The Ingress sends `/api/v1/*` to the **backend** service and everything else to
the **frontend** service, so the browser talks to one origin (cookies work) and
the Next.js in-cluster proxy isn't needed.

## Dependencies

`postgres`, `redis`, and `opensearch` are bundled as single-replica Deployments
for convenience. **For production, disable them and point at managed services:**

```bash
helm upgrade --install cb deploy/helm/kirmya \
  --set postgres.enabled=false \
  --set secrets.databaseUrl='postgres://user:pass@my-rds:5432/kirmya?sslmode=require' \
  --set redis.enabled=false \        # app degrades to a no-op cache
  --set opensearch.enabled=false     # app degrades to DB-fallback search
```

The backend runs migrations automatically on startup. Redis/OpenSearch URLs are
injected automatically when those subcharts are enabled.

## Key values

| Key | Default | Notes |
|---|---|---|
| `backend.image.repository` / `.tag` | `kirmya/backend:latest` | push your built image |
| `frontend.image.repository` / `.tag` | `kirmya/frontend:latest` | |
| `backend.autoscaling.enabled` | `true` | CPU HPA 2–10 |
| `ingress.host` | `kirmya.local` | |
| `config.otelEndpoint` | `""` | set to an OTLP/HTTP collector to enable tracing |
| `secrets.jwtSecret` | `change-me-in-production` | **override this** |
| `postgres/redis/opensearch.enabled` | `true` | disable for managed services |

Validate locally without a cluster:

```bash
helm lint deploy/helm/kirmya
helm template cb deploy/helm/kirmya | less
```
