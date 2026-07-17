package connections

import (
	"errors"
	"net/http"
	"workspace-app/internal/common"
)

var (
	ErrSelfRequest      = errors.New("cannot request connection to yourself")
	ErrAlreadyPending   = errors.New("connection request already pending")
	ErrAlreadyConnected = errors.New("already connected")
	ErrBlocked          = errors.New("cannot connect because a user is blocked")
	ErrCooldown         = errors.New("cooldown active after decline")
	ErrRateLimited      = errors.New("rate limit exceeded")
	ErrNotFound         = errors.New("connection not found")
	ErrForbidden        = errors.New("forbidden")
)

type ErrCooldownActive struct {
	RetryAfter int
}

func (e *ErrCooldownActive) Error() string {
	return "cooldown active after decline"
}

// MapError translates internal errors to their respective API AppError envelopes
func MapError(err error) error {
	if err == nil {
		return nil
	}

	var cooldownErr *ErrCooldownActive
	if errors.As(err, &cooldownErr) {
		return &common.AppError{
			Code:    "COOLDOWN_ACTIVE",
			Message: "A 30-day cooldown is active. You cannot request a connection yet.",
			Status:  http.StatusForbidden,
		}
	}

	switch {
	case errors.Is(err, ErrSelfRequest):
		return &common.AppError{
			Code:    "FORBIDDEN",
			Message: "You cannot send a connection request to yourself.",
			Status:  http.StatusForbidden,
		}
	case errors.Is(err, ErrAlreadyPending):
		return &common.AppError{
			Code:    "ALREADY_PENDING",
			Message: "A connection request is already pending for this user pair.",
			Status:  http.StatusConflict,
		}
	case errors.Is(err, ErrAlreadyConnected):
		return &common.AppError{
			Code:    "ALREADY_CONNECTED",
			Message: "You are already connected to this user.",
			Status:  http.StatusConflict,
		}
	case errors.Is(err, ErrBlocked):
		return &common.AppError{
			Code:    "BLOCKED",
			Message: "You cannot connect with this user.",
			Status:  http.StatusForbidden,
		}
	case errors.Is(err, ErrCooldown):
		return &common.AppError{
			Code:    "COOLDOWN_ACTIVE",
			Message: "A 30-day cooldown is active. You cannot request a connection yet.",
			Status:  http.StatusForbidden,
		}
	case errors.Is(err, ErrRateLimited):
		return &common.AppError{
			Code:    "RATE_LIMITED",
			Message: "Daily request limit exceeded. Please try again tomorrow.",
			Status:  http.StatusTooManyRequests,
		}
	case errors.Is(err, ErrNotFound):
		return &common.AppError{
			Code:    "NOT_FOUND",
			Message: "Connection or request not found.",
			Status:  http.StatusNotFound,
		}
	case errors.Is(err, ErrForbidden):
		return &common.AppError{
			Code:    "FORBIDDEN",
			Message: "You do not have permission to modify this connection.",
			Status:  http.StatusForbidden,
		}
	default:
		return err
	}
}
