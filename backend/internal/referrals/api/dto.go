package api

import (
	"time"

	"workspace-app/internal/referrals/domain"
)

type referralResponse struct {
	ID         string `json:"id"`
	SeekerID   string `json:"seeker_id"`
	ReferrerID string `json:"referrer_id,omitempty"`
	JobID      string `json:"job_id,omitempty"`
	Company    string `json:"company,omitempty"`
	Message    string `json:"message,omitempty"`
	Status     string `json:"status"`
	Outcome    string `json:"outcome,omitempty"`
	DecidedAt  string `json:"decided_at,omitempty"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type requestRequest struct {
	ReferrerID string `json:"referrer_id"`
	JobID      string `json:"job_id"`
	Company    string `json:"company"`
	Message    string `json:"message"`
}

type outcomeRequest struct {
	Outcome string `json:"outcome"`
}

func toResponse(r *domain.Referral) referralResponse {
	resp := referralResponse{
		ID: r.ID, SeekerID: r.SeekerID, ReferrerID: r.ReferrerID, JobID: r.JobID,
		Company: r.Company, Message: r.Message, Status: r.Status, Outcome: r.Outcome,
		CreatedAt: r.CreatedAt.UTC().Format(time.RFC3339), UpdatedAt: r.UpdatedAt.UTC().Format(time.RFC3339),
	}
	if r.DecidedAt != nil {
		resp.DecidedAt = r.DecidedAt.UTC().Format(time.RFC3339)
	}
	return resp
}

func toResponses(refs []domain.Referral) []referralResponse {
	out := make([]referralResponse, 0, len(refs))
	for i := range refs {
		out = append(out, toResponse(&refs[i]))
	}
	return out
}
