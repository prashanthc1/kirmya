// Package profile is the composition root for the Profile bounded context.
package profile

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"workspace-app/internal/profile/api"
	"workspace-app/internal/profile/application"
	"workspace-app/internal/profile/infrastructure/postgres"
)

// RegisterRoutes wires the profile module (Postgres repo + service + HTTP api)
// and mounts its routes behind the provided auth middleware. events and cache
// may be nil.
func RegisterRoutes(mux *http.ServeMux, db *sql.DB, authMiddleware func(http.Handler) http.Handler, events application.EventPublisher, cache application.Cache, visibility api.VisibilityReader) {
	repo := postgres.NewRepository(db)
	svc := application.NewService(repo, events, cache)
	h := api.NewHandler(svc)
	h.SetVisibilityReader(visibility)
	api.RegisterRoutes(mux, h, authMiddleware)
}

// RegisterGinRoutes wires the profile module and mounts its routes on the Gin engine
func RegisterGinRoutes(r *gin.Engine, db *sql.DB, authMiddleware func(http.Handler) http.Handler, events application.EventPublisher, cache application.Cache) {
	repo := postgres.NewRepository(db)
	svc := application.NewService(repo, events, cache)
	api.RegisterGinRoutes(r, db, authMiddleware, svc, cache)
}
