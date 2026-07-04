package api

import (
	"net/http"
	"strings"

	"workspace-app/internal/common"
	"workspace-app/internal/identity/infrastructure/jwtauth"
)

// Middleware provides JWT authentication and RBAC for the platform router.
type Middleware struct {
	tokens *jwtauth.Factory
}

func NewMiddleware(tokens *jwtauth.Factory) *Middleware { return &Middleware{tokens: tokens} }

// Authenticate validates the Bearer access token and injects the auth user.
func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := m.parse(r)
		if !ok {
			common.WriteUnauthorizedError(w, "invalid or missing access token")
			return
		}
		r = r.WithContext(common.ContextWithAuthUser(r.Context(), common.AuthUser{
			ID:    claims.Subject,
			Email: claims.Email,
			Role:  primaryRole(claims.Roles),
		}))
		next.ServeHTTP(w, r)
	})
}

// RequireRole enforces that the caller holds the given role.
func (m *Middleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := m.parse(r)
			if !ok {
				common.WriteUnauthorizedError(w, "invalid or missing access token")
				return
			}
			if !hasRole(claims.Roles, role) {
				common.WriteForbiddenError(w, "insufficient permissions")
				return
			}
			r = r.WithContext(common.ContextWithAuthUser(r.Context(), common.AuthUser{
				ID:    claims.Subject,
				Email: claims.Email,
				Role:  primaryRole(claims.Roles),
			}))
			next.ServeHTTP(w, r)
		})
	}
}

func (m *Middleware) parse(r *http.Request) (*jwtauth.Claims, bool) {
	parts := strings.Fields(r.Header.Get("Authorization"))
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return nil, false
	}
	claims, err := m.tokens.Parse(parts[1])
	if err != nil {
		return nil, false
	}
	return claims, true
}

func primaryRole(roles []string) string {
	if len(roles) == 0 {
		return ""
	}
	return roles[0]
}

func hasRole(roles []string, want string) bool {
	for _, r := range roles {
		if r == want {
			return true
		}
	}
	return false
}
