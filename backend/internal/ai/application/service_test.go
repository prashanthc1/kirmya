package application

import (
	"context"
	"errors"
	"testing"

	"workspace-app/internal/ai/domain"
)

func TestReviewResumeParsesFencedJSON(t *testing.T) {
	llm := &fakeLLM{ready: true, resp: "```json\n{\"summary\":\"Solid\",\"ats_score\":150,\"strengths\":[\"clear\"],\"improvements\":[\"add metrics\"]}\n```"}
	svc := NewService(newFakeRepo(), llm)

	review, err := svc.ReviewResume(context.Background(), "u1", "A reasonably long resume body with enough characters.")
	if err != nil {
		t.Fatalf("review: %v", err)
	}
	if review.Summary != "Solid" {
		t.Errorf("summary = %q", review.Summary)
	}
	if review.ATSScore != 100 {
		t.Errorf("ats score should clamp to 100, got %d", review.ATSScore)
	}
}

func TestReviewResumeShortRejected(t *testing.T) {
	svc := NewService(newFakeRepo(), &fakeLLM{ready: true, resp: "{}"})
	var ve ValidationError
	if _, err := svc.ReviewResume(context.Background(), "u1", "short"); !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
}

func TestAIDegradesWhenNotConfigured(t *testing.T) {
	svc := NewService(newFakeRepo(), &fakeLLM{ready: false})
	if svc.Available() {
		t.Fatal("expected Available() == false")
	}
	_, err := svc.ReviewResume(context.Background(), "u1", "A reasonably long resume body with enough characters.")
	if !errors.Is(err, domain.ErrLLMNotReady) {
		t.Fatalf("expected ErrLLMNotReady, got %v", err)
	}
}

func TestSkillGapFillsTargetRole(t *testing.T) {
	llm := &fakeLLM{ready: true, resp: `{"summary":"reachable","missing_skills":["PMP"],"learning_path":[{"skill":"PMP","resource":"course","why":"required"}]}`}
	svc := NewService(newFakeRepo(), llm)

	gap, err := svc.SkillGap(context.Background(), "u1", "Coordinator", "Project Manager", []string{"Scheduling"})
	if err != nil {
		t.Fatalf("skillgap: %v", err)
	}
	if gap.TargetRole != "Project Manager" {
		t.Errorf("target role should default to input, got %q", gap.TargetRole)
	}
	if len(gap.MissingSkills) != 1 || len(gap.LearningPath) != 1 {
		t.Errorf("unexpected gap: %+v", gap)
	}
}

func TestCoachCreatesThreadAndAppends(t *testing.T) {
	llm := &fakeLLM{ready: true, resp: "Here's some advice."}
	svc := NewService(newFakeRepo(), llm)
	ctx := context.Background()

	reply, threadID, err := svc.Coach(ctx, "owner", "", "I lost my job, where do I start?")
	if err != nil {
		t.Fatalf("coach: %v", err)
	}
	if reply != "Here's some advice." || threadID == "" {
		t.Fatalf("unexpected reply=%q thread=%q", reply, threadID)
	}

	// Second turn on the same thread should include prior history in the LLM call.
	if _, _, err := svc.Coach(ctx, "owner", threadID, "What about my resume?"); err != nil {
		t.Fatalf("coach 2: %v", err)
	}
	// History sent on turn 2 = [user1, assistant1] then current user appended = 3 messages.
	if len(llm.lastMsgs) != 3 {
		t.Fatalf("expected 3 messages sent to LLM on turn 2, got %d", len(llm.lastMsgs))
	}

	// Ownership enforced.
	if _, _, err := svc.Coach(ctx, "intruder", threadID, "hi"); !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}
}
