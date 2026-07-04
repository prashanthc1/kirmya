// Package career is the composition root for the Career bounded context. Career
// ladders are curated reference data, so the module needs no database — it wires
// the application service to the HTTP api.
package career

import (
	"net/http"

	"workspace-app/internal/career/api"
	"workspace-app/internal/career/application"
)

// RegisterRoutes wires the career module. auth authenticates the caller.
func RegisterRoutes(mux *http.ServeMux, auth func(http.Handler) http.Handler) {
	api.RegisterRoutes(mux, api.NewHandler(application.NewService()), auth)
}
