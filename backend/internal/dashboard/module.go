// Package dashboard is the composition root for the Dashboard bounded context.
// It wires the PostgreSQL read adapter to the application service and the HTTP
// api. The dashboard is a read-only backend-for-frontend projection.
package dashboard

import (
	"database/sql"
	"net/http"

	"workspace-app/internal/dashboard/api"
	"workspace-app/internal/dashboard/application"
	"workspace-app/internal/dashboard/infrastructure/postgres"
)

// RegisterRoutes wires the dashboard module. auth authenticates the caller.
func RegisterRoutes(mux *http.ServeMux, db *sql.DB, auth func(http.Handler) http.Handler) {
	repo := postgres.NewRepository(db)
	svc := application.NewService(repo)
	api.RegisterRoutes(mux, api.NewHandler(svc), auth)
}
