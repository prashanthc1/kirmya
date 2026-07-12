// Package postgres implements community/domain.Repository on PostgreSQL.
package postgres

import (
	"context"
	"database/sql"
	"errors"

	"workspace-app/internal/community/domain"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

func (r *Repository) ListCommunities(ctx context.Context) ([]domain.Community, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT c.id, c.slug, c.name, COALESCE(c.description,''), COALESCE(c.category,''), c.created_at,
		       (SELECT COUNT(*) FROM community_members m WHERE m.community_id = c.id)
		FROM communities c ORDER BY c.name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Community{}
	for rows.Next() {
		var c domain.Community
		if err := rows.Scan(&c.ID, &c.Slug, &c.Name, &c.Description, &c.Category, &c.CreatedAt, &c.MemberCount); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *Repository) GetBySlug(ctx context.Context, slug string) (*domain.Community, error) {
	var c domain.Community
	err := r.db.QueryRowContext(ctx, `
		SELECT c.id, c.slug, c.name, COALESCE(c.description,''), COALESCE(c.category,''), c.created_at,
		       (SELECT COUNT(*) FROM community_members m WHERE m.community_id = c.id)
		FROM communities c WHERE c.slug = $1`, slug).
		Scan(&c.ID, &c.Slug, &c.Name, &c.Description, &c.Category, &c.CreatedAt, &c.MemberCount)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repository) CreateCommunity(ctx context.Context, c *domain.Community, creatorUserID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, `
		INSERT INTO communities (slug, name, description, category)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`,
		c.Slug, c.Name, c.Description, c.Category).Scan(&c.ID, &c.CreatedAt)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO community_members (community_id, user_id, role)
		VALUES ($1, $2, 'moderator')`,
		c.ID, creatorUserID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) Join(ctx context.Context, communityID, userID string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO community_members (community_id, user_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`,
		communityID, userID)
	return err
}

func (r *Repository) Leave(ctx context.Context, communityID, userID string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM community_members WHERE community_id=$1 AND user_id=$2`, communityID, userID)
	return err
}

func (r *Repository) IsMember(ctx context.Context, communityID, userID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM community_members WHERE community_id=$1 AND user_id=$2)`,
		communityID, userID).Scan(&exists)
	return exists, err
}

const postCols = `p.id, p.community_id, p.author_id, p.title, COALESCE(p.body,''), p.created_at,
	(SELECT COUNT(*) FROM comments cm WHERE cm.post_id = p.id),
	(SELECT COUNT(*) FROM reactions rx WHERE rx.post_id = p.id)`

func scanPost(s interface{ Scan(...any) error }) (domain.Post, error) {
	var p domain.Post
	err := s.Scan(&p.ID, &p.CommunityID, &p.AuthorID, &p.Title, &p.Body, &p.CreatedAt, &p.CommentCount, &p.ReactionCount)
	return p, err
}

func (r *Repository) ListPosts(ctx context.Context, communityID string, limit int) ([]domain.Post, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT `+postCols+` FROM posts p WHERE p.community_id = $1 ORDER BY p.created_at DESC LIMIT $2`,
		communityID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Post{}
	for rows.Next() {
		p, err := scanPost(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *Repository) CreatePost(ctx context.Context, p *domain.Post) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO posts (community_id, author_id, title, body) VALUES ($1,$2,$3,NULLIF($4,'')) RETURNING id, created_at`,
		p.CommunityID, p.AuthorID, p.Title, p.Body).Scan(&p.ID, &p.CreatedAt)
}

