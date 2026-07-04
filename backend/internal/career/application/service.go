// Package application implements the career use cases. Career data is curated
// reference data owned by the domain, so the service is a thin query facade.
package application

import (
	"context"

	"workspace-app/internal/career/domain"
)

// Service is the career query service.
type Service struct{}

// NewService constructs the career service.
func NewService() *Service { return &Service{} }

// Paths returns the career ladder for the given starting role.
func (s *Service) Paths(_ context.Context, from string) domain.Path {
	return domain.Paths(from)
}
