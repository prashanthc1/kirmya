// Package mentorship is the composition root for the Mentorship context.
package mentorship

import (
	"database/sql"
	"net/http"

	"workspace-app/internal/mentorship/api"
	"workspace-app/internal/mentorship/application"
	"workspace-app/internal/mentorship/infrastructure/postgres"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB, authMiddleware func(http.Handler) http.Handler, events application.EventPublisher) {
	svc := application.NewService(postgres.NewRepository(db), events)
	api.RegisterRoutes(mux, api.NewHandler(svc), authMiddleware)
}
