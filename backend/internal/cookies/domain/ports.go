package domain

import "context"

// Repository represents the persistence port for CookiePreferences.
type Repository interface {
	GetByUserID(ctx context.Context, userID string) (*CookiePreferences, error)
	GetByAnonymousID(ctx context.Context, anonymousID string) (*CookiePreferences, error)
	Save(ctx context.Context, p *CookiePreferences) error
	Delete(ctx context.Context, userID string, anonymousID string) error
	Merge(ctx context.Context, anonymousID string, userID string) error
}
