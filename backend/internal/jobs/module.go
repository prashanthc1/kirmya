// Package jobs is the composition root for the Jobs bounded context.
package jobs

import (
	"database/sql"
	"net/http"

	"workspace-app/internal/jobs/api"
	"workspace-app/internal/jobs/application"
	"workspace-app/internal/jobs/infrastructure/aimatch"
	"workspace-app/internal/jobs/infrastructure/postgres"
)

// RegisterRoutes wires the jobs module (Postgres repo + service + HTTP api).
// authMiddleware gates seeker actions; recruiterOnly gates posting/managing jobs
// and applicants. events and cache may be nil. Job matching uses the AI matcher
// (Claude) with a heuristic fallback, so it works with or without an API key.
func RegisterRoutes(mux *http.ServeMux, db *sql.DB, authMiddleware, recruiterOnly func(http.Handler) http.Handler, events application.EventPublisher, cache application.Cache) {
	repo := postgres.NewRepository(db)
	svc := application.NewService(repo, events, cache)

	matcher := aimatch.New(application.NewHeuristicMatcher())
	matchSvc := application.NewMatchService(repo, postgres.NewSkillReader(db), matcher)

	api.RegisterRoutes(mux, api.NewHandler(svc, matchSvc), authMiddleware, recruiterOnly)
}
