// Package ai is the composition root for the AI bounded context (resume
// reviewer, career coach, skill-gap engine) backed by Claude.
package ai

import (
	"database/sql"
	"net/http"

	"workspace-app/internal/ai/api"
	"workspace-app/internal/ai/application"
	"workspace-app/internal/ai/infrastructure/anthropic"
	"workspace-app/internal/ai/infrastructure/postgres"
)

// RegisterRoutes wires the AI module. The LLM provider (Claude) reads
// ANTHROPIC_API_KEY; if unset, AI endpoints return 503 rather than failing hard.
func RegisterRoutes(mux *http.ServeMux, db *sql.DB, authMiddleware func(http.Handler) http.Handler) {
	svc := application.NewService(postgres.NewRepository(db), anthropic.New())
	api.RegisterRoutes(mux, api.NewHandler(svc), authMiddleware)
}
