package postgres

import (
	"context"
	"database/sql"
)

// SkillReader reads a seeker's profile skills for job matching. It lives in the
// jobs context as a read-only projection over the profile-owned skills tables.
type SkillReader struct{ db *sql.DB }

func NewSkillReader(db *sql.DB) *SkillReader { return &SkillReader{db: db} }

func (r *SkillReader) SeekerSkills(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT s.name
		FROM profile_skills ps
		JOIN skills s ON s.id = ps.skill_id
		WHERE ps.user_id = $1
		ORDER BY s.name`, userID)
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
