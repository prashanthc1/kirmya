package api

import (
	"net/http"
	"strings"

	"workspace-app/internal/common"
	"workspace-app/internal/identity/infrastructure/jwtauth"
)

type Routes struct {
	handler *Handler
	tokens  *jwtauth.Factory
}

func NewRoutes(h *Handler, tokens *jwtauth.Factory) *Routes {
	return &Routes{handler: h, tokens: tokens}
}

func (rt *Routes) Register(mux *http.ServeMux) {
	// Custom middleware for optional authentication
	optAuth := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			parts := strings.Fields(req.Header.Get("Authorization"))
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				claims, err := rt.tokens.Parse(parts[1])
				if err == nil {
					req = req.WithContext(common.ContextWithAuthUser(req.Context(), common.AuthUser{
						ID:    claims.Subject,
						Email: claims.Email,
					}))
				}
			}
			next.ServeHTTP(w, req)
		})
	}

	// GET, POST, PUT, DELETE endpoints mounted at both /api/cookies/preferences and /api/v1/cookies/preferences
	for _, prefix := range []string{"/api/cookies/preferences", "/api/v1/cookies/preferences"} {
		mux.Handle("GET "+prefix, optAuth(http.HandlerFunc(rt.handler.Get)))
		mux.Handle("POST "+prefix, optAuth(http.HandlerFunc(rt.handler.Save)))
		mux.Handle("PUT "+prefix, optAuth(http.HandlerFunc(rt.handler.Save)))
		mux.Handle("DELETE "+prefix, optAuth(http.HandlerFunc(rt.handler.Delete)))
	}
}
