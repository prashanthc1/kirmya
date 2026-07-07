package platform

import (
	"database/sql"
	"net/http"
	"time"

	"workspace-app/internal/admin"
	"workspace-app/internal/ai"
	"workspace-app/internal/career"
	"workspace-app/internal/common"
	"workspace-app/internal/community"
	"workspace-app/internal/dashboard"
	"workspace-app/internal/identity"
	identitydomain "workspace-app/internal/identity/domain"
	"workspace-app/internal/jobs"
	"workspace-app/internal/mentorship"
	"workspace-app/internal/messaging"
	"workspace-app/internal/network"
	"workspace-app/internal/notifications"
	platformcache "workspace-app/internal/platform/cache"
	"workspace-app/internal/platform/eventbus"
	platformmiddleware "workspace-app/internal/platform/middleware"
	"workspace-app/internal/platform/observability"
	"workspace-app/internal/platform/outbox"
	platformsearch "workspace-app/internal/platform/search"
	"workspace-app/internal/profile"
	"workspace-app/internal/referrals"
	"workspace-app/internal/resume"
	"workspace-app/internal/search"
	"workspace-app/internal/settings"
	settingsapp "workspace-app/internal/settings/application"
	settingspg "workspace-app/internal/settings/infrastructure/postgres"
)

// NewRouter builds the route map for the Kirmya platform.
func NewRouter(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		common.WriteSuccess(w, http.StatusOK, map[string]string{
			"status":  "healthy",
			"service": "kirmya",
		})
	})

	// Prometheus metrics endpoint + DB connection-pool gauges.
	mux.Handle("GET /metrics", observability.Handler())
	observability.RegisterDBStats(db)

	// Event bus (in-process; NATS-ready). Modules publish/subscribe here.
	bus := eventbus.New()

	// Outbox publisher handles transactional event writing to event_outbox table.
	pub := outbox.NewPublisher(db)
	outboxBus := outbox.NewBus(bus, pub)

	// Outbox relay polls the event_outbox table and processes/publishes them onto NATS/EventBus.
	relay := outbox.NewRelay(db, bus)
	relay.Start(250 * time.Millisecond)

	// Cache-aside layer (Redis when REDIS_URL is set; no-op otherwise).
	cache := platformcache.New()

	// Full-text search engine (OpenSearch when OPENSEARCH_URL is set; DB fallback
	// otherwise).
	searchEngine := platformsearch.New()

	// Stateless settings read-service shared by other modules to enforce a user's
	// privacy and notification preferences.
	settingsReader := settingsapp.NewService(settingspg.NewRepository(db), outboxBus)

	// Identity is the composition root for auth. It replaces the former auth +
	// user modules and provides the shared JWT auth middleware.
	identityMod := identity.NewModule(db, cache, outboxBus)
	identityMod.RegisterRoutes(mux)

	// Redis-backed token bucket rate limiter for connections and messaging
	redisLimiter := platformmiddleware.NewRedisRateLimiter(db)

	// Feature modules — all on Postgres/DDD, sharing identity's auth middleware.
	profile.RegisterRoutes(mux, db, identityMod.AuthMiddleware, outboxBus, cache, settingsReader)
	jobs.RegisterRoutes(mux, db, identityMod.AuthMiddleware, identityMod.RoleMiddleware(identitydomain.RoleRecruiter), outboxBus, cache)
	referrals.RegisterRoutes(mux, db, identityMod.AuthMiddleware, outboxBus)
	resume.RegisterRoutes(mux, db, identityMod.AuthMiddleware, outboxBus)
	ai.RegisterRoutes(mux, db, identityMod.AuthMiddleware)
	messaging.RegisterRoutes(mux, db, identityMod.AuthMiddleware, outboxBus, settingsReader, redisLimiter.Limit)
	mentorship.RegisterRoutes(mux, db, identityMod.AuthMiddleware, outboxBus)
	community.RegisterRoutes(mux, db, identityMod.AuthMiddleware, outboxBus)
	notifications.RegisterRoutes(mux, db, identityMod.AuthMiddleware, bus, settingsReader)
	settings.RegisterRoutes(mux, db, identityMod.AuthMiddleware, outboxBus)
	admin.RegisterRoutes(mux, db, identityMod.AuthMiddleware, identityMod.AdminMiddleware)
	search.RegisterRoutes(mux, db, identityMod.AuthMiddleware, bus, searchEngine)
	dashboard.RegisterRoutes(mux, db, identityMod.AuthMiddleware)
	career.RegisterRoutes(mux, identityMod.AuthMiddleware)
	network.RegisterRoutes(mux, db, identityMod.AuthMiddleware, bus, redisLimiter.Limit)

	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("web/swagger-ui"))))
	mux.Handle("/openapi.yaml", http.FileServer(http.Dir("docs")))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger-ui/", http.StatusTemporaryRedirect)
	})

	return mux
}
