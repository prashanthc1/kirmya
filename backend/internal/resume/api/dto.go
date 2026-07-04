package api

import (
	"time"

	"workspace-app/internal/resume/domain"
)

type resumeResponse struct {
	ID        string            `json:"id"`
	Title     string            `json:"title"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
	Latest    *versionResponse  `json:"latest_version,omitempty"`
	Score     *scoreResponse    `json:"score,omitempty"`
	Versions  []versionResponse `json:"versions,omitempty"`
}

type versionResponse struct {
	ID          string `json:"id"`
	VersionNo   int    `json:"version_no"`
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	SizeBytes   int64  `json:"size_bytes"`
	CreatedAt   string `json:"created_at"`
}

type scoreResponse struct {
	Overall     int      `json:"overall"`
	Formatting  int      `json:"formatting"`
	Keywords    int      `json:"keywords"`
	ATS         int      `json:"ats"`
	Suggestions []string `json:"suggestions"`
}

func toVersion(v *domain.Version) versionResponse {
	return versionResponse{
		ID: v.ID, VersionNo: v.VersionNo, Filename: v.Filename, ContentType: v.ContentType,
		SizeBytes: v.SizeBytes, CreatedAt: v.CreatedAt.UTC().Format(time.RFC3339),
	}
}

func toScore(s *domain.Score) *scoreResponse {
	if s == nil {
		return nil
	}
	sugg := s.Suggestions
	if sugg == nil {
		sugg = []string{}
	}
	return &scoreResponse{Overall: s.Overall, Formatting: s.Formatting, Keywords: s.Keywords, ATS: s.ATS, Suggestions: sugg}
}

func toResume(r *domain.Resume) resumeResponse {
	resp := resumeResponse{
		ID: r.ID, Title: r.Title,
		CreatedAt: r.CreatedAt.UTC().Format(time.RFC3339), UpdatedAt: r.UpdatedAt.UTC().Format(time.RFC3339),
		Score: toScore(r.Score),
	}
	if r.Latest != nil {
		lv := toVersion(r.Latest)
		resp.Latest = &lv
	}
	for i := range r.Versions {
		resp.Versions = append(resp.Versions, toVersion(&r.Versions[i]))
	}
	return resp
}

func toResumes(list []domain.Resume) []resumeResponse {
	out := make([]resumeResponse, 0, len(list))
	for i := range list {
		out = append(out, toResume(&list[i]))
	}
	return out
}
