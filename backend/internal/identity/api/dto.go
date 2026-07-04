package api

import (
	"errors"
	"net/mail"
	"strings"

	"workspace-app/internal/identity/domain"
)

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

// selfRegisterableRoles is the allowlist of roles a caller may grant itself via
// the public registration endpoint. Privileged roles (recruiter, admin) must be
// provisioned by an admin and are never self-assignable — the API, not the
// frontend, is the trust boundary here.
var selfRegisterableRoles = map[string]bool{
	domain.RoleJobSeeker: true,
	domain.RoleReferrer:  true,
	domain.RoleMentor:    true,
}

func (r registerRequest) validate() error {
	if _, err := mail.ParseAddress(strings.TrimSpace(r.Email)); err != nil {
		return errors.New("a valid email is required")
	}
	if len(r.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	// Empty role is allowed and defaults to job_seeker in the service layer.
	if role := strings.TrimSpace(r.Role); role != "" && !selfRegisterableRoles[role] {
		return errors.New("role must be one of: job_seeker, referrer, mentor")
	}
	return nil
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     string `json:"code"` // optional TOTP for MFA challenge
}

type tokenRequest struct {
	Token string `json:"token"`
}

type emailRequest struct {
	Email string `json:"email"`
}

type resetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type oauthCallbackRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

type mfaVerifyRequest struct {
	Code string `json:"code"`
}
