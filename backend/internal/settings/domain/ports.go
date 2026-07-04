package domain

import "context"

// Repository is the persistence port for the Settings context. The concrete
// adapter lives in infrastructure/postgres.
type Repository interface {
	// Get returns the user's settings, or ErrNotFound if no row exists yet.
	Get(ctx context.Context, userID string) (*UserSettings, error)

	// EnsureDefaults inserts the default row for a user if one does not already
	// exist and returns the current settings. It is idempotent.
	EnsureDefaults(ctx context.Context, userID string) (*UserSettings, error)

	// Update persists the aggregate. It is version-checked and returns
	// ErrOptimisticLock on a stale write.
	Update(ctx context.Context, s *UserSettings) error
}
