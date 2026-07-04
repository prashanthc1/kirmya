// Package resume is the composition root for the Resume bounded context.
package resume

import (
	"database/sql"
	"net/http"

	"workspace-app/internal/resume/api"
	"workspace-app/internal/resume/application"
	"workspace-app/internal/resume/infrastructure/parser"
	"workspace-app/internal/resume/infrastructure/postgres"
	"workspace-app/internal/resume/infrastructure/storage"
)

// RegisterRoutes wires the resume module (Postgres repo + filesystem storage +
// parser + service + HTTP api). events may be nil.
func RegisterRoutes(mux *http.ServeMux, db *sql.DB, authMiddleware func(http.Handler) http.Handler, events application.EventPublisher) {
	svc := application.NewService(
		postgres.NewRepository(db),
		storage.NewFileSystem(),
		parser.New(),
		events,
	)
	api.RegisterRoutes(mux, api.NewHandler(svc), authMiddleware)
}
