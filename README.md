# Recession Recovery Workspace

A platform for connecting job seekers, recruiters, founders, freelancers, mentors, and collaborators during recession and transition periods.

This is a monorepo split into two projects:

- [`backend/`](backend/) — Go modular-monolith HTTP API (see [backend/README.md](backend/README.md))
- [`frontend/`](frontend/) — Next.js web application (see [frontend/README.md](frontend/README.md))

## Getting started

### Backend

```
cd backend
go run ./cmd/workspace-app
```

API runs at `http://localhost:8080`. See [backend/README.md](backend/README.md) for details.

### Frontend

```
cd frontend
npm install
npm run dev
```

App runs at `http://localhost:3000`.
