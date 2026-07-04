// Package api is the HTTP delivery layer for the career context.
package api

import (
	"net/http"
	"strings"

	"workspace-app/internal/career/application"
	"workspace-app/internal/common"
)

// Handler holds the career HTTP handlers.
type Handler struct{ svc *application.Service }

// NewHandler constructs the career handler.
func NewHandler(svc *application.Service) *Handler { return &Handler{svc: svc} }

// Paths handles GET /career/paths?from=<role> (authenticated): returns the role
// ladder, pay bands and skill gap reachable from the given starting role.
func (h *Handler) Paths(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	if strings.TrimSpace(from) == "" {
		common.WriteValidationError(w, "from (current role) is required")
		return
	}
	common.WriteSuccess(w, http.StatusOK, h.svc.Paths(r.Context(), from))
}
