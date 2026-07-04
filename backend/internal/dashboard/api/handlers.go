// Package api is the HTTP delivery layer for the dashboard context.
package api

import (
	"net/http"

	"workspace-app/internal/common"
	"workspace-app/internal/dashboard/application"
)

// Handler holds the dashboard HTTP handlers.
type Handler struct{ svc *application.Service }

// NewHandler constructs the dashboard handler.
func NewHandler(svc *application.Service) *Handler { return &Handler{svc: svc} }

// Summary handles GET /me/dashboard (authenticated): returns the caller's
// aggregated, role-segmented dashboard counts.
func (h *Handler) Summary(w http.ResponseWriter, r *http.Request) {
	s, err := h.svc.Summary(r.Context(), common.UserIDFromContext(r.Context()))
	if err != nil {
		common.WriteInternalError(w, "could not load dashboard")
		return
	}
	common.WriteSuccess(w, http.StatusOK, s)
}
