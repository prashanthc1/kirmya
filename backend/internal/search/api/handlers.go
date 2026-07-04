package api

import (
	"net/http"
	"strconv"
	"strings"

	"workspace-app/internal/common"
	"workspace-app/internal/search/application"
)

type Handler struct {
	svc *application.Service
}

func NewHandler(svc *application.Service) *Handler { return &Handler{svc: svc} }

// Search handles GET /search?q=&type=&limit=. type may be repeated or comma-
// separated (user, job, community, skill).
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	hits, err := h.svc.Query(r.Context(), q.Get("q"), parseTypes(q), parseLimit(q.Get("limit")))
	if err != nil {
		common.WriteInternalError(w, "search failed")
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]any{
		"results": hits,
		"engine":  engineLabel(h.svc.Ready()),
	})
}

// Autocomplete handles GET /search/autocomplete?q=&limit=.
func (h *Handler) Autocomplete(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	hits, err := h.svc.Suggest(r.Context(), q.Get("q"), parseLimit(q.Get("limit")))
	if err != nil {
		common.WriteInternalError(w, "search failed")
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]any{"results": hits})
}

func parseTypes(q map[string][]string) []string {
	var out []string
	for _, raw := range q["type"] {
		for _, t := range strings.Split(raw, ",") {
			if t = strings.TrimSpace(t); t != "" {
				out = append(out, t)
			}
		}
	}
	return out
}

func parseLimit(s string) int {
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	return 0
}

func engineLabel(ready bool) string {
	if ready {
		return "opensearch"
	}
	return "database"
}
