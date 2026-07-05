package network

import (
	"database/sql"
	"net/http"

	"workspace-app/internal/network/api"
	"workspace-app/internal/network/application"
	"workspace-app/internal/network/infrastructure/postgres"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB, authMiddleware func(http.Handler) http.Handler) {
	repo := postgres.NewRepository(db)
	svc := application.NewService(repo)
	h := api.NewHandler(svc)
	api.RegisterRoutes(mux, h, authMiddleware)
}
