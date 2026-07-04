package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"workspace-app/internal/common"
	"workspace-app/internal/notifications/application"
	"workspace-app/internal/notifications/domain"
)

type Handler struct{ svc *application.Service }

func NewHandler(svc *application.Service) *Handler { return &Handler{svc: svc} }

// Stream handles GET /notifications/stream — a Server-Sent Events stream that
// pushes notifications to the authenticated user in real time. Clients connect
// with `fetch` (so the Bearer token rides the Authorization header) and parse
// the text/event-stream body.
func (h *Handler) Stream(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // disable proxy (nginx) buffering

	rc := http.NewResponseController(w)
	// SSE connections are long-lived; clear the server's WriteTimeout for this one.
	_ = rc.SetWriteDeadline(time.Time{})

	ch, cancel := h.svc.Subscribe(uid)
	defer cancel()

	_, _ = fmt.Fprint(w, ": connected\n\n")
	if rc.Flush() != nil {
		return
	}

	ping := time.NewTicker(25 * time.Second)
	defer ping.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case n, ok := <-ch:
			if !ok {
				return
			}
			data, _ := json.Marshal(toDTO(&n))
			_, _ = fmt.Fprintf(w, "event: notification\ndata: %s\n\n", data)
			if rc.Flush() != nil {
				return
			}
		case <-ping.C:
			_, _ = fmt.Fprint(w, ": ping\n\n")
			if rc.Flush() != nil {
				return
			}
		}
	}
}

type notificationDTO struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Body      string `json:"body,omitempty"`
	Link      string `json:"link,omitempty"`
	Read      bool   `json:"read"`
	CreatedAt string `json:"created_at"`
}

func toDTO(n *domain.Notification) notificationDTO {
	return notificationDTO{
		ID: n.ID, Type: n.Type, Title: n.Title, Body: n.Body, Link: n.Link,
		Read: n.ReadAt != nil, CreatedAt: n.CreatedAt.UTC().Format(time.RFC3339),
	}
}

// List handles GET /notifications.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	limit, offset := pageParams(r)
	items, err := h.svc.List(r.Context(), uid, limit, offset)
	if err != nil {
		common.WriteInternalError(w, "could not load notifications")
		return
	}
	unread, _ := h.svc.UnreadCount(r.Context(), uid)
	out := make([]notificationDTO, 0, len(items))
	for i := range items {
		out = append(out, toDTO(&items[i]))
	}
	common.WriteSuccess(w, http.StatusOK, map[string]any{
		"notifications": out, "unread": unread, "limit": limit, "offset": offset,
	})
}

// pageParams parses ?limit=&offset= with a default page size of 50 (capped 100).
func pageParams(r *http.Request) (limit, offset int) {
	limit, offset = 50, 0
	if v, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && v > 0 {
		limit = v
	}
	if limit > 100 {
		limit = 100
	}
	if v, err := strconv.Atoi(r.URL.Query().Get("offset")); err == nil && v > 0 {
		offset = v
	}
	return limit, offset
}

// MarkRead handles POST /notifications/{id}/read.
func (h *Handler) MarkRead(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.MarkRead(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("id")); err != nil {
		common.WriteInternalError(w, "could not update notification")
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]bool{"read": true})
}

// MarkAllRead handles POST /notifications/read-all.
func (h *Handler) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.MarkAllRead(r.Context(), common.UserIDFromContext(r.Context())); err != nil {
		common.WriteInternalError(w, "could not update notifications")
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]bool{"read": true})
}
