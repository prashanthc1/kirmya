package application

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"workspace-app/internal/mentorship/domain"
)

type fakeRepo struct {
	seq      int
	mentors  map[string]*domain.MentorProfile // by id
	byUser   map[string]string                // userID -> mentorID
	sessions map[string]*domain.Session
	reviews  []domain.Review
	slots    map[string]*domain.AvailabilitySlot
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{mentors: map[string]*domain.MentorProfile{}, byUser: map[string]string{}, sessions: map[string]*domain.Session{}, slots: map[string]*domain.AvailabilitySlot{}}
}

func (r *fakeRepo) AddAvailability(_ context.Context, slot *domain.AvailabilitySlot) error {
	r.seq++
	slot.ID = fmt.Sprintf("slot-%d", r.seq)
	cp := *slot
	r.slots[slot.ID] = &cp
	return nil
}
func (r *fakeRepo) ListAvailability(_ context.Context, mentorID string, openOnly bool) ([]domain.AvailabilitySlot, error) {
	out := []domain.AvailabilitySlot{}
	for _, sl := range r.slots {
		if sl.MentorID == mentorID && (!openOnly || !sl.IsBooked) {
			out = append(out, *sl)
		}
	}
	return out, nil
}
func (r *fakeRepo) GetSlot(_ context.Context, id string) (*domain.AvailabilitySlot, error) {
	sl, ok := r.slots[id]
	if !ok {
		return nil, domain.ErrSlotNotFound
	}
	cp := *sl
	return &cp, nil
}

// CreateSessionWithSlot mirrors the real adapter: claim the slot only if it is
// still open, otherwise reject — both effects happen together.
func (r *fakeRepo) CreateSessionWithSlot(ctx context.Context, s *domain.Session, slotID string) error {
	sl, ok := r.slots[slotID]
	if !ok || sl.IsBooked {
		return domain.ErrSlotUnavailable
	}
	sl.IsBooked = true
	return r.CreateSession(ctx, s)
}

