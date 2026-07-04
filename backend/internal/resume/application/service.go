// Package application implements the Resume use cases: upload, parse, score,
// versioning, and improvement suggestions.
package application

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"workspace-app/internal/resume/domain"
)

// EventPublisher publishes domain events (the platform bus satisfies this).
type EventPublisher interface {
	Publish(ctx context.Context, eventType, aggregateID string, payload map[string]any) error
}

const (
	eventResumeUploaded = "ResumeUploaded"
	eventResumeParsed   = "ResumeParsed"
)

// ValidationError is returned for invalid input (mapped to HTTP 400).
type ValidationError struct{ Msg string }

func (e ValidationError) Error() string { return e.Msg }

// ErrForbidden indicates the caller does not own the resume.
var ErrForbidden = errors.New("forbidden")

type Service struct {
	repo    domain.Repository
	storage domain.Storage
	parser  domain.Parser
	events  EventPublisher
}

func NewService(repo domain.Repository, storage domain.Storage, parser domain.Parser, events EventPublisher) *Service {
	return &Service{repo: repo, storage: storage, parser: parser, events: events}
}

// UploadInput carries an uploaded file.
type UploadInput struct {
	Title       string
	Filename    string
	ContentType string
	Data        []byte
}

// Upload creates a new resume (version 1), parses and scores it.
func (s *Service) Upload(ctx context.Context, userID string, in UploadInput) (*domain.Resume, error) {
	if len(in.Data) == 0 {
		return nil, ValidationError{"uploaded file is empty"}
	}
	title := strings.TrimSpace(in.Title)
	if title == "" {
		title = in.Filename
	}
	if title == "" {
		title = "My Resume"
	}
	res := &domain.Resume{UserID: userID, Title: title}
	if err := s.repo.CreateResume(ctx, res); err != nil {
		return nil, err
	}
	if err := s.addVersion(ctx, res, in); err != nil {
		return nil, err
	}
	return s.detail(ctx, res.ID)
}

// AddVersion appends a new version to an owned resume.
func (s *Service) AddVersion(ctx context.Context, userID, resumeID string, in UploadInput) (*domain.Resume, error) {
	if len(in.Data) == 0 {
		return nil, ValidationError{"uploaded file is empty"}
	}
	res, err := s.requireOwner(ctx, userID, resumeID)
	if err != nil {
		return nil, err
	}
	if err := s.addVersion(ctx, res, in); err != nil {
		return nil, err
	}
	return s.detail(ctx, res.ID)
}

// addVersion stores the file, extracts text, persists the version and its score.
func (s *Service) addVersion(ctx context.Context, res *domain.Resume, in UploadInput) error {
	versionNo, err := s.repo.NextVersionNo(ctx, res.ID)
	if err != nil {
		return err
	}
	ext := filepath.Ext(in.Filename)
	key := fmt.Sprintf("%s/v%d%s", res.ID, versionNo, ext)
	if err := s.storage.Save(ctx, key, in.Data); err != nil {
		return err
	}

	text, perr := s.parser.ExtractText(in.Filename, in.ContentType, in.Data)
	if perr != nil && !errors.Is(perr, domain.ErrEmptyUpload) {
		// Unsupported type is a client error; other parse issues yield empty text.
		if errors.Is(perr, domain.ErrUnsupported) {
			return ValidationError{"unsupported file type — upload a PDF, DOCX, or TXT"}
		}
	}

	v := &domain.Version{
		ResumeID: res.ID, VersionNo: versionNo, Filename: in.Filename,
		ContentType: in.ContentType, SizeBytes: int64(len(in.Data)), StorageKey: key, ExtractedText: text,
	}
	if err := s.repo.AddVersion(ctx, v); err != nil {
		return err
	}

	score := domain.ScoreText(text)
	score.VersionID = v.ID
	if err := s.repo.SaveScore(ctx, &score); err != nil {
		return err
	}

	s.publish(ctx, eventResumeUploaded, res.ID, map[string]any{"user_id": res.UserID, "version_no": versionNo})
	s.publish(ctx, eventResumeParsed, res.ID, map[string]any{"version_id": v.ID, "overall": score.Overall})
	return nil
}

// List returns the user's resumes with their latest score attached.
func (s *Service) List(ctx context.Context, userID string) ([]domain.Resume, error) {
	resumes, err := s.repo.ListResumesByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	for i := range resumes {
		if score, err := s.repo.LatestScore(ctx, resumes[i].ID); err == nil {
			resumes[i].Score = score
		}
	}
	return resumes, nil
}

// Get returns a resume with its versions, latest version, and latest score.
func (s *Service) Get(ctx context.Context, userID, resumeID string) (*domain.Resume, error) {
	if _, err := s.requireOwner(ctx, userID, resumeID); err != nil {
		return nil, err
	}
	return s.detail(ctx, resumeID)
}

// Versions returns an owned resume's version history.
func (s *Service) Versions(ctx context.Context, userID, resumeID string) ([]domain.Version, error) {
	if _, err := s.requireOwner(ctx, userID, resumeID); err != nil {
		return nil, err
	}
	return s.repo.ListVersions(ctx, resumeID)
}

// Score returns the latest score for an owned resume.
func (s *Service) Score(ctx context.Context, userID, resumeID string) (*domain.Score, error) {
	if _, err := s.requireOwner(ctx, userID, resumeID); err != nil {
		return nil, err
	}
	return s.repo.LatestScore(ctx, resumeID)
}

// Review re-evaluates the latest version and returns its score + suggestions.
// (Heuristic for now; the AI module will enrich this later.)
func (s *Service) Review(ctx context.Context, userID, resumeID string) (*domain.Score, error) {
	if _, err := s.requireOwner(ctx, userID, resumeID); err != nil {
		return nil, err
	}
	v, err := s.repo.LatestVersion(ctx, resumeID)
	if err != nil {
		return nil, err
	}
	score := domain.ScoreText(v.ExtractedText)
	score.VersionID = v.ID
	if err := s.repo.SaveScore(ctx, &score); err != nil {
		return nil, err
	}
	return &score, nil
}

// Delete soft-deletes an owned resume.
func (s *Service) Delete(ctx context.Context, userID, resumeID string) error {
	if _, err := s.requireOwner(ctx, userID, resumeID); err != nil {
		return err
	}
	return s.repo.SoftDeleteResume(ctx, resumeID)
}

func (s *Service) detail(ctx context.Context, resumeID string) (*domain.Resume, error) {
	res, err := s.repo.GetResume(ctx, resumeID)
	if err != nil {
		return nil, err
	}
	if res.Versions, err = s.repo.ListVersions(ctx, resumeID); err != nil {
		return nil, err
	}
	if latest, err := s.repo.LatestVersion(ctx, resumeID); err == nil {
		res.Latest = latest
	}
	if score, err := s.repo.LatestScore(ctx, resumeID); err == nil {
		res.Score = score
	}
	return res, nil
}

func (s *Service) requireOwner(ctx context.Context, userID, resumeID string) (*domain.Resume, error) {
	res, err := s.repo.GetResume(ctx, resumeID)
	if err != nil {
		return nil, err
	}
	if res.UserID != userID {
		return nil, ErrForbidden
	}
	return res, nil
}

func (s *Service) publish(ctx context.Context, evt, aggID string, payload map[string]any) {
	if s.events != nil {
		_ = s.events.Publish(ctx, evt, aggID, payload)
	}
}
