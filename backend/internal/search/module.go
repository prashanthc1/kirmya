// Package search is the composition root for the Search bounded context. It
// wires the OpenSearch engine (or DB fallback) to the application service and
// HTTP api, subscribes to domain events to keep the index fresh, and backfills
// the index on startup.
package search

import (
	"context"
	"net/http"
	"time"

	"database/sql"

	"workspace-app/internal/platform/eventbus"
	"workspace-app/internal/platform/search"
	"workspace-app/internal/search/api"
	"workspace-app/internal/search/application"
	"workspace-app/internal/search/infrastructure/postgres"
)

// RegisterRoutes wires the search module, subscribes to events, and kicks off a
// one-time index backfill. engine and bus may be nil-equivalent (a Noop engine
// disables indexing and uses the DB fallback for queries).
func RegisterRoutes(mux *http.ServeMux, db *sql.DB, authMiddleware func(http.Handler) http.Handler, bus *eventbus.Bus, engine search.Engine) {
	svc := application.NewService(engine, postgres.NewSource(db))
	api.RegisterRoutes(mux, api.NewHandler(svc), authMiddleware)

	if bus != nil {
		subscribe(bus, svc)
	}

	// Backfill in the background so startup isn't blocked by indexing.
	if engine.Ready() {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()
			svc.Backfill(ctx)
		}()
	}
}

// subscribe keeps the index fresh as entities change.
func subscribe(bus *eventbus.Bus, svc *application.Service) {
	indexUser := func(ctx context.Context, e eventbus.Event) { svc.IndexUser(ctx, e.AggregateID) }
	bus.Subscribe("UserRegistered", indexUser)
	bus.Subscribe("ProfileUpdated", indexUser)

	bus.Subscribe("JobPosted", func(ctx context.Context, e eventbus.Event) {
		svc.IndexJob(ctx, e.AggregateID)
	})
}
