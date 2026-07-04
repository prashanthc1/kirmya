# Recession Recovery Workspace

A modular-monolith Go platform for connecting job seekers, recruiters, founders, freelancers, mentors, and collaborators during recession and transition periods.

## Architecture

- `cmd/workspace-app` — entrypoint for the HTTP API server
- `internal/platform` — core server wiring, routing, middleware
- `internal/auth` — authentication and user management routes
- `internal/jobs` — job posting and application routes
- `internal/workspace` — workspace and collaboration endpoints
- `internal/chat` — realtime chat APIs
- `internal/meetings` — meeting and video call scaffolding
- `internal/notifications` — notification delivery endpoints
- `internal/matching` — AI-style matching and recommendations
- `internal/community` — community discussion endpoints
- `internal/ideas` — startup idea collaboration endpoints
- `web/swagger-ui` — Swagger UI HTML frontend
- `docs/openapi.yaml` — OpenAPI 3.0 API contract

## Run

1. Open a terminal at the repository root.
2. Run `go run ./cmd/workspace-app`.
3. Visit `http://localhost:8080/api/v1/health`.
4. View the API contract at `http://localhost:8080/openapi.yaml`.
5. Open the Swagger UI at `http://localhost:8080/swagger-ui/`.

## Notes

This repository is designed as a modular monolith with microservice-ready principles:

- clear package boundaries
- separate API modules for each domain
- lightweight HTTP handlers with route registration
- static Swagger/OpenAPI contract included

The implementation is a scaffold ready for database, authentication, messaging, and AI integration.
