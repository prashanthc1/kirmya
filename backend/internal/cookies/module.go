package cookies

import (
	"database/sql"
	"net/http"

	"workspace-app/internal/cookies/api"
	"workspace-app/internal/cookies/application"
	"workspace-app/internal/cookies/infrastructure/postgres"
	"workspace-app/internal/identity/infrastructure/jwtauth"
)

type Module struct {
	Svc *application.Service
	h   *api.Handler
}

func NewModule(db *sql.DB) *Module {
	repo := postgres.NewRepository(db)
	svc := application.NewService(repo)
	h := api.NewHandler(svc)
	return &Module{
		Svc: svc,
		h:   h,
	}
}

func (m *Module) RegisterRoutes(mux *http.ServeMux, tokens *jwtauth.Factory) {
	routes := api.NewRoutes(m.h, tokens)
	routes.Register(mux)
}
