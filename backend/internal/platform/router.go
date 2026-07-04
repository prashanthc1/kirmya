package platform

import (
	"database/sql"
	"net/http"

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
	"workspace-app/internal/notifications"
	platformcache "workspace-app/internal/platform/cache"
	"workspace-app/internal/platform/eventbus"
	"workspace-app/internal/platform/observability"
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

	// Cache-aside layer (Redis when REDIS_URL is set; no-op otherwise).
	cache := platformcache.New()

	// Full-text search engine (OpenSearch when OPENSEARCH_URL is set; DB fallback
	// otherwise).
	searchEngine := platformsearch.New()

	// Stateless settings read-service shared by other modules to enforce a user's
	// privacy and notification preferences.
	settingsReader := settingsapp.NewService(settingspg.NewRepository(db), bus)

	// Identity is the composition root for auth. It replaces the former auth +
	// user modules and provides the shared JWT auth middleware.
	identityMod := identity.NewModule(db, bus)
	identityMod.RegisterRoutes(mux)

	// Feature modules — all on Postgres/DDD, sharing identity's auth middleware.
	profile.RegisterRoutes(mux, db, identityMod.AuthMiddleware, bus, cache, settingsReader)
	jobs.RegisterRoutes(mux, db, identityMod.AuthMiddleware, identityMod.RoleMiddleware(identitydomain.RoleRecruiter), bus, cache)
	referrals.RegisterRoutes(mux, db, identityMod.AuthMiddleware, bus)
	resume.RegisterRoutes(mux, db, identityMod.AuthMiddleware, bus)
	ai.RegisterRoutes(mux, db, identityMod.AuthMiddleware)
	messaging.RegisterRoutes(mux, db, identityMod.AuthMiddleware, bus, settingsReader)
	mentorship.RegisterRoutes(mux, db, identityMod.AuthMiddleware, bus)
	community.RegisterRoutes(mux, db, identityMod.AuthMiddleware, bus)
	notifications.RegisterRoutes(mux, db, identityMod.AuthMiddleware, bus, settingsReader)
	settings.RegisterRoutes(mux, db, identityMod.AuthMiddleware, bus)
	admin.RegisterRoutes(mux, db, identityMod.AuthMiddleware, identityMod.AdminMiddleware)
	search.RegisterRoutes(mux, db, identityMod.AuthMiddleware, bus, searchEngine)
	dashboard.RegisterRoutes(mux, db, identityMod.AuthMiddleware)
	career.RegisterRoutes(mux, identityMod.AuthMiddleware)

	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("web/swagger-ui"))))
	mux.Handle("/openapi.yaml", http.FileServer(http.Dir("docs")))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger-ui/", http.StatusTemporaryRedirect)
	})

	return mux
}
