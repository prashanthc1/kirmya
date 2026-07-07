package application

import (
	"context"
	"fmt"
	"sync"

	"workspace-app/internal/ai/domain"
)

// fakeLLM returns a canned completion and records the last messages it saw.
type fakeLLM struct {
	ready    bool
	resp     string
	err      error
	lastMsgs []domain.LLMMessage
}

func (f *fakeLLM) Ready() bool { return f.ready }

func (f *fakeLLM) Complete(_ context.Context, _ string, msgs []domain.LLMMessage, _ int) (domain.Completion, error) {
	f.lastMsgs = msgs
	if f.err != nil {
		return domain.Completion{}, f.err
	}
	return domain.Completion{Text: f.resp, Model: "test-model", InputTokens: 10, OutputTokens: 20}, nil
}

func (f *fakeLLM) StreamComplete(_ context.Context, _ string, msgs []domain.LLMMessage, _ int) (chan string, error) {
	f.lastMsgs = msgs
	if f.err != nil {
		return nil, f.err
	}
	out := make(chan string, 1)
	out <- f.resp
	close(out)
	return out, nil
}

// fakeRepo is an in-memory ai/domain.Repository.
type fakeRepo struct {
	mu       sync.Mutex
	seq      int
	threads  map[string]*domain.CoachThread
	messages map[string][]domain.CoachMessage
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{threads: map[string]*domain.CoachThread{}, messages: map[string][]domain.CoachMessage{}}
}

func (r *fakeRepo) id(p string) string { r.seq++; return fmt.Sprintf("%s-%d", p, r.seq) }

func (r *fakeRepo) CreateThread(_ context.Context, t *domain.CoachThread) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	t.ID = r.id("thread")
	c := *t
	r.threads[t.ID] = &c
	return nil
}

func (r *fakeRepo) GetThread(_ context.Context, id string) (*domain.CoachThread, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.threads[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	c := *t
	return &c, nil
}

func (r *fakeRepo) ListThreads(_ context.Context, userID string) ([]domain.CoachThread, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []domain.CoachThread{}
	for _, t := range r.threads {
		if t.UserID == userID {
			out = append(out, *t)
		}
	}
	return out, nil
}

func (r *fakeRepo) AddMessage(_ context.Context, m *domain.CoachMessage) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	m.ID = r.id("msg")
	r.messages[m.ThreadID] = append(r.messages[m.ThreadID], *m)
	return nil
}

func (r *fakeRepo) ListMessages(_ context.Context, threadID string) ([]domain.CoachMessage, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]domain.CoachMessage, len(r.messages[threadID]))
	copy(out, r.messages[threadID])
	return out, nil
}

func (r *fakeRepo) LogInteraction(_ context.Context, _, _ string, _ domain.Completion) error {
	return nil
}
