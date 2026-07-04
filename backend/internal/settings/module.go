// Package settings is the composition root for the Settings context.
package settings

import (
	"database/sql"
	"net/http"

	"workspace-app/internal/settings/api"
	"workspace-app/internal/settings/application"
	"workspace-app/internal/settings/infrastructure/postgres"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB, authMiddleware func(http.Handler) http.Handler, events application.EventPublisher) {
	svc := application.NewService(postgres.NewRepository(db), events)
	api.RegisterRoutes(mux, api.NewHandler(svc), authMiddleware)
}
