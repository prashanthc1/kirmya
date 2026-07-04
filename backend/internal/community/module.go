// Package community is the composition root for the Communities context.
package community

import (
	"database/sql"
	"net/http"

	"workspace-app/internal/community/api"
	"workspace-app/internal/community/application"
	"workspace-app/internal/community/infrastructure/postgres"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB, authMiddleware func(http.Handler) http.Handler, events application.EventPublisher) {
	svc := application.NewService(postgres.NewRepository(db), events)
	api.RegisterRoutes(mux, api.NewHandler(svc), authMiddleware)
}
