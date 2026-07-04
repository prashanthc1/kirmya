package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"workspace-app/internal/admin/application"
	"workspace-app/internal/admin/domain"
	"workspace-app/internal/common"
)

type Handler struct {
	svc *application.Service
}

func NewHandler(svc *application.Service) *Handler { return &Handler{svc: svc} }

func decode(r *http.Request, dst any) bool {
	return json.NewDecoder(http.MaxBytesReader(nil, r.Body, 1<<20)).Decode(dst) == nil
}

func (h *Handler) writeErr(w http.ResponseWriter, err error) {
	var ve application.ValidationError
	switch {
	case errors.As(err, &ve):
		common.WriteValidationError(w, ve.Msg)
	case errors.Is(err, domain.ErrNotFound):
		common.WriteNotFoundError(w, "not found")
	default:
		common.WriteInternalError(w, "something went wrong")
	}
}

func adminID(r *http.Request) string { return common.UserIDFromContext(r.Context()) }

// ----- Users -------------------------------------------------------------

// ListUsers handles GET /admin/users.
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	f := domain.UserFilter{
		Query:  q.Get("q"),
		Status: q.Get("status"),
		Role:   q.Get("role"),
		Limit:  atoiDefault(q.Get("limit"), 25),
		Offset: atoiDefault(q.Get("offset"), 0),
	}
	users, total, err := h.svc.ListUsers(r.Context(), f)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]any{
		"users":  toUsers(users),
		"total":  total,
		"limit":  f.Limit,
		"offset": f.Offset,
	})
}

// GetUser handles GET /admin/users/{id}.
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	u, err := h.svc.GetUser(r.Context(), r.PathValue("id"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, toUser(u))
}

// SetUserStatus handles PATCH /admin/users/{id}/status.
func (h *Handler) SetUserStatus(w http.ResponseWriter, r *http.Request) {
	var req statusRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	u, err := h.svc.SetUserStatus(r.Context(), adminID(r), r.PathValue("id"), req.Status)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, toUser(u))
}

// AssignRole handles POST /admin/users/{id}/roles.
func (h *Handler) AssignRole(w http.ResponseWriter, r *http.Request) {
	var req roleRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	u, err := h.svc.AssignRole(r.Context(), adminID(r), r.PathValue("id"), req.Role)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, toUser(u))
}

// RevokeRole handles DELETE /admin/users/{id}/roles/{role}.
func (h *Handler) RevokeRole(w http.ResponseWriter, r *http.Request) {
	u, err := h.svc.RevokeRole(r.Context(), adminID(r), r.PathValue("id"), r.PathValue("role"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, toUser(u))
}

// ----- Moderation --------------------------------------------------------

// RemovePost handles DELETE /admin/posts/{id}.
func (h *Handler) RemovePost(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.RemovePost(r.Context(), adminID(r), r.PathValue("id")); err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]string{"status": "removed"})
}

// RemoveComment handles DELETE /admin/comments/{id}.
func (h *Handler) RemoveComment(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.RemoveComment(r.Context(), adminID(r), r.PathValue("id")); err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]string{"status": "removed"})
}

// ----- Reports -----------------------------------------------------------

// FileReport handles POST /reports (any authenticated user).
func (h *Handler) FileReport(w http.ResponseWriter, r *http.Request) {
	var req fileReportRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	rep, err := h.svc.FileReport(r.Context(), common.UserIDFromContext(r.Context()), application.ReportInput(req))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusCreated, toReport(rep))
}

// ListReports handles GET /admin/reports.
func (h *Handler) ListReports(w http.ResponseWriter, r *http.Request) {
	reports, err := h.svc.ListReports(r.Context(), r.URL.Query().Get("status"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]any{"reports": toReports(reports)})
}

// ResolveReport handles PATCH /admin/reports/{id}.
func (h *Handler) ResolveReport(w http.ResponseWriter, r *http.Request) {
	var req resolveReportRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	rep, err := h.svc.ResolveReport(r.Context(), adminID(r), r.PathValue("id"), req.Status, req.ActionTaken)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, toReport(rep))
}

// ----- Analytics ---------------------------------------------------------

// Stats handles GET /admin/stats.
func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	a, err := h.svc.Analytics(r.Context())
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, a)
}

func atoiDefault(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}
