// Package messaging is the composition root for the Messaging context.
package messaging

import (
	"database/sql"
	"net/http"

	"workspace-app/internal/messaging/api"
	"workspace-app/internal/messaging/application"
	"workspace-app/internal/messaging/infrastructure/postgres"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB, authMiddleware func(http.Handler) http.Handler, bus application.Bus, policy application.MessagePolicyReader) {
	svc := application.NewService(postgres.NewRepository(db), bus, application.NewHub(bus))
	svc.SetMessagePolicyReader(policy)
	api.RegisterRoutes(mux, api.NewHandler(svc), authMiddleware)
}
