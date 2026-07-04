// Package domain holds the Mentorship bounded context.
package domain

import (
	"errors"
	"time"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrNotMentor       = errors.New("you are not the mentor for this session")
	ErrNotMentee       = errors.New("you are not the mentee for this session")
	ErrNotComplete     = errors.New("session must be completed before review")
	ErrSlotNotFound    = errors.New("availability slot not found")
	ErrSlotUnavailable = errors.New("availability slot is already booked")
	ErrSlotNotOwned    = errors.New("you do not own this availability slot")
)

const (
	StatusRequested = "requested"
	StatusConfirmed = "confirmed"
	StatusCompleted = "completed"
	StatusCancelled = "cancelled"
)

var ValidStatuses = map[string]bool{
	StatusRequested: true, StatusConfirmed: true, StatusCompleted: true, StatusCancelled: true,
}

type MentorProfile struct {
	ID        string
	UserID    string
	Headline  string
	Bio       string
	Expertise string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Session struct {
	ID          string
	MentorID    string
	MenteeID    string
	Topic       string
	Status      string
	ScheduledAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Review struct {
	ID        string
	SessionID string
	Rating    int
	Comment   string
	CreatedAt time.Time
}

// AvailabilitySlot is a window a mentor opens for booking. A slot is consumed
// (IsBooked=true) when a mentee books a session against it.
type AvailabilitySlot struct {
	ID        string
	MentorID  string
	StartsAt  time.Time
	EndsAt    time.Time
	IsBooked  bool
	CreatedAt time.Time
}

// The Repository port lives in ports.go (gold-standard DDD layout).
