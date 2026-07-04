// Package aimatch implements jobs/application.JobMatcher using the AI module's
// LLM. When the LLM is unavailable or its response can't be parsed, it falls
// back to the injected heuristic matcher, so job matching always works.
package aimatch

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	aidomain "workspace-app/internal/ai/domain"
	"workspace-app/internal/ai/infrastructure/anthropic"
	"workspace-app/internal/jobs/application"
	"workspace-app/internal/jobs/domain"
)

const system = `You are a career-matching assistant for a job-seeker recovery platform.
Given a candidate's skills and a list of jobs, score how well each job fits the
candidate (0-100) considering their skills and the role.
Return ONLY JSON of the form:
{"matches":[{"id":"<job id>","score":<0-100>,"matched_skills":["..."],"missing_skills":["..."],"reason":"<one sentence>"}]}
Include every job id exactly once. missing_skills are skills the role needs that
the candidate lacks. Keep reasons to one concise sentence.`

// Matcher ranks jobs with the LLM, falling back to fallback on any problem.
type Matcher struct {
	llm      aidomain.LLM
	fallback application.JobMatcher
}

// New builds the AI matcher backed by Claude (anthropic). fallback is used
// whenever the LLM is not ready or returns an unparseable response.
func New(fallback application.JobMatcher) *Matcher {
	return &Matcher{llm: anthropic.New(), fallback: fallback}
}

type aiMatch struct {
	ID            string   `json:"id"`
	Score         int      `json:"score"`
	MatchedSkills []string `json:"matched_skills"`
	MissingSkills []string `json:"missing_skills"`
	Reason        string   `json:"reason"`
}

func (m *Matcher) Rank(ctx context.Context, skills []string, jobs []domain.Job) ([]domain.Match, error) {
	if !m.llm.Ready() {
		return m.fallback.Rank(ctx, skills, jobs)
	}

	prompt := buildPrompt(skills, jobs)
	c, err := m.llm.Complete(ctx, system, []aidomain.LLMMessage{{Role: aidomain.RoleUser, Content: prompt}}, 1500)
	if err != nil {
		return m.fallback.Rank(ctx, skills, jobs)
	}

	var parsed struct {
		Matches []aiMatch `json:"matches"`
	}
	if err := decodeJSON(c.Text, &parsed); err != nil || len(parsed.Matches) == 0 {
		return m.fallback.Rank(ctx, skills, jobs)
	}

	byID := make(map[string]aiMatch, len(parsed.Matches))
	for _, am := range parsed.Matches {
		byID[am.ID] = am
	}

	out := make([]domain.Match, 0, len(jobs))
	for _, j := range jobs {
		am, ok := byID[j.ID]
		if !ok {
			// The model skipped this job — fill it from the heuristic.
			out = append(out, application.ScoreJob(skills, j))
			continue
		}
		out = append(out, domain.Match{
			Job:           j,
			Score:         clamp(am.Score),
			MatchedSkills: nonNil(am.MatchedSkills),
			MissingSkills: nonNil(am.MissingSkills),
			Reason:        am.Reason,
		})
	}
	return out, nil
}

func buildPrompt(skills []string, jobs []domain.Job) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Candidate skills: %s\n\nJobs:\n", strings.Join(skills, ", "))
	for _, j := range jobs {
		desc := j.Description
		if len(desc) > 400 {
			desc = desc[:400]
		}
		fmt.Fprintf(&b, "- id=%s | %s at %s (%s) | %s\n", j.ID, j.Title, j.Company, j.Location, desc)
	}
	return b.String()
}

// decodeJSON tolerantly extracts a JSON object from an LLM response that may be
// wrapped in code fences or prose.
func decodeJSON(text string, dst any) error {
	t := strings.TrimSpace(text)
	t = strings.TrimPrefix(t, "```json")
	t = strings.TrimPrefix(t, "```")
	t = strings.TrimSuffix(t, "```")
	if i := strings.Index(t, "{"); i >= 0 {
		if j := strings.LastIndex(t, "}"); j > i {
			t = t[i : j+1]
		}
	}
	return json.Unmarshal([]byte(t), dst)
}

func clamp(v int) int {
	if v < 0 {
		return 0
	}
	if v > 100 {
		return 100
	}
	return v
}

func nonNil(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}
