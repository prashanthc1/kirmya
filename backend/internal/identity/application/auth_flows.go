package application

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"workspace-app/internal/identity/domain"
)

// newFamilyID returns a random refresh-token family identifier.
func newFamilyID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// sendVerification generates, stores and emails an email-verification token.
func (s *Service) sendVerification(ctx context.Context, u *domain.User) error {
	raw, hash, err := s.tokens.GenerateOpaqueToken()
	if err != nil {
		return err
	}
	if err := s.verif.StoreEmailToken(ctx, u.ID, hash, time.Now().Add(verificationTTL)); err != nil {
		return err
	}
	return s.mailer.SendVerificationEmail(ctx, u.Email, raw)
}

// VerifyEmail consumes a verification token and marks the email verified.
func (s *Service) VerifyEmail(ctx context.Context, rawToken string) error {
	userID, err := s.verif.ConsumeEmailToken(ctx, s.tokens.HashOpaqueToken(rawToken))
	if errors.Is(err, domain.ErrTokenNotFound) {
		return ErrInvalidCredentials
	}
	if err != nil {
		return err
	}
	if err := s.users.SetEmailVerified(ctx, userID); err != nil {
		return err
	}
	s.record(ctx, userID, "user.verify_email", "user", userID, "")
	_ = s.events.Publish(ctx, domain.EventEmailVerified, userID, nil)
	return nil
}

// ResendVerification re-issues a verification email if the account exists and
// is unverified. It never reveals whether the email exists.
func (s *Service) ResendVerification(ctx context.Context, email string) error {
	u, err := s.users.GetByEmail(ctx, email)
	if errors.Is(err, domain.ErrUserNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	if u.EmailVerified {
		return nil
	}
	return s.sendVerification(ctx, u)
}

// ForgotPassword emails a reset link if the account exists. Always returns nil
// to avoid account enumeration.
func (s *Service) ForgotPassword(ctx context.Context, email string) error {
	u, err := s.users.GetByEmail(ctx, email)
	if errors.Is(err, domain.ErrUserNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	raw, hash, err := s.tokens.GenerateOpaqueToken()
	if err != nil {
		return err
	}
	if err := s.verif.StorePasswordToken(ctx, u.ID, hash, time.Now().Add(resetTTL)); err != nil {
		return err
	}
	return s.mailer.SendPasswordResetEmail(ctx, u.Email, raw)
}

// ResetPassword consumes a reset token, sets a new password, and revokes every
// existing refresh token so a reset (often triggered because the old password
// may be compromised) invalidates all active sessions on every device.
func (s *Service) ResetPassword(ctx context.Context, rawToken, newPassword string) error {
	userID, err := s.verif.ConsumePasswordToken(ctx, s.tokens.HashOpaqueToken(rawToken))
	if errors.Is(err, domain.ErrTokenNotFound) {
		return ErrInvalidCredentials
	}
	if err != nil {
		return err
	}
	hash, err := s.hasher.Hash(newPassword)
	if err != nil {
		return err
	}
	if err := s.users.SetPasswordHash(ctx, userID, hash); err != nil {
		return err
	}
	// Force re-authentication everywhere; mirrors ChangePassword's behaviour.
	if err := s.refresh.RevokeAllForUser(ctx, userID); err != nil {
		return err
	}
	s.record(ctx, userID, "user.reset_password", "user", userID, "")
	_ = s.events.Publish(ctx, domain.EventPasswordReset, userID, nil)
	return nil
}

// SetupMFA provisions a TOTP secret (unconfirmed) and returns the otpauth URL
// for QR enrollment.
func (s *Service) SetupMFA(ctx context.Context, userID string) (otpauthURL string, err error) {
	u, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return "", err
	}
	secretEnc, url, err := s.totp.Generate(u.Email)
	if err != nil {
		return "", err
	}
	if err := s.mfa.Upsert(ctx, &domain.MFACredential{UserID: userID, SecretEnc: secretEnc}); err != nil {
		return "", err
	}
	return url, nil
}

// ConfirmMFA validates the first TOTP code and enables MFA on the account.
func (s *Service) ConfirmMFA(ctx context.Context, userID, code string) error {
	cred, err := s.mfa.Get(ctx, userID)
	if err != nil {
		return ErrInvalidMFACode
	}
	if !s.totp.Validate(cred.SecretEnc, code) {
		return ErrInvalidMFACode
	}
	if err := s.mfa.Confirm(ctx, userID); err != nil {
		return err
	}
	if err := s.users.SetMFAEnabled(ctx, userID, true); err != nil {
		return err
	}
	s.record(ctx, userID, "user.enable_mfa", "user", userID, "")
	_ = s.events.Publish(ctx, domain.EventMFAEnabled, userID, nil)
	return nil
}

// DisableMFA turns off TOTP for an account after validating a current code, so
// only someone in possession of the authenticator can disable the second factor.
// It is a no-op (returns nil) if MFA is already disabled.
func (s *Service) DisableMFA(ctx context.Context, userID, code string) error {
	u, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if !u.MFAEnabled {
		return nil
	}
	cred, err := s.mfa.Get(ctx, userID)
	if err != nil || !s.totp.Validate(cred.SecretEnc, code) {
		return ErrInvalidMFACode
	}
	if err := s.users.SetMFAEnabled(ctx, userID, false); err != nil {
		return err
	}
	s.record(ctx, userID, "user.disable_mfa", "user", userID, "")
	_ = s.events.Publish(ctx, domain.EventMFADisabled, userID, nil)
	return nil
}

// ChangePassword updates the password for an authenticated user after verifying
// their current password, then revokes every refresh token so all other
// sessions are forced to re-authenticate. The new password is validated by the
// api layer (length/complexity) before this is called.
func (s *Service) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	u, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if !u.HasPassword() {
		// OAuth-only accounts have no password to change.
		return ErrInvalidCredentials
	}
	ok, err := s.hasher.Verify(currentPassword, u.PasswordHash)
	if err != nil || !ok {
		return ErrInvalidCredentials
	}
	hash, err := s.hasher.Hash(newPassword)
	if err != nil {
		return err
	}
	if err := s.users.SetPasswordHash(ctx, userID, hash); err != nil {
		return err
	}
	// Force re-auth everywhere: a password change should invalidate sessions an
	// attacker (or an old device) might still hold.
	if err := s.refresh.RevokeAllForUser(ctx, userID); err != nil {
		return err
	}
	s.record(ctx, userID, "user.change_password", "user", userID, "")
	_ = s.events.Publish(ctx, domain.EventPasswordChanged, userID, nil)
	return nil
}