func (r *fakeRepo) UpsertMentorProfile(_ context.Context, p *domain.MentorProfile) error {
	if id, ok := r.byUser[p.UserID]; ok {
		p.ID = id
	} else {
		r.seq++
		p.ID = fmt.Sprintf("mentor-%d", r.seq)
		r.byUser[p.UserID] = p.ID
	}
	cp := *p
	r.mentors[p.ID] = &cp
	return nil
}
func (r *fakeRepo) GetMentorByID(_ context.Context, id string) (*domain.MentorProfile, error) {
	m, ok := r.mentors[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return m, nil
}
func (r *fakeRepo) GetMentorByUserID(_ context.Context, userID string) (*domain.MentorProfile, error) {
	id, ok := r.byUser[userID]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return r.mentors[id], nil
}
func (r *fakeRepo) ListMentors(_ context.Context) ([]domain.MentorProfile, error) {
	out := []domain.MentorProfile{}
	for _, m := range r.mentors {
		out = append(out, *m)
	}
	return out, nil
}
func (r *fakeRepo) CreateSession(_ context.Context, s *domain.Session) error {
	r.seq++
	s.ID = fmt.Sprintf("sess-%d", r.seq)
	cp := *s
	r.sessions[s.ID] = &cp
	return nil
}
func (r *fakeRepo) GetSession(_ context.Context, id string) (*domain.Session, error) {
	s, ok := r.sessions[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *s
	return &cp, nil
}
func (r *fakeRepo) ListSessionsForMentee(_ context.Context, menteeID string) ([]domain.Session, error) {
	out := []domain.Session{}
	for _, s := range r.sessions {
		if s.MenteeID == menteeID {
			out = append(out, *s)
		}
	}
	return out, nil
}
func (r *fakeRepo) ListSessionsForMentor(_ context.Context, mentorID string) ([]domain.Session, error) {
	out := []domain.Session{}
	for _, s := range r.sessions {
		if s.MentorID == mentorID {
			out = append(out, *s)
		}
	}
	return out, nil
}
func (r *fakeRepo) UpdateSessionStatus(_ context.Context, id, status string) error {
	if s, ok := r.sessions[id]; ok {
		s.Status = status
	}
	return nil
}
func (r *fakeRepo) CreateReview(_ context.Context, rv *domain.Review) error {
	r.seq++
	rv.ID = fmt.Sprintf("rev-%d", r.seq)
	r.reviews = append(r.reviews, *rv)
	return nil
}
func (r *fakeRepo) ListReviewsForMentor(_ context.Context, mentorID string) ([]domain.Review, error) {
	out := []domain.Review{}
	for _, rv := range r.reviews {
		if s, ok := r.sessions[rv.SessionID]; ok && s.MentorID == mentorID {
			out = append(out, rv)
		}
	}
	return out, nil
}

func setup(t *testing.T) (*Service, *fakeRepo, string) {
	t.Helper()
	repo := newFakeRepo()
	svc := NewService(repo, nil)
	m, err := svc.BecomeMentor(context.Background(), "mentor-user", MentorInput{Headline: "Ops leader"})
	if err != nil {
		t.Fatalf("become mentor: %v", err)
	}
	return svc, repo, m.ID
}

func TestBookRejectsSelf(t *testing.T) {
	svc, _, mentorID := setup(t)
	if _, err := svc.Book(context.Background(), "mentor-user", BookInput{MentorID: mentorID, Topic: "topic", ScheduledAt: time.Now().Add(time.Hour)}); err == nil {
		t.Fatal("expected error booking yourself")
	}
}

func TestStatusOnlyByMentor(t *testing.T) {
	svc, _, mentorID := setup(t)
	ctx := context.Background()
	sess, err := svc.Book(ctx, "mentee", BookInput{MentorID: mentorID, Topic: "interview prep", ScheduledAt: time.Now().Add(time.Hour)})
	if err != nil {
		t.Fatalf("book: %v", err)
	}
	if _, err := svc.UpdateStatus(ctx, "mentee", sess.ID, domain.StatusConfirmed); !errors.Is(err, domain.ErrNotMentor) {
		t.Fatalf("expected ErrNotMentor, got %v", err)
	}
	if _, err := svc.UpdateStatus(ctx, "mentor-user", sess.ID, domain.StatusConfirmed); err != nil {
		t.Fatalf("mentor update: %v", err)
	}
}

func TestReviewRequiresCompletedAndMentee(t *testing.T) {
	svc, _, mentorID := setup(t)
	ctx := context.Background()
	sess, _ := svc.Book(ctx, "mentee", BookInput{MentorID: mentorID, Topic: "topic", ScheduledAt: time.Now().Add(time.Hour)})

	if _, err := svc.Review(ctx, "mentee", sess.ID, 5, "great"); !errors.Is(err, domain.ErrNotComplete) {
		t.Fatalf("expected ErrNotComplete, got %v", err)
	}
	_, _ = svc.UpdateStatus(ctx, "mentor-user", sess.ID, domain.StatusCompleted)
	if _, err := svc.Review(ctx, "stranger", sess.ID, 5, "x"); !errors.Is(err, domain.ErrNotMentee) {
		t.Fatalf("expected ErrNotMentee, got %v", err)
	}
	if _, err := svc.Review(ctx, "mentee", sess.ID, 5, "great session"); err != nil {
		t.Fatalf("review: %v", err)
	}
}

func TestMentorReviewsAverage(t *testing.T) {
	svc, _, mentorID := setup(t)
	ctx := context.Background()

	for i, rating := range []int{4, 5} {
		sess, err := svc.Book(ctx, fmt.Sprintf("mentee-%d", i), BookInput{MentorID: mentorID, Topic: "topic", ScheduledAt: time.Now().Add(time.Hour)})
		if err != nil {
			t.Fatalf("book: %v", err)
		}
		if _, err := svc.UpdateStatus(ctx, "mentor-user", sess.ID, domain.StatusCompleted); err != nil {
			t.Fatalf("complete: %v", err)
		}
		if _, err := svc.Review(ctx, fmt.Sprintf("mentee-%d", i), sess.ID, rating, "ok"); err != nil {
			t.Fatalf("review: %v", err)
		}
	}

	stats, err := svc.MentorReviews(ctx, mentorID)
	if err != nil {
		t.Fatalf("mentor reviews: %v", err)
	}
	if stats.Count != 2 {
		t.Fatalf("expected 2 reviews, got %d", stats.Count)
	}
	if stats.AverageRating != 4.5 {
		t.Fatalf("expected average 4.5, got %v", stats.AverageRating)
	}
}

func TestMentorReviewsUnknownMentor(t *testing.T) {
	svc, _, _ := setup(t)
	if _, err := svc.MentorReviews(context.Background(), "nope"); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestAddAvailabilityValidation(t *testing.T) {
	svc, _, _ := setup(t)
	ctx := context.Background()
	start := time.Now().Add(time.Hour)
	// end before start
	if _, err := svc.AddAvailability(ctx, "mentor-user", start, start.Add(-time.Hour)); err == nil {
		t.Fatal("expected validation error when end precedes start")
	}
	// non-mentor caller
	if _, err := svc.AddAvailability(ctx, "stranger", start, start.Add(time.Hour)); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for non-mentor, got %v", err)
	}
	slot, err := svc.AddAvailability(ctx, "mentor-user", start, start.Add(time.Hour))
	if err != nil {
		t.Fatalf("add availability: %v", err)
	}
	if slot.ID == "" || slot.IsBooked {
		t.Fatalf("expected an open slot, got %+v", slot)
	}
}

func TestBookConsumesSlotAndRejectsDoubleBook(t *testing.T) {
	svc, _, mentorID := setup(t)
	ctx := context.Background()
	start := time.Now().Add(2 * time.Hour)
	slot, err := svc.AddAvailability(ctx, "mentor-user", start, start.Add(time.Hour))
	if err != nil {
		t.Fatalf("add availability: %v", err)
	}

	sess, err := svc.Book(ctx, "mentee", BookInput{MentorID: mentorID, Topic: "topic", SlotID: slot.ID})
	if err != nil {
		t.Fatalf("book against slot: %v", err)
	}
	if !sess.ScheduledAt.Equal(start) {
		t.Fatalf("expected scheduled time to follow the slot, got %v", sess.ScheduledAt)
	}
	// the slot is now booked; a second booking against it must fail
	if _, err := svc.Book(ctx, "mentee-2", BookInput{MentorID: mentorID, Topic: "topic", SlotID: slot.ID}); !errors.Is(err, domain.ErrSlotUnavailable) {
		t.Fatalf("expected ErrSlotUnavailable on double-book, got %v", err)
	}
	// the open availability list should now be empty
	open, err := svc.MentorAvailability(ctx, mentorID)
	if err != nil {
		t.Fatalf("availability: %v", err)
	}
	if len(open) != 0 {
		t.Fatalf("expected no open slots after booking, got %d", len(open))
	}
}

type fakePublisher struct{ events []string }

func (p *fakePublisher) Publish(_ context.Context, eventType, _ string, _ map[string]any) error {
	p.events = append(p.events, eventType)
	return nil
}

func TestEventsPublishedOnLifecycle(t *testing.T) {
	repo := newFakeRepo()
	pub := &fakePublisher{}
	svc := NewService(repo, pub)
	ctx := context.Background()
	mentor, err := svc.BecomeMentor(ctx, "mentor-user", MentorInput{Headline: "Ops leader"})
	if err != nil {
		t.Fatalf("become mentor: %v", err)
	}
	sess, err := svc.Book(ctx, "mentee", BookInput{MentorID: mentor.ID, Topic: "t", ScheduledAt: time.Now().Add(time.Hour)})
	if err != nil {
		t.Fatalf("book: %v", err)
	}
	if _, err := svc.UpdateStatus(ctx, "mentor-user", sess.ID, domain.StatusConfirmed); err != nil {
		t.Fatalf("confirm: %v", err)
	}
	if _, err := svc.UpdateStatus(ctx, "mentor-user", sess.ID, domain.StatusCompleted); err != nil {
		t.Fatalf("complete: %v", err)
	}
	if _, err := svc.Review(ctx, "mentee", sess.ID, 5, "great"); err != nil {
		t.Fatalf("review: %v", err)
	}
	want := []string{domain.EventMentorshipBooked, domain.EventSessionConfirmed, domain.EventSessionCompleted, domain.EventReviewLeft}
	if len(pub.events) != len(want) {
		t.Fatalf("expected events %v, got %v", want, pub.events)
	}
	for i, e := range want {
		if pub.events[i] != e {
			t.Fatalf("event %d: expected %s, got %s", i, e, pub.events[i])
		}
	}
}
