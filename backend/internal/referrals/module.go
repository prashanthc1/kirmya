// Package referrals is the composition root for the Referrals bounded context.
package referrals

import (
	"database/sql"
	"net/http"

	"workspace-app/internal/referrals/api"
	"workspace-app/internal/referrals/application"
	"workspace-app/internal/referrals/infrastructure/postgres"
)

// RegisterRoutes wires the referrals module and mounts its routes behind the
// provided auth middleware. events may be nil.
func RegisterRoutes(mux *http.ServeMux, db *sql.DB, authMiddleware func(http.Handler) http.Handler, events application.EventPublisher) {
	repo := postgres.NewRepository(db)
	svc := application.NewService(repo, events)
	api.RegisterRoutes(mux, api.NewHandler(svc), authMiddleware)
}
