package application

import (
	"context"
	"errors"
	"strings"

	"workspace-app/internal/identity/domain"
)

// OAuthAuthURL returns the provider authorization URL for the given state.
func (s *Service) OAuthAuthURL(provider, state string) (string, error) {
	p, ok := s.providers[provider]
	if !ok {
		return "", ErrUnknownProvider
	}
	return p.AuthURL(state), nil
}

// OAuthCallback exchanges an authorization code, provisions or links the local
// account, and issues a session. First OAuth login auto-creates the account.
func (s *Service) OAuthCallback(ctx context.Context, provider, code, ip string) (AuthResult, error) {
	p, ok := s.providers[provider]
	if !ok {
		return AuthResult{}, ErrUnknownProvider
	}
	uid, email, fullName, err := p.Exchange(ctx, code)
	if err != nil {
		return AuthResult{}, err
	}
	email = strings.TrimSpace(strings.ToLower(email))

	// 1) Already linked?
	if userID, found, err := s.oauth.FindUserIDByProvider(ctx, provider, uid); err != nil {
		return AuthResult{}, err
	} else if found {
		u, err := s.users.GetByID(ctx, userID)
		if err != nil {
			return AuthResult{}, err
		}
		return s.issueSession(ctx, u)
	}

	// 2) Existing account with same email? Otherwise create one.
	u, err := s.users.GetByEmail(ctx, email)
	if errors.Is(err, domain.ErrUserNotFound) {
		u = &domain.User{Email: email, FullName: fullName, EmailVerified: true, Status: domain.StatusActive}
		if err := s.users.Create(ctx, u); err != nil {
			return AuthResult{}, err
		}
		if err := s.users.AssignRole(ctx, u.ID, domain.RoleJobSeeker); err != nil {
			return AuthResult{}, err
		}
		u.Roles = []string{domain.RoleJobSeeker}
		_ = s.events.Publish(ctx, domain.EventUserRegistered, u.ID, map[string]any{"email": u.Email, "provider": provider})
	} else if err != nil {
		return AuthResult{}, err
	}

	// 3) Link the external identity.
	if err := s.oauth.Link(ctx, u.ID, provider, uid); err != nil {
		return AuthResult{}, err
	}
	s.record(ctx, u.ID, "user.oauth_login", "user", u.ID, ip)
	_ = s.events.Publish(ctx, domain.EventOAuthLinked, u.ID, map[string]any{"provider": provider})
	return s.issueSession(ctx, u)
}
