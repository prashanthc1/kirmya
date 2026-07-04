package application

import (
	"context"
	"sort"
	"strings"

	"workspace-app/internal/jobs/domain"
)

// SkillReader reads a seeker's skills (from their profile) for matching.
type SkillReader interface {
	SeekerSkills(ctx context.Context, userID string) ([]string, error)
}

// JobMatcher ranks jobs for a seeker given their skills. Implementations: the
// pure-heuristic HeuristicMatcher (always available) and an AI-backed matcher
// that falls back to the heuristic.
type JobMatcher interface {
	Rank(ctx context.Context, skills []string, jobs []domain.Job) ([]domain.Match, error)
}

const matchCandidateLimit = 50

// MatchService produces ranked job recommendations for a seeker.
type MatchService struct {
	repo    domain.Repository
	skills  SkillReader
	matcher JobMatcher
}

func NewMatchService(repo domain.Repository, skills SkillReader, matcher JobMatcher) *MatchService {
	return &MatchService{repo: repo, skills: skills, matcher: matcher}
}

// Matches returns the seeker's job recommendations, best first.
func (s *MatchService) Matches(ctx context.Context, userID string) ([]domain.Match, error) {
	skills, err := s.skills.SeekerSkills(ctx, userID)
	if err != nil {
		return nil, err
	}
	jobs, err := s.repo.ListJobs(ctx, domain.Filter{Limit: matchCandidateLimit})
	if err != nil {
		return nil, err
	}
	if len(jobs) == 0 {
		return []domain.Match{}, nil
	}
	matches, err := s.matcher.Rank(ctx, skills, jobs)
	if err != nil {
		return nil, err
	}
	sort.SliceStable(matches, func(i, j int) bool { return matches[i].Score > matches[j].Score })
	return matches, nil
}

// HeuristicMatcher scores jobs by how many of the seeker's skills appear in the
// job's title/description. It needs no external services, so it is the default
// and the fallback for the AI matcher.
type HeuristicMatcher struct{}

func NewHeuristicMatcher() *HeuristicMatcher { return &HeuristicMatcher{} }

func (HeuristicMatcher) Rank(_ context.Context, skills []string, jobs []domain.Job) ([]domain.Match, error) {
	out := make([]domain.Match, 0, len(jobs))
	for _, j := range jobs {
		out = append(out, ScoreJob(skills, j))
	}
	return out, nil
}

// ScoreJob is the pure heuristic: skills found in the job text are "matched",
// and the score scales with coverage of the seeker's skills (with a small base
// so every listing is still rankable).
func ScoreJob(skills []string, j domain.Job) domain.Match {
	haystack := strings.ToLower(j.Title + " " + j.Company + " " + j.Description)
	matched := []string{}
	for _, sk := range skills {
		if s := strings.TrimSpace(sk); s != "" && strings.Contains(haystack, strings.ToLower(s)) {
			matched = append(matched, sk)
		}
	}

	score := 35 // base relevance for any open role
	if len(skills) > 0 {
		score = 20 + int(80.0*float64(len(matched))/float64(len(skills)))
	}
	if score > 100 {
		score = 100
	}

	reason := "Open role you can apply to."
	if len(matched) > 0 {
		reason = "Matches your skills: " + strings.Join(matched, ", ") + "."
	}
	return domain.Match{Job: j, Score: score, MatchedSkills: matched, MissingSkills: []string{}, Reason: reason}
}
