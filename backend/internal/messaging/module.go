// Package messaging is the composition root for the Messaging context.
package messaging

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"workspace-app/internal/messaging/api"
	"workspace-app/internal/messaging/application"
	"workspace-app/internal/messaging/infrastructure/postgres"
)

type dbConnectionChecker struct {
	db *sql.DB
}

func (c *dbConnectionChecker) AreConnected(ctx context.Context, userA, userB string) (bool, error) {
	var status string
	err := c.db.QueryRowContext(ctx, `
		SELECT status FROM user_connections
		WHERE (requester_id = $1 AND receiver_id = $2) OR (requester_id = $2 AND receiver_id = $1)
	`, userA, userB).Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return status == "accepted", nil
}

func RegisterRoutes(mux *http.ServeMux, db *sql.DB, authMiddleware func(http.Handler) http.Handler, bus application.Bus, policy application.MessagePolicyReader, limit func(string) func(http.Handler) http.Handler) {
	hub := application.NewHub(bus)
	repo := postgres.NewRepository(db)
	svc := application.NewService(repo, bus, hub, &dbConnectionChecker{db: db})
	svc.SetMessagePolicyReader(policy)
	api.RegisterRoutes(mux, api.NewHandler(svc), authMiddleware, limit)
}
