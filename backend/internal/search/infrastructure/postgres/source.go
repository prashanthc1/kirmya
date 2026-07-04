// Package postgres implements search/application.Source on PostgreSQL. It reads
// across context-owned tables (users/profiles, jobs, communities, skills) to
// build search documents and provides the ILIKE fallback used when OpenSearch
// is unavailable.
package postgres

import (
	"context"
	"database/sql"
	"errors"

	"workspace-app/internal/platform/search"
	"workspace-app/internal/search/application"
)

type Source struct{ db *sql.DB }

func NewSource(db *sql.DB) *Source { return &Source{db: db} }

func (s *Source) User(ctx context.Context, id string) (*application.UserRow, error) {
	var u application.UserRow
	err := s.db.QueryRowContext(ctx, `
		SELECT u.id, u.full_name, u.email, COALESCE(p.headline,'')
		FROM users u LEFT JOIN profiles p ON p.user_id = u.id
		WHERE u.id = $1 AND u.deleted_at IS NULL`, id).
		Scan(&u.ID, &u.FullName, &u.Email, &u.Headline)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Source) Job(ctx context.Context, id string) (*application.JobRow, error) {
	var j application.JobRow
	err := s.db.QueryRowContext(ctx, `
		SELECT id, title, company, COALESCE(location,''), COALESCE(description,'')
		FROM jobs WHERE id = $1`, id).
		Scan(&j.ID, &j.Title, &j.Company, &j.Location, &j.Description)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &j, nil
}

func (s *Source) AllUsers(ctx context.Context) ([]application.UserRow, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT u.id, u.full_name, u.email, COALESCE(p.headline,'')
		FROM users u LEFT JOIN profiles p ON p.user_id = u.id
		WHERE u.deleted_at IS NULL`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []application.UserRow{}
	for rows.Next() {
		var u application.UserRow
		if err := rows.Scan(&u.ID, &u.FullName, &u.Email, &u.Headline); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

func (s *Source) AllJobs(ctx context.Context) ([]application.JobRow, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, title, company, COALESCE(location,''), COALESCE(description,'') FROM jobs`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []application.JobRow{}
	for rows.Next() {
		var j application.JobRow
		if err := rows.Scan(&j.ID, &j.Title, &j.Company, &j.Location, &j.Description); err != nil {
			return nil, err
		}
		out = append(out, j)
	}
	return out, rows.Err()
}

func (s *Source) AllCommunities(ctx context.Context) ([]application.CommunityRow, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, slug, name, COALESCE(description,''), COALESCE(category,'') FROM communities`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []application.CommunityRow{}
	for rows.Next() {
		var c application.CommunityRow
		if err := rows.Scan(&c.ID, &c.Slug, &c.Name, &c.Description, &c.Category); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (s *Source) AllSkills(ctx context.Context) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT name FROM skills ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []string{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		out = append(out, name)
	}
	return out, rows.Err()
}

// FallbackSearch runs ILIKE queries across the requested types (all when none
// are given) and merges the results up to limit. It mirrors what the engine
// returns so callers are agnostic to which path served the query.
func (s *Source) FallbackSearch(ctx context.Context, query string, types []string, limit int) ([]search.Hit, error) {
	want := func(t string) bool {
		if len(types) == 0 {
			return true
		}
		for _, x := range types {
			if x == t {
				return true
			}
		}
		return false
	}
	pat := "%" + query + "%"
	hits := []search.Hit{}

	if want("user") {
		rows, err := s.db.QueryContext(ctx, `
			SELECT u.id, u.full_name, u.email, COALESCE(p.headline,'')
			FROM users u LEFT JOIN profiles p ON p.user_id = u.id
			WHERE u.deleted_at IS NULL AND (u.full_name ILIKE $1 OR u.email ILIKE $1 OR p.headline ILIKE $1)
			ORDER BY u.full_name NULLS LAST LIMIT $2`, pat, limit)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var id, name, email, headline string
			if err := rows.Scan(&id, &name, &email, &headline); err != nil {
				rows.Close()
				return nil, err
			}
			title := name
			if title == "" {
				title = email
			}
			hits = append(hits, search.Hit{Type: "user", RefID: id, Title: title, Subtitle: headline})
		}
		rows.Close()
	}

	if want("job") {
		rows, err := s.db.QueryContext(ctx, `
			SELECT id, title, company, COALESCE(location,'')
			FROM jobs WHERE title ILIKE $1 OR company ILIKE $1 OR description ILIKE $1
			ORDER BY created_at DESC LIMIT $2`, pat, limit)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var id, title, company, location string
			if err := rows.Scan(&id, &title, &company, &location); err != nil {
				rows.Close()
				return nil, err
			}
			hits = append(hits, search.Hit{Type: "job", RefID: id, Title: title, Subtitle: company + " · " + location, URL: "/jobs"})
		}
		rows.Close()
	}

	if want("community") {
		rows, err := s.db.QueryContext(ctx, `
			SELECT id, slug, name, COALESCE(category,'')
			FROM communities WHERE name ILIKE $1 OR description ILIKE $1 OR category ILIKE $1
			LIMIT $2`, pat, limit)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var id, slug, name, category string
			if err := rows.Scan(&id, &slug, &name, &category); err != nil {
				rows.Close()
				return nil, err
			}
			hits = append(hits, search.Hit{Type: "community", RefID: id, Title: name, Subtitle: category, URL: "/communities/" + slug})
		}
		rows.Close()
	}

	if want("skill") {
		rows, err := s.db.QueryContext(ctx, `
			SELECT name FROM skills WHERE name ILIKE $1 ORDER BY name LIMIT $2`, pat, limit)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				rows.Close()
				return nil, err
			}
			hits = append(hits, search.Hit{Type: "skill", RefID: name, Title: name, Subtitle: "Skill", URL: "/jobs?q=" + name})
		}
		rows.Close()
	}

	if len(hits) > limit {
		hits = hits[:limit]
	}
	return hits, nil
}
