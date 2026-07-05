// Package identity is the composition root for the Identity bounded context. It
// wires the domain ports to their PostgreSQL/crypto/JWT/OAuth adapters and the
// application service to the HTTP api, exposing a Service for other modules and
// auth Middleware for the platform router.
package identity

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"workspace-app/internal/identity/api"
	"workspace-app/internal/identity/application"
	"workspace-app/internal/identity/domain"
	"workspace-app/internal/identity/infrastructure/crypto"
	"workspace-app/internal/identity/infrastructure/jwtauth"
	"workspace-app/internal/identity/infrastructure/mailer"
	"workspace-app/internal/identity/infrastructure/oauth"
	"workspace-app/internal/identity/infrastructure/postgres"
	"workspace-app/internal/platform/tx"
)

// Module is the assembled identity context.
type Module struct {
	Service    *application.Service
	Middleware *api.Middleware
	handler    *api.Handler
}

// NewModule builds the identity module. events must satisfy
// domain.EventPublisher (the platform event bus does).
func NewModule(db *sql.DB, events domain.EventPublisher) *Module {
	tokens := jwtauth.NewFactory()

	deps := application.Deps{
		Users:     postgres.NewUserRepository(db),
		Refresh:   postgres.NewRefreshTokenRepository(db),
		Verif:     postgres.NewVerificationRepository(db),
		OAuth:     postgres.NewOAuthRepository(db),
		MFA:       postgres.NewMFARepository(db),
		Audit:     postgres.NewAuditRepository(db),
		Hasher:    crypto.NewArgon2Hasher(),
		Tokens:    tokens,
		TOTP:      crypto.NewTOTPService("Kirmya"),
		Mailer:    buildMailer(),
		Events:    events,
		Providers: buildProviders(),
		Tx:        tx.NewTxManager(db),
	}

	svc := application.NewService(deps)
	return &Module{
		Service:    svc,
		Middleware: api.NewMiddleware(tokens),
		handler:    api.NewHandler(svc),
	}
}

// RegisterRoutes mounts identity HTTP routes.
func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	api.RegisterRoutes(mux, m.handler, m.Middleware)
}

// AuthMiddleware returns the JWT authentication middleware for other modules.
func (m *Module) AuthMiddleware(next http.Handler) http.Handler {
	return m.Middleware.Authenticate(next)
}

// AdminMiddleware returns middleware that authenticates and requires the admin
// role. Used to gate the admin module's routes.
func (m *Module) AdminMiddleware(next http.Handler) http.Handler {
	return m.Middleware.RequireRole(domain.RoleAdmin)(next)
}

// RoleMiddleware returns middleware that authenticates and requires the given
// role, for modules that gate role-specific actions (e.g. recruiter-only).
func (m *Module) RoleMiddleware(role string) func(http.Handler) http.Handler {
	return m.Middleware.RequireRole(role)
}

// buildMailer selects the transactional mailer, in priority order: the Resend
// HTTP API (RESEND_API_KEY) — preferred because it sends over HTTPS and works
// where outbound SMTP is blocked (e.g. Railway below the Pro plan) — then a real
// SMTP adapter (SMTP_HOST), otherwise the development log-only mailer (which
// fails closed in production).
func buildMailer() domain.Mailer {
	// Prefer the Resend HTTP API: it delivers over HTTPS (443), so it works on
	// hosts that block outbound SMTP (e.g. Railway below the Pro plan).
	if m, err := mailer.NewResendMailer(); err == nil {
		log.Printf("[identity] Resend mailer enabled (HTTP API)")
		return m
	}
	if m, err := mailer.NewSMTPMailer(); err == nil {
		log.Printf("[identity] SMTP mailer enabled (host=%s)", os.Getenv("SMTP_HOST"))
		return m
	}
	if os.Getenv("APP_ENV") == "production" {
		log.Printf("[identity] WARNING: no RESEND_API_KEY or SMTP_HOST set in production; verification and password-reset emails will be refused, not sent")
	}
	return mailer.NewLogMailer()
}

// buildProviders constructs OAuth providers whose credentials are configured.
func buildProviders() map[string]application.OAuthProvider {
	providers := map[string]application.OAuthProvider{}
	if id := os.Getenv("GOOGLE_CLIENT_ID"); id != "" {
		providers[domain.ProviderGoogle] = oauth.NewGoogle(
			id, os.Getenv("GOOGLE_CLIENT_SECRET"), os.Getenv("GOOGLE_REDIRECT_URI"))
	}
	if id := os.Getenv("LINKEDIN_CLIENT_ID"); id != "" {
		providers[domain.ProviderLinkedIn] = oauth.NewLinkedIn(
			id, os.Getenv("LINKEDIN_CLIENT_SECRET"), os.Getenv("LINKEDIN_REDIRECT_URI"))
	}
	return providers
}
