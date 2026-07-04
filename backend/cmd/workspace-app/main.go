// Command workspace-app is the composition root for the Kirmya backend: it loads
// secrets, opens PostgreSQL, applies migrations, optionally seeds demo data,
// wires the platform HTTP server (which assembles every bounded context), and
// serves the REST/JSON API on $PORT (default 8080) with graceful shutdown.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"workspace-app/internal/platform"
	"workspace-app/internal/platform/observability"
	"workspace-app/internal/platform/secrets"
	"workspace-app/internal/platform/seed"
)

func main() {
	log.SetFlags(log.LstdFlags | log.LUTC)

	// Secrets first: this may export values (e.g. JWT_SECRET, DATABASE_URL) into
	// the environment before anything reads them. With the default backend it is
	// a no-op.
	if err := secrets.Load(log.Default()); err != nil {
		log.Fatalf("[startup] load secrets: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Tracing: OTLP exporter when OTEL_EXPORTER_OTLP_ENDPOINT is set, else no-op.
	shutdownTracing, err := observability.InitTracing(ctx)
	if err != nil {
		log.Fatalf("[startup] init tracing: %v", err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdownTracing(shutdownCtx); err != nil {
			log.Printf("[shutdown] flush traces: %v", err)
		}
	}()

	db, err := platform.OpenDatabase()
	if err != nil {
		log.Fatalf("[startup] open database: %v", err)
	}
	defer func() { _ = db.Close() }()

	if err := platform.RunMigrations(db); err != nil {
		log.Fatalf("[startup] run migrations: %v", err)
	}

	// Demo data — seed.Run is a no-op unless SEED_DEMO_DATA=true.
	if err := seed.Run(db); err != nil {
		log.Fatalf("[startup] seed demo data: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := platform.NewServer(port, db)

	// Serve until a signal arrives, then drain in-flight requests.
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("[startup] listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		log.Fatalf("[server] listen: %v", err)
	case <-ctx.Done():
		log.Printf("[shutdown] signal received; draining connections")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("[shutdown] graceful shutdown failed: %v", err)
		_ = srv.Close()
	}
	log.Printf("[shutdown] stopped")
}
