package application

import (
	"context"
	"testing"

	"workspace-app/internal/jobs/domain"
)

type fakeSkills struct{ skills []string }

func (f fakeSkills) SeekerSkills(context.Context, string) ([]string, error) { return f.skills, nil }

func TestScoreJob(t *testing.T) {
	skills := []string{"Go", "PostgreSQL", "Leadership"}

	hit := ScoreJob(skills, domain.Job{Title: "Go Backend Engineer", Description: "Build services in Go with PostgreSQL."})
	if hit.Score <= 35 {
		t.Fatalf("expected boosted score for matches, got %d", hit.Score)
	}
	if len(hit.MatchedSkills) != 2 { // Go, PostgreSQL (not Leadership)
		t.Fatalf("expected 2 matched skills, got %v", hit.MatchedSkills)
	}

	miss := ScoreJob(skills, domain.Job{Title: "Pastry Chef", Description: "Bake croissants."})
	if len(miss.MatchedSkills) != 0 {
		t.Fatalf("expected no matches, got %v", miss.MatchedSkills)
	}
	if hit.Score <= miss.Score {
		t.Fatalf("a matching job (%d) should outrank a non-matching one (%d)", hit.Score, miss.Score)
	}
}

func TestScoreJobNoSkills(t *testing.T) {
	m := ScoreJob(nil, domain.Job{Title: "Anything"})
	if m.Score != 35 {
		t.Fatalf("expected base score 35 with no skills, got %d", m.Score)
	}
	if m.MatchedSkills == nil || m.MissingSkills == nil {
		t.Fatal("skill slices must be non-nil for stable JSON")
	}
}

func TestMatchServiceRanksBestFirst(t *testing.T) {
	repo := newFakeRepo()
	ctx := context.Background()
	// Two jobs: one matches the seeker's skills, one doesn't.
	_ = repo.CreateJob(ctx, &domain.Job{Title: "Operations Manager", Description: "Lead operations and logistics teams.", PostedBy: "r"})
	_ = repo.CreateJob(ctx, &domain.Job{Title: "Marine Biologist", Description: "Study coral reefs.", PostedBy: "r"})

	svc := NewMatchService(repo, fakeSkills{[]string{"Operations", "Logistics"}}, NewHeuristicMatcher())

	matches, err := svc.Matches(ctx, "seeker")
	if err != nil {
		t.Fatalf("matches: %v", err)
	}
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}
	if matches[0].Job.Title != "Operations Manager" {
		t.Fatalf("expected the operations role ranked first, got %q (%d) then %q (%d)",
			matches[0].Job.Title, matches[0].Score, matches[1].Job.Title, matches[1].Score)
	}
	if matches[0].Score < matches[1].Score {
		t.Fatal("results must be sorted by score descending")
	}
}

func TestMatchServiceEmptyWhenNoJobs(t *testing.T) {
	svc := NewMatchService(newFakeRepo(), fakeSkills{[]string{"Go"}}, NewHeuristicMatcher())
	matches, err := svc.Matches(context.Background(), "seeker")
	if err != nil || len(matches) != 0 {
		t.Fatalf("expected empty, got %v err=%v", matches, err)
	}
}