func (r *Repository) GetPost(ctx context.Context, id string) (*domain.Post, error) {
	p, err := scanPost(r.db.QueryRowContext(ctx, `SELECT `+postCols+` FROM posts p WHERE p.id = $1`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) AddComment(ctx context.Context, c *domain.Comment) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO comments (post_id, author_id, body) VALUES ($1,$2,$3) RETURNING id, created_at`,
		c.PostID, c.AuthorID, c.Body).Scan(&c.ID, &c.CreatedAt)
}

func (r *Repository) ListComments(ctx context.Context, postID string) ([]domain.Comment, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, post_id, author_id, body, created_at FROM comments WHERE post_id=$1 ORDER BY created_at`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Comment{}
	for rows.Next() {
		var c domain.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.AuthorID, &c.Body, &c.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *Repository) ToggleReaction(ctx context.Context, postID, userID, kind string) (bool, error) {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM reactions WHERE post_id=$1 AND user_id=$2 AND kind=$3`, postID, userID, kind)
	if err != nil {
		return false, err
	}
	if n, _ := res.RowsAffected(); n > 0 {
		return false, nil // removed an existing reaction
	}
	_, err = r.db.ExecContext(ctx,
		`INSERT INTO reactions (post_id, user_id, kind) VALUES ($1,$2,$3) ON CONFLICT DO NOTHING`,
		postID, userID, kind)
	return true, err
}

func (r *Repository) MemberRole(ctx context.Context, communityID, userID string) (string, error) {
	var role string
	err := r.db.QueryRowContext(ctx,
		`SELECT role FROM community_members WHERE community_id=$1 AND user_id=$2`, communityID, userID).Scan(&role)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return role, nil
}

func (r *Repository) DeletePost(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM posts WHERE id=$1`, id)
	return err
}

// --- tags ---

func (r *Repository) SetPostTags(ctx context.Context, postID string, tags []string) error {
	if _, err := r.db.ExecContext(ctx, `DELETE FROM post_tags WHERE post_id=$1`, postID); err != nil {
		return err
	}
	for _, t := range tags {
		if _, err := r.db.ExecContext(ctx,
			`INSERT INTO post_tags (post_id, tag) VALUES ($1,$2) ON CONFLICT DO NOTHING`, postID, t); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) ListTags(ctx context.Context, communityID string) ([]domain.Tag, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT pt.tag, COUNT(*)
		FROM post_tags pt
		JOIN posts p ON p.id = pt.post_id
		WHERE p.community_id = $1
		GROUP BY pt.tag
		ORDER BY COUNT(*) DESC, pt.tag`, communityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Tag{}
	for rows.Next() {
		var t domain.Tag
		if err := rows.Scan(&t.Name, &t.Count); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *Repository) ListPostsByTag(ctx context.Context, communityID, tag string, limit int) ([]domain.Post, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT `+postCols+` FROM posts p
		 JOIN post_tags pt ON pt.post_id = p.id
		 WHERE p.community_id = $1 AND pt.tag = $2
		 ORDER BY p.created_at DESC LIMIT $3`,
		communityID, tag, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Post{}
	for rows.Next() {
		p, err := scanPost(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// --- polls ---

func (r *Repository) CreatePoll(ctx context.Context, p *domain.Poll) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if err := tx.QueryRowContext(ctx,
		`INSERT INTO polls (post_id, question) VALUES ($1,$2) RETURNING id, created_at`,
		p.PostID, p.Question).Scan(&p.ID, &p.CreatedAt); err != nil {
		return err
	}
	for i := range p.Options {
		p.Options[i].PollID = p.ID
		if err := tx.QueryRowContext(ctx,
			`INSERT INTO poll_options (poll_id, label) VALUES ($1,$2) RETURNING id`,
			p.ID, p.Options[i].Label).Scan(&p.Options[i].ID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *Repository) GetPollByPost(ctx context.Context, postID string) (*domain.Poll, error) {
	var p domain.Poll
	err := r.db.QueryRowContext(ctx,
		`SELECT id, post_id, question, created_at FROM polls WHERE post_id=$1`, postID).
		Scan(&p.ID, &p.PostID, &p.Question, &p.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrPollNotFound
	}
	if err != nil {
		return nil, err
	}
	if err := r.loadOptions(ctx, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) GetPoll(ctx context.Context, pollID string) (*domain.Poll, error) {
	var p domain.Poll
	err := r.db.QueryRowContext(ctx,
		`SELECT id, post_id, question, created_at FROM polls WHERE id=$1`, pollID).
		Scan(&p.ID, &p.PostID, &p.Question, &p.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrPollNotFound
	}
	if err != nil {
		return nil, err
	}
	if err := r.loadOptions(ctx, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) loadOptions(ctx context.Context, p *domain.Poll) error {
	rows, err := r.db.QueryContext(ctx, `
		SELECT o.id, o.poll_id, o.label,
		       (SELECT COUNT(*) FROM poll_votes v WHERE v.option_id = o.id)
		FROM poll_options o WHERE o.poll_id = $1 ORDER BY o.id`, p.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	p.Options = []domain.PollOption{}
	for rows.Next() {
		var o domain.PollOption
		if err := rows.Scan(&o.ID, &o.PollID, &o.Label, &o.VoteCount); err != nil {
			return err
		}
		p.Options = append(p.Options, o)
	}
	return rows.Err()
}

func (r *Repository) Vote(ctx context.Context, pollID, optionID, userID string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO poll_votes (poll_id, option_id, user_id) VALUES ($1,$2,$3)
		ON CONFLICT (poll_id, user_id) DO UPDATE SET option_id = EXCLUDED.option_id`,
		pollID, optionID, userID)
	return err
}

// --- reports ---

func (r *Repository) CreateReport(ctx context.Context, rep *domain.Report) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO content_reports (reporter_id, target_type, target_id, reason)
		VALUES ($1,$2,$3,$4)
		RETURNING id, status, created_at`,
		rep.ReporterID, rep.TargetType, rep.TargetID, rep.Reason).
		Scan(&rep.ID, &rep.Status, &rep.CreatedAt)
}

func (r *Repository) ListOpenReports(ctx context.Context, communityID string) ([]domain.Report, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT cr.id, cr.reporter_id, cr.target_type, cr.target_id, cr.reason, cr.status, cr.created_at
		FROM content_reports cr
		JOIN posts p ON p.id = cr.target_id
		WHERE cr.target_type = 'post' AND cr.status = 'open' AND p.community_id = $1
		ORDER BY cr.created_at DESC`, communityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Report{}
	for rows.Next() {
		var rep domain.Report
		if err := rows.Scan(&rep.ID, &rep.ReporterID, &rep.TargetType, &rep.TargetID, &rep.Reason, &rep.Status, &rep.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, rep)
	}
	return out, rows.Err()
}
