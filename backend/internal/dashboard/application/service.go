// Package application implements the dashboard use cases. It depends only on the
// domain Repository port; the PostgreSQL adapter is injected in module.go.
package application

import (
	"context"

	"workspace-app/internal/dashboard/domain"
)

// Service is the dashboard query service.
type Service struct {
	repo domain.Repository
}

// NewService constructs the dashboard service.
func NewService(repo domain.Repository) *Service { return &Service{repo: repo} }

// Summary returns the aggregated per-user dashboard projection.
func (s *Service) Summary(ctx context.Context, userID string) (domain.Summary, error) {
	return s.repo.Summary(ctx, userID)
}
