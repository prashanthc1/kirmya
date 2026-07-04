# Kirmya — Backend Folder Structure

> Go 1.26 · module `workspace-app` · modular monolith · Clean Architecture / DDD per module.

## 1. Top-Level Layout

```
backend/
├── cmd/
│   └── workspace-app/
│       └── main.go               # composition root: load config, open DB, wire modules, start server
├── internal/
│   ├── platform/                 # framework & wiring (not a domain module)
│   │   ├── config.go             # env/config loading
│   │   ├── database.go           # PostgreSQL pool (pgx) + health
│   │   ├── migrations.go         # migration runner
│   │   ├── router.go             # global router + module registration
│   │   ├── server.go             # http.Server lifecycle, graceful shutdown
│   │   ├── middleware/           # request-id, recover, security headers, cors, ratelimit, csrf, auth, rbac, otel, audit
│   │   ├── eventbus/             # in-process bus (NATS-ready) + outbox relay
│   │   ├── cache/                # Redis client + helpers
│   │   ├── search/               # OpenSearch client + indexers
│   │   └── observability/        # OTel/Prometheus setup
│   ├── common/                   # shared kernel: errors, response, context, pagination, validation, ids
│   └── <module>/                 # one folder per bounded context (see §2)
├── migrations/                   # PostgreSQL .up.sql/.down.sql + seed/
├── docs/                         # openapi.yaml (served via /swagger-ui)
├── web/swagger-ui/               # static Swagger UI
├── go.mod
└── go.sum
```

## 2. Per-Module Layout (DDD / Clean Architecture)

Target layout for every bounded context. Example: `identity`.

```
internal/identity/
├── module.go                     # public wiring: NewModule(deps) -> registers routes, subscribes events
├── domain/                       # pure domain — NO infra imports
│   ├── user.go                   # entities / aggregates + invariants
│   ├── role.go
│   ├── token.go                  # refresh-token rotation rules (value objects)
│   ├── events.go                 # UserRegistered, EmailVerified, ...
│   └── ports.go                  # repository + gateway interfaces (UserRepository, TokenStore, Mailer, ...)
├── application/                  # use cases (CQRS-ready): commands + queries
│   ├── register_user.go
│   ├── login.go
│   ├── refresh_token.go
│   ├── verify_email.go
│   ├── reset_password.go
│   ├── oauth_login.go
│   └── service.go                # Service interface other modules depend on
├── infrastructure/               # adapters implementing domain ports
│   ├── postgres/
│   │   ├── user_repository.go
│   │   ├── token_repository.go
│   │   └── oauth_repository.go
│   ├── oauth/                    # google.go, linkedin.go
│   ├── crypto/                   # argon2 hashing, totp
│   └── mailer/                   # email sender adapter
└── api/                          # delivery layer (HTTP)
    ├── routes.go                 # registers handlers on the mux
    ├── handlers.go               # http.HandlerFunc per use case
    ├── dto.go                    # request/response structs + validation
    └── middleware.go             # module-specific (e.g. auth-only here)
```

### Dependency rule
`api → application → domain`. `infrastructure` depends on `domain` (implements its ports). `domain` imports nothing from other layers/modules. Modules depend on each other only via another module's `application.Service` interface, injected in `main.go`.

## 3. Existing → Target Mapping (incremental migration)

The current modules use `domain/dto/handler/repository/routes/service`. Map:

| Existing | Target |
|---|---|
| `service/` | `application/` |
| `repository/` + external clients | `infrastructure/` |
| `handler/` + `routes/` | `api/` |
| `dto/` | `api/dto.go` |
| `domain/` | `domain/` (add `ports.go`, `events.go`) |

Modules are migrated as we touch them; **identity is migrated/built first** to the target layout. Untouched modules keep working in the old layout (both compile under the same binary).

## 4. Composition Root (`main.go`) sketch

```go
cfg := platform.LoadConfig()
db := platform.OpenDatabase(cfg)          // pgx pool
platform.RunMigrations(db)
bus := eventbus.New()                      // in-process; NATS later
cache := cache.NewRedis(cfg)
search := search.NewOpenSearch(cfg)

identityMod := identity.NewModule(identity.Deps{DB: db, Bus: bus, Cache: cache, Cfg: cfg})
profileMod  := profile.NewModule(profile.Deps{DB: db, Bus: bus, Identity: identityMod.Service})
// ... other modules, injecting interfaces only

router := platform.NewRouter(platform.Modules{Identity: identityMod, Profile: profileMod /* ... */})
server := platform.NewServer(cfg.Port, router)
server.Start()                             // + graceful shutdown
```

## 5. Configuration (12-factor)

All config via env (`internal/platform/config.go`), with `.env` for local dev. Keys: `PORT`, `DATABASE_URL`, `REDIS_URL`, `OPENSEARCH_URL`, `JWT_SIGNING_KEY`, `JWT_ACCESS_TTL`, `JWT_REFRESH_TTL`, `GOOGLE_CLIENT_ID/SECRET`, `LINKEDIN_CLIENT_ID/SECRET`, `CLAUDE_API_KEY`, `OPENAI_API_KEY`, `MAIL_*`, `OTEL_EXPORTER_*`.

## 6. Testing Layout
- Unit tests beside code (`*_test.go`) — domain + application (table-driven, no DB).
- Repository/integration tests against a throwaway Postgres (testcontainers or dockerized) under `infrastructure/...`.
- API tests via `httptest` exercising the router.
- See [docs/09_... testing strategy in MVP roadmap]. Run with `go test ./...`.
