package application

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"workspace-app/internal/community/domain"
)

type fakeRepo struct {
	seq        int
	bySlug     map[string]*domain.Community
	members    map[string]bool   // communityID|userID
	roles      map[string]string // communityID|userID -> role
	posts      map[string]*domain.Post
	reactions  map[string]bool // postID|userID
	tags       map[string][]string
	polls      map[string]*domain.Poll      // by poll id
	pollByPost map[string]string            // postID -> pollID
	votes      map[string]map[string]string // pollID -> userID -> optionID
	reports    []domain.Report
}

func newFakeRepo() *fakeRepo {
	r := &fakeRepo{
		bySlug: map[string]*domain.Community{}, members: map[string]bool{}, roles: map[string]string{},
		posts: map[string]*domain.Post{}, reactions: map[string]bool{}, tags: map[string][]string{},
		polls: map[string]*domain.Poll{}, pollByPost: map[string]string{}, votes: map[string]map[string]string{},
	}
	r.bySlug["technology"] = &domain.Community{ID: "comm-1", Slug: "technology", Name: "Technology"}
	return r
}

func (r *fakeRepo) ListCommunities(_ context.Context) ([]domain.Community, error) {
	out := []domain.Community{}
	for _, c := range r.bySlug {
		out = append(out, *c)
	}
	return out, nil
}
func (r *fakeRepo) GetBySlug(_ context.Context, slug string) (*domain.Community, error) {
	c, ok := r.bySlug[slug]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return c, nil
}
func (r *fakeRepo) Join(_ context.Context, cid, uid string) error {
	r.members[cid+"|"+uid] = true
	return nil
}
func (r *fakeRepo) Leave(_ context.Context, cid, uid string) error {
	delete(r.members, cid+"|"+uid)
	return nil
}
func (r *fakeRepo) IsMember(_ context.Context, cid, uid string) (bool, error) {
	return r.members[cid+"|"+uid], nil
}
func (r *fakeRepo) ListPosts(_ context.Context, cid string, _ int) ([]domain.Post, error) {
	out := []domain.Post{}
	for _, p := range r.posts {
		if p.CommunityID == cid {
			out = append(out, *p)
		}
	}
	return out, nil
}
func (r *fakeRepo) CreatePost(_ context.Context, p *domain.Post) error {
	r.seq++
	p.ID = fmt.Sprintf("post-%d", r.seq)
	cp := *p
	r.posts[p.ID] = &cp
	return nil
}
func (r *fakeRepo) GetPost(_ context.Context, id string) (*domain.Post, error) {
	p, ok := r.posts[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return p, nil
}
func (r *fakeRepo) AddComment(_ context.Context, c *domain.Comment) error { c.ID = "cmt"; return nil }
func (r *fakeRepo) ListComments(_ context.Context, _ string) ([]domain.Comment, error) {
	return nil, nil
}
func (r *fakeRepo) ToggleReaction(_ context.Context, postID, userID, _ string) (bool, error) {
	k := postID + "|" + userID
	if r.reactions[k] {
		delete(r.reactions, k)
		return false, nil
	}
	r.reactions[k] = true
	return true, nil
}
func (r *fakeRepo) MemberRole(_ context.Context, cid, uid string) (string, error) {
	return r.roles[cid+"|"+uid], nil
}
func (r *fakeRepo) DeletePost(_ context.Context, id string) error { delete(r.posts, id); return nil }
func (r *fakeRepo) SetPostTags(_ context.Context, postID string, tags []string) error {
	r.tags[postID] = tags
	return nil
}
func (r *fakeRepo) ListTags(_ context.Context, cid string) ([]domain.Tag, error) {
	counts := map[string]int{}
	for postID, tags := range r.tags {
		if p, ok := r.posts[postID]; ok && p.CommunityID == cid {
			for _, t := range tags {
				counts[t]++
			}
		}
	}
	out := []domain.Tag{}
	for name, n := range counts {
		out = append(out, domain.Tag{Name: name, Count: n})
	}
	return out, nil
}
func (r *fakeRepo) ListPostsByTag(_ context.Context, cid, tag string, _ int) ([]domain.Post, error) {
	out := []domain.Post{}
	for postID, tags := range r.tags {
		p, ok := r.posts[postID]
		if !ok || p.CommunityID != cid {
			continue
		}
		for _, t := range tags {
			if t == tag {
				out = append(out, *p)
				break
			}
		}
	}
	return out, nil
}
func (r *fakeRepo) CreatePoll(_ context.Context, p *domain.Poll) error {
	r.seq++
	p.ID = fmt.Sprintf("poll-%d", r.seq)
	for i := range p.Options {
		p.Options[i].ID = fmt.Sprintf("%s-opt-%d", p.ID, i)
		p.Options[i].PollID = p.ID
	}
	cp := *p
	r.polls[p.ID] = &cp
	r.pollByPost[p.PostID] = p.ID
	return nil
}
func (r *fakeRepo) GetPollByPost(_ context.Context, postID string) (*domain.Poll, error) {
	id, ok := r.pollByPost[postID]
	if !ok {
		return nil, domain.ErrPollNotFound
	}
	return r.polls[id], nil
}
func (r *fakeRepo) GetPoll(_ context.Context, pollID string) (*domain.Poll, error) {
	p, ok := r.polls[pollID]
	if !ok {
		return nil, domain.ErrPollNotFound
	}
	cp := *p
	cp.Options = append([]domain.PollOption(nil), p.Options...)
	for i := range cp.Options {
		n := 0
		for _, opt := range r.votes[pollID] {
			if opt == cp.Options[i].ID {
				n++
			}
		}
		cp.Options[i].VoteCount = n
	}
	return &cp, nil
}
func (r *fakeRepo) Vote(_ context.Context, pollID, optionID, userID string) error {
	if r.votes[pollID] == nil {
		r.votes[pollID] = map[string]string{}
	}
	r.votes[pollID][userID] = optionID
	return nil
}
func (r *fakeRepo) CreateReport(_ context.Context, rep *domain.Report) error {
	r.seq++
	rep.ID = fmt.Sprintf("report-%d", r.seq)
	rep.Status = "open"
	r.reports = append(r.reports, *rep)
	return nil
}
func (r *fakeRepo) ListOpenReports(_ context.Context, cid string) ([]domain.Report, error) {
	out := []domain.Report{}
	for _, rep := range r.reports {
		if p, ok := r.posts[rep.TargetID]; ok && p.CommunityID == cid && rep.Status == "open" {
			out = append(out, rep)
		}
	}
	return out, nil
}

func TestToggleJoin(t *testing.T) {
	svc := NewService(newFakeRepo(), nil)
	ctx := context.Background()
	joined, err := svc.ToggleJoin(ctx, "u1", "technology")
	if err != nil || !joined {
		t.Fatalf("expected joined=true, got %v err=%v", joined, err)
	}
	joined, _ = svc.ToggleJoin(ctx, "u1", "technology")
	if joined {
		t.Fatal("expected joined=false after second toggle")
	}
}

func TestCreatePostAutoJoinsAndValidates(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo, nil)
	ctx := context.Background()

	if _, err := svc.CreatePost(ctx, "u1", "technology", "", "body", nil); err == nil {
		t.Fatal("expected validation error for empty title")
	}
	p, err := svc.CreatePost(ctx, "u1", "technology", "Re-skilling tips", "body", []string{"Go", "go", " career "})
	if err != nil {
		t.Fatalf("create post: %v", err)
	}
	if len(p.Tags) != 2 {
		t.Fatalf("expected normalized tags [go career], got %v", p.Tags)
	}
	if p.ID == "" {
		t.Fatal("expected post id")
	}
	if member, _ := repo.IsMember(ctx, "comm-1", "u1"); !member {
		t.Fatal("expected author auto-joined the community")
	}
}

