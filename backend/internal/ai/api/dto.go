package api

import (
	"time"

	"workspace-app/internal/ai/domain"
)

type resumeReviewRequest struct {
	ResumeText string `json:"resume_text"`
}

type skillGapRequest struct {
	CurrentRole   string   `json:"current_role"`
	TargetRole    string   `json:"target_role"`
	CurrentSkills []string `json:"current_skills"`
}

type coachRequest struct {
	ThreadID string `json:"thread_id"`
	Message  string `json:"message"`
}

type threadSummary struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type messageDTO struct {
	ID        string `json:"id"`
	Role      string `json:"role"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type threadDetail struct {
	threadSummary
	Messages []messageDTO `json:"messages"`
}

func toThreadSummary(t *domain.CoachThread) threadSummary {
	return threadSummary{
		ID: t.ID, Title: t.Title,
		CreatedAt: t.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: t.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func toThreadDetail(t *domain.CoachThread) threadDetail {
	d := threadDetail{threadSummary: toThreadSummary(t)}
	for _, m := range t.Messages {
		d.Messages = append(d.Messages, messageDTO{
			ID: m.ID, Role: m.Role, Content: m.Content,
			CreatedAt: m.CreatedAt.UTC().Format(time.RFC3339),
		})
	}
	return d
}
