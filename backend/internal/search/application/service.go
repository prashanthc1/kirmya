// Package application implements the Search use cases: querying (with a database
// fallback when the engine is unavailable) and keeping the index fresh by
// re-indexing entities in response to domain events.
package application

import (
	"context"
	"log"
	"strings"

	"workspace-app/internal/platform/search"
)

// Row types are the minimal entity projections the Source returns; the service
// maps them into search.Doc values.
type UserRow struct {
	ID, FullName, Email, Headline string
}

type JobRow struct {
	ID, Title, Company, Location, Description string
}

type CommunityRow struct {
	ID, Slug, Name, Description, Category string
}

// Source reads entities from the database for (re)indexing and provides the
// ILIKE fallback search used when the engine is not ready.
type Source interface {
	User(ctx context.Context, id string) (*UserRow, error)
	Job(ctx context.Context, id string) (*JobRow, error)
	AllUsers(ctx context.Context) ([]UserRow, error)
	AllJobs(ctx context.Context) ([]JobRow, error)
	AllCommunities(ctx context.Context) ([]CommunityRow, error)
	AllSkills(ctx context.Context) ([]string, error)
	FallbackSearch(ctx context.Context, query string, types []string, limit int) ([]search.Hit, error)
}

type Service struct {
	engine search.Engine
	src    Source
}

func NewService(engine search.Engine, src Source) *Service {
	return &Service{engine: engine, src: src}
}

// Ready reports whether the full-text engine backs queries (vs. DB fallback).
func (s *Service) Ready() bool { return s.engine.Ready() }

// Query runs a fuzzy search, falling back to the database when the engine is
// unavailable or errors.
func (s *Service) Query(ctx context.Context, q string, types []string, limit int) ([]search.Hit, error) {
	q = strings.TrimSpace(q)
	if q == "" {
		return []search.Hit{}, nil
	}
	limit = clampLimit(limit)
	if s.engine.Ready() {
		if hits, err := s.engine.Search(ctx, q, types, limit); err == nil {
			return hits, nil
		} else {
			log.Printf("[search] engine query failed, falling back to DB: %v", err)
		}
	}
	return s.src.FallbackSearch(ctx, q, types, limit)
}

// Suggest powers autocomplete (prefix match), with the same DB fallback.
func (s *Service) Suggest(ctx context.Context, q string, limit int) ([]search.Hit, error) {
	q = strings.TrimSpace(q)
	if q == "" {
		return []search.Hit{}, nil
	}
	limit = clampLimit(limit)
	if s.engine.Ready() {
		if hits, err := s.engine.Suggest(ctx, q, limit); err == nil {
			return hits, nil
		} else {
			log.Printf("[search] engine suggest failed, falling back to DB: %v", err)
		}
	}
	return s.src.FallbackSearch(ctx, q, nil, limit)
}

// ----- Indexing ----------------------------------------------------------

// IndexUser (re)indexes a single user by id (best-effort).
func (s *Service) IndexUser(ctx context.Context, id string) {
	row, err := s.src.User(ctx, id)
	if err != nil || row == nil {
		return
	}
	s.engine.Index(ctx, userDoc(*row))
}

// IndexJob (re)indexes a single job by id (best-effort).
func (s *Service) IndexJob(ctx context.Context, id string) {
	row, err := s.src.Job(ctx, id)
	if err != nil || row == nil {
		return
	}
	s.engine.Index(ctx, jobDoc(*row))
}

// Backfill indexes every entity. Run once on startup; upserts are idempotent.
func (s *Service) Backfill(ctx context.Context) {
	if !s.engine.Ready() {
		return
	}
	if users, err := s.src.AllUsers(ctx); err == nil {
		for _, u := range users {
			s.engine.Index(ctx, userDoc(u))
		}
	}
	if jobs, err := s.src.AllJobs(ctx); err == nil {
		for _, j := range jobs {
			s.engine.Index(ctx, jobDoc(j))
		}
	}
	if comms, err := s.src.AllCommunities(ctx); err == nil {
		for _, c := range comms {
			s.engine.Index(ctx, communityDoc(c))
		}
	}
	if skills, err := s.src.AllSkills(ctx); err == nil {
		for _, sk := range skills {
			s.engine.Index(ctx, skillDoc(sk))
		}
	}
	log.Printf("[search] backfill complete")
}

// ----- Doc mapping -------------------------------------------------------

func userDoc(u UserRow) search.Doc {
	title := u.FullName
	if title == "" {
		title = u.Email
	}
	return search.Doc{
		Type: "user", RefID: u.ID, Title: title, Subtitle: u.Headline,
		Body: u.Email + " " + u.Headline, URL: "/profile/" + u.ID,
	}
}

func jobDoc(j JobRow) search.Doc {
	return search.Doc{
		Type: "job", RefID: j.ID, Title: j.Title,
		Subtitle: strings.TrimSpace(j.Company + " · " + j.Location),
		Body:     j.Description, URL: "/jobs",
	}
}

func communityDoc(c CommunityRow) search.Doc {
	return search.Doc{
		Type: "community", RefID: c.ID, Title: c.Name, Subtitle: c.Category,
		Body: c.Description, URL: "/communities/" + c.Slug,
	}
}

func skillDoc(skill string) search.Doc {
	return search.Doc{
		Type: "skill", RefID: strings.ToLower(skill), Title: skill,
		Subtitle: "Skill", URL: "/jobs?q=" + skill,
	}
}

func clampLimit(limit int) int {
	if limit <= 0 {
		return 20
	}
	if limit > 50 {
		return 50
	}
	return limit
}
