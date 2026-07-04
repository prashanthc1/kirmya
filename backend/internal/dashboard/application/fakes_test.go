package application

import (
	"context"

	"workspace-app/internal/dashboard/domain"
)

// fakeRepo is an in-memory dashboard Repository for unit tests.
type fakeRepo struct {
	summary domain.Summary
	err     error
	gotUser string
}

func (f *fakeRepo) Summary(_ context.Context, userID string) (domain.Summary, error) {
	f.gotUser = userID
	return f.summary, f.err
}
