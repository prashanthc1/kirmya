package application

import (
	"context"
	"errors"
	"testing"

	"workspace-app/internal/ai/domain"
)

func TestSkillGapRequiresTargetRole(t *testing.T) {
	svc := NewService(newFakeRepo(), &fakeLLM{ready: true, resp: "{}"})
	var ve ValidationError
	if _, err := svc.SkillGap(context.Background(), "u1", "Coordinator", "  ", nil); !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError for missing target role, got %v", err)
	}
}

func TestCoachRequiresMessage(t *testing.T) {
	svc := NewService(newFakeRepo(), &fakeLLM{ready: true, resp: "ok"})
	var ve ValidationError
	if _, _, err := svc.Coach(context.Background(), "u1", "", "   "); !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError for empty message, got %v", err)
	}
}

func TestCoachPropagatesLLMError(t *testing.T) {
	boom := errors.New("llm exploded")
	svc := NewService(newFakeRepo(), &fakeLLM{ready: true, err: boom})
	if _, _, err := svc.Coach(context.Background(), "u1", "", "help me"); !errors.Is(err, boom) {
		t.Fatalf("expected LLM error to propagate, got %v", err)
	}
}

func TestReviewResumePropagatesLLMError(t *testing.T) {
	boom := errors.New("timeout")
	svc := NewService(newFakeRepo(), &fakeLLM{ready: true, err: boom})
	_, err := svc.ReviewResume(context.Background(), "u1", "A reasonably long resume body with enough characters.")
	if !errors.Is(err, boom) {
		t.Fatalf("expected LLM error to propagate, got %v", err)
	}
}

func TestReviewResumeUnparseableResponse(t *testing.T) {
	svc := NewService(newFakeRepo(), &fakeLLM{ready: true, resp: "I cannot help with that."})
	_, err := svc.ReviewResume(context.Background(), "u1", "A reasonably long resume body with enough characters.")
	if err == nil {
		t.Fatal("expected a parse error for a non-JSON response")
	}
}

func TestGetThreadOwnershipAndMessages(t *testing.T) {
	svc := NewService(newFakeRepo(), &fakeLLM{ready: true, resp: "advice"})
	ctx := context.Background()

	_, threadID, err := svc.Coach(ctx, "owner", "", "first message")
	if err != nil {
		t.Fatalf("coach: %v", err)
	}

	if _, err := svc.GetThread(ctx, "intruder", threadID); !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}
	thread, err := svc.GetThread(ctx, "owner", threadID)
	if err != nil {
		t.Fatalf("get thread: %v", err)
	}
	// One user turn + one assistant reply.
	if len(thread.Messages) != 2 {
		t.Fatalf("expected 2 messages in thread, got %d", len(thread.Messages))
	}
}

func TestListThreadsFiltersByUser(t *testing.T) {
	svc := NewService(newFakeRepo(), &fakeLLM{ready: true, resp: "advice"})
	ctx := context.Background()

	if _, _, err := svc.Coach(ctx, "owner", "", "mine one"); err != nil {
		t.Fatalf("coach: %v", err)
	}
	if _, _, err := svc.Coach(ctx, "owner", "", "mine two"); err != nil {
		t.Fatalf("coach: %v", err)
	}
	if _, _, err := svc.Coach(ctx, "other", "", "theirs"); err != nil {
		t.Fatalf("coach: %v", err)
	}

	threads, err := svc.ListThreads(ctx, "owner")
	if err != nil {
		t.Fatalf("list threads: %v", err)
	}
	if len(threads) != 2 {
		t.Fatalf("expected 2 threads for owner, got %d", len(threads))
	}
}

func TestGetThreadNotFound(t *testing.T) {
	svc := NewService(newFakeRepo(), &fakeLLM{ready: true, resp: "advice"})
	if _, err := svc.GetThread(context.Background(), "owner", "missing"); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