func TestToggleReaction(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo, nil)
	ctx := context.Background()
	p, _ := svc.CreatePost(ctx, "u1", "technology", "Title", "", nil)

	reacted, err := svc.ToggleReaction(ctx, "u2", p.ID)
	if err != nil || !reacted {
		t.Fatalf("expected reacted=true, got %v err=%v", reacted, err)
	}
	reacted, _ = svc.ToggleReaction(ctx, "u2", p.ID)
	if reacted {
		t.Fatal("expected reacted=false after toggle")
	}
}

func TestCreatePollValidationAndOwnership(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo, nil)
	ctx := context.Background()
	p, _ := svc.CreatePost(ctx, "author", "technology", "Title", "", nil)

	if _, err := svc.CreatePoll(ctx, "author", p.ID, "Pick one", []string{"only-one"}); !errors.Is(err, domain.ErrInvalidPoll) {
		t.Fatalf("expected ErrInvalidPoll for <2 options, got %v", err)
	}
	if _, err := svc.CreatePoll(ctx, "intruder", p.ID, "Pick one", []string{"a", "b"}); err == nil {
		t.Fatal("expected non-author to be rejected")
	}
	poll, err := svc.CreatePoll(ctx, "author", p.ID, "Pick one", []string{"a", "b", "a"})
	if err != nil {
		t.Fatalf("create poll: %v", err)
	}
	if len(poll.Options) != 2 {
		t.Fatalf("expected duplicate option removed, got %d", len(poll.Options))
	}
	if _, err := svc.CreatePoll(ctx, "author", p.ID, "Again", []string{"a", "b"}); !errors.Is(err, domain.ErrPollExists) {
		t.Fatalf("expected ErrPollExists, got %v", err)
	}
}

