package api

import (
	"time"

	"workspace-app/internal/jobs/domain"
)

type jobResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Company     string `json:"company"`
	Location    string `json:"location"`
	Description string `json:"description"`
	Salary      string `json:"salary"`
	JobType     string `json:"job_type"`
	PostedBy    string `json:"posted_by"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type applicationResponse struct {
	ID          string `json:"id"`
	JobID       string `json:"job_id"`
	UserID      string `json:"user_id"`
	Status      string `json:"status"`
	CoverLetter string `json:"cover_letter"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type postJobRequest struct {
	Title       string `json:"title"`
	Company     string `json:"company"`
	Location    string `json:"location"`
	Description string `json:"description"`
	Salary      string `json:"salary"`
	JobType     string `json:"job_type"`
}

type applyRequest struct {
	CoverLetter string `json:"cover_letter"`
}

type statusRequest struct {
	Status string `json:"status"`
}

type matchResponse struct {
	Job           jobResponse `json:"job"`
	Score         int         `json:"score"`
	MatchedSkills []string    `json:"matched_skills"`
	MissingSkills []string    `json:"missing_skills"`
	Reason        string      `json:"reason"`
}

func toMatch(m *domain.Match) matchResponse {
	return matchResponse{
		Job:           toJobResponse(&m.Job),
		Score:         m.Score,
		MatchedSkills: m.MatchedSkills,
		MissingSkills: m.MissingSkills,
		Reason:        m.Reason,
	}
}

func toMatches(matches []domain.Match) []matchResponse {
	out := make([]matchResponse, 0, len(matches))
	for i := range matches {
		out = append(out, toMatch(&matches[i]))
	}
	return out
}

func toJobResponse(j *domain.Job) jobResponse {
	return jobResponse{
		ID: j.ID, Title: j.Title, Company: j.Company, Location: j.Location,
		Description: j.Description, Salary: j.Salary, JobType: j.JobType, PostedBy: j.PostedBy,
		CreatedAt: j.CreatedAt.UTC().Format(time.RFC3339), UpdatedAt: j.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func toJobResponses(jobs []domain.Job) []jobResponse {
	out := make([]jobResponse, 0, len(jobs))
	for i := range jobs {
		out = append(out, toJobResponse(&jobs[i]))
	}
	return out
}

func toAppResponse(a *domain.Application) applicationResponse {
	return applicationResponse{
		ID: a.ID, JobID: a.JobID, UserID: a.UserID, Status: a.Status, CoverLetter: a.CoverLetter,
		CreatedAt: a.CreatedAt.UTC().Format(time.RFC3339), UpdatedAt: a.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func toAppResponses(apps []domain.Application) []applicationResponse {
	out := make([]applicationResponse, 0, len(apps))
	for i := range apps {
		out = append(out, toAppResponse(&apps[i]))
	}
	return out
}
