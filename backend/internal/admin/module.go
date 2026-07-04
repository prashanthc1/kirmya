// Package admin is the composition root for the Admin bounded context. It wires
// the PostgreSQL repository to the application service and the HTTP api. Admin
// routes are gated by the RBAC admin role; the report-filing route is open to
// any authenticated user.
package admin

import (
	"database/sql"
	"net/http"

	"workspace-app/internal/admin/api"
	"workspace-app/internal/admin/application"
	"workspace-app/internal/admin/infrastructure/postgres"
)

// RegisterRoutes wires the admin module. auth authenticates any user (for filing
// reports); adminOnly additionally enforces the admin role (for /admin/*).
func RegisterRoutes(mux *http.ServeMux, db *sql.DB, auth, adminOnly func(http.Handler) http.Handler) {
	repo := postgres.NewRepository(db)
	svc := application.NewService(repo)
	api.RegisterRoutes(mux, api.NewHandler(svc), auth, adminOnly)
}