func TestVoteMovesSingleVote(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo, nil)
	ctx := context.Background()
	p, _ := svc.CreatePost(ctx, "author", "technology", "Title", "", nil)
	poll, _ := svc.CreatePoll(ctx, "author", p.ID, "Pick one", []string{"a", "b"})

	if _, err := svc.Vote(ctx, "voter", poll.ID, "missing-option"); !errors.Is(err, domain.ErrOptionNotInPoll) {
		t.Fatalf("expected ErrOptionNotInPoll, got %v", err)
	}
	after, err := svc.Vote(ctx, "voter", poll.ID, poll.Options[0].ID)
	if err != nil {
		t.Fatalf("vote: %v", err)
	}
	if total(after) != 1 {
		t.Fatalf("expected 1 vote, got %d", total(after))
	}
	// same voter moves their vote — still one total vote
	after, _ = svc.Vote(ctx, "voter", poll.ID, poll.Options[1].ID)
	if total(after) != 1 {
		t.Fatalf("expected vote to move (still 1 total), got %d", total(after))
	}
	if after.Options[1].VoteCount != 1 {
		t.Fatalf("expected the moved vote on option 2, got %+v", after.Options)
	}
}

func total(p *domain.Poll) int {
	n := 0
	for _, o := range p.Options {
		n += o.VoteCount
	}
	return n
}

func TestReportAndModeration(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo, nil)
	ctx := context.Background()
	p, _ := svc.CreatePost(ctx, "author", "technology", "Spam", "", nil)

	if _, err := svc.ReportPost(ctx, "reporter", p.ID, ""); err == nil {
		t.Fatal("expected validation error for empty reason")
	}
	if _, err := svc.ReportPost(ctx, "reporter", p.ID, "spam"); err != nil {
		t.Fatalf("report: %v", err)
	}

	// non-moderator cannot list reports or hide
	if _, err := svc.CommunityReports(ctx, "reporter", "technology"); !errors.Is(err, domain.ErrNotModerator) {
		t.Fatalf("expected ErrNotModerator, got %v", err)
	}
	repo.roles["comm-1|mod"] = domain.RoleModerator
	reports, err := svc.CommunityReports(ctx, "mod", "technology")
	if err != nil {
		t.Fatalf("list reports: %v", err)
	}
	if len(reports) != 1 {
		t.Fatalf("expected 1 open report, got %d", len(reports))
	}
	if err := svc.HidePost(ctx, "mod", "technology", p.ID); err != nil {
		t.Fatalf("hide post: %v", err)
	}
	if _, err := repo.GetPost(ctx, p.ID); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected post removed, got %v", err)
	}
}
