// Package messaging is the composition root for the Messaging context.
package messaging

import (
	"context"
	"database/sql"
	"net/http"

	"workspace-app/internal/connections"
	"workspace-app/internal/messaging/api"
	"workspace-app/internal/messaging/application"
	"workspace-app/internal/messaging/infrastructure/postgres"
)

type dbConnectionChecker struct {
	db *sql.DB
}

func (c *dbConnectionChecker) AreConnected(ctx context.Context, userA, userB string) (bool, error) {
	return connections.CanMessage(ctx, c.db, userA, userB)
}

func RegisterRoutes(mux *http.ServeMux, db *sql.DB, authMiddleware func(http.Handler) http.Handler, bus application.Bus, policy application.MessagePolicyReader, limit func(string) func(http.Handler) http.Handler) {
	hub := application.NewHub(bus)
	repo := postgres.NewRepository(db)
	svc := application.NewService(repo, bus, hub, &dbConnectionChecker{db: db})
	svc.SetMessagePolicyReader(policy)
	api.RegisterRoutes(mux, api.NewHandler(svc), authMiddleware, limit)
}
