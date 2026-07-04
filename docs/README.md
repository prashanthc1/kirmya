# Kirmya — Documentation

Kirmya is a professional networking + **career-recovery** platform solving:
> "I lost my job. Help me get my next opportunity faster."

It is built by extending and rebranding the existing `workspace-app` modular monolith, re-platformed onto PostgreSQL.

## Index

| # | Doc | Contents |
|---|---|---|
| 01 | [Product Requirements (PRD)](01_PRD.md) | Problem, personas, journeys, MVP scope, metrics, requirements |
| 02 | [System Architecture](02_SYSTEM_ARCHITECTURE.md) | Modular monolith + DDD, high-level diagram, events, CQRS, extraction path |
| 03 | [Database Design](03_DATABASE_DESIGN.md) | PostgreSQL conventions, ERD, identity schema, migrations, seed |
| 04 | [Backend Structure](04_BACKEND_STRUCTURE.md) | Go folder layout, per-module DDD layers, composition root |
| 05 | [Frontend Structure](05_FRONTEND_STRUCTURE.md) | Next.js App Router layout, feature modules, API client |
| 06 | [API Contracts](06_API_CONTRACTS.md) | REST endpoints per module, examples, errors, rate limits |
| 10 | [MVP Roadmap & Production Strategy](10_MVP_ROADMAP.md) | Phases, testing, DevOps, deployment, scalability |
| 11 | [Running with Docker](11_RUNNING_WITH_DOCKER.md) | One-command local stack (Postgres + API + frontend) |

Machine-readable API spec: `backend/docs/openapi.yaml` (served at `/swagger-ui/`).

## Build order
Foundational docs (this set) → **Identity module** → Profile → Resume/Jobs/Referrals → Career/Communities/Mentorship → Messaging/Notifications/Admin → hardening. See the roadmap for details.
