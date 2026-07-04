// Package domain holds the Resume bounded context entities, ports, and the
// pure resume-scoring logic.
package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound    = errors.New("resume not found")
	ErrUnsupported = errors.New("unsupported file type")
	ErrEmptyUpload = errors.New("uploaded file is empty")
)

type Resume struct {
	ID        string
	UserID    string
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
	Versions  []Version // optional, loaded on detail
	Latest    *Version  // optional
	Score     *Score    // optional, latest score
}

type Version struct {
	ID            string
	ResumeID      string
	VersionNo     int
	Filename      string
	ContentType   string
	SizeBytes     int64
	StorageKey    string
	ExtractedText string
	CreatedAt     time.Time
}

// Score is a parsed resume's assessment (each sub-score 0-100).
type Score struct {
	VersionID   string
	Overall     int
	Formatting  int
	Keywords    int
	ATS         int
	Suggestions []string
	CreatedAt   time.Time
}

// Repository is the persistence port for resumes/versions/scores.
type Repository interface {
	CreateResume(ctx context.Context, r *Resume) error
	GetResume(ctx context.Context, id string) (*Resume, error)
	ListResumesByUser(ctx context.Context, userID string) ([]Resume, error)
	SoftDeleteResume(ctx context.Context, id string) error

	NextVersionNo(ctx context.Context, resumeID string) (int, error)
	AddVersion(ctx context.Context, v *Version) error
	ListVersions(ctx context.Context, resumeID string) ([]Version, error)
	LatestVersion(ctx context.Context, resumeID string) (*Version, error)

	SaveScore(ctx context.Context, s *Score) error
	LatestScore(ctx context.Context, resumeID string) (*Score, error)
}

// Storage persists raw resume bytes (filesystem in MVP, S3 later).
type Storage interface {
	Save(ctx context.Context, key string, data []byte) error
}

// Parser extracts plain text from an uploaded resume.
type Parser interface {
	ExtractText(filename, contentType string, data []byte) (string, error)
}
