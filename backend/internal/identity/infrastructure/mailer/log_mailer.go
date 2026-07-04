// Package mailer provides a development Mailer that logs the verification or
// reset link to the console instead of delivering an email, plus a real SMTP
// adapter (smtp_mailer.go) for production. In development the LogMailer prints
// the full clickable link so the flow is testable without an SMTP server; in
// production it fails closed (refuses to send) rather than leaking an
// account-takeover token into production logs.
package mailer

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"net/url"
	"os"
	"strings"
)

// ErrNoProductionMailer is returned when the log-only mailer is asked to send in
// production, where it fails closed rather than silently dropping (or leaking)
// account-takeover tokens.
var ErrNoProductionMailer = errors.New("no production mailer configured")

type LogMailer struct {
	appURL   string
	prod     bool
	devLinks bool // log the full (token-bearing) link only in explicit development
}

func NewLogMailer() *LogMailer {
	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		appURL = "http://localhost:3000"
	}
	env := os.Getenv("APP_ENV")
	return &LogMailer{
		appURL: strings.TrimRight(appURL, "/"),
		prod:   env == "production",
		// Allowlist (fail-safe): only the explicit dev/empty env prints raw links;
		// any other value (staging, test, etc.) logs only the non-reversible ref.
		devLinks: env == "" || env == "development" || env == "dev" || env == "local",
	}
}

func (m *LogMailer) SendVerificationEmail(_ context.Context, email, rawToken string) error {
	if m.prod {
		log.Printf("[mailer] REFUSING to send verification email in production: no real mailer configured")
		return ErrNoProductionMailer
	}
	if m.devLinks {
		// Dev only: print the full link so the verification flow can be completed
		// locally without an SMTP server. Dev tokens are local-only.
		link := m.appURL + "/verify-email?token=" + url.QueryEscape(rawToken)
		log.Printf("[mailer] DEV verification link for %s (ref %s): %s", email, tokenRef(rawToken), link)
		return nil
	}
	// Non-production, non-dev (e.g. staging): never log the raw token.
	log.Printf("[mailer] verification email queued for %s (ref %s)", email, tokenRef(rawToken))
	return nil
}

func (m *LogMailer) SendPasswordResetEmail(_ context.Context, email, rawToken string) error {
	if m.prod {
		log.Printf("[mailer] REFUSING to send password reset email in production: no real mailer configured")
		return ErrNoProductionMailer
	}
	if m.devLinks {
		// Dev only: print the full link so the reset flow can be completed
		// locally without an SMTP server. Dev tokens are local-only.
		link := m.appURL + "/reset-password?token=" + url.QueryEscape(rawToken)
		log.Printf("[mailer] DEV password reset link for %s (ref %s): %s", email, tokenRef(rawToken), link)
		return nil
	}
	// Non-production, non-dev (e.g. staging): never log the raw token.
	log.Printf("[mailer] password reset email queued for %s (ref %s)", email, tokenRef(rawToken))
	return nil
}

// tokenRef returns a short, non-reversible reference to rawToken (the first 4
// bytes of its SHA-256 hash, hex-encoded) suitable for correlating log lines
// without leaking the token itself.
func tokenRef(rawToken string) string {
	sum := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(sum[:4])
}
