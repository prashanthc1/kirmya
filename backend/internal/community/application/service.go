// Package application implements Communities use cases.
package application

import (
	"context"
	"strings"

	"workspace-app/internal/community/domain"
)

// EventPublisher is the best-effort port onto the platform event bus.
type EventPublisher interface {
	Publish(ctx context.Context, eventType, aggregateID string, payload map[string]any) error
}

type ValidationError struct{ Msg string }

func (e ValidationError) Error() string { return e.Msg }

type Service struct {
	repo   domain.Repository
	events EventPublisher
}

func NewService(repo domain.Repository, events EventPublisher) *Service {
	return &Service{repo: repo, events: events}
}

func (s *Service) publish(ctx context.Context, eventType, aggregateID string, payload map[string]any) {
	if s.events != nil {
		_ = s.events.Publish(ctx, eventType, aggregateID, payload)
	}
}

func (s *Service) ListCommunities(ctx context.Context) ([]domain.Community, error) {
	return s.repo.ListCommunities(ctx)
}

func (s *Service) GetCommunity(ctx context.Context, slug string) (*domain.Community, error) {
	return s.repo.GetBySlug(ctx, slug)
}

// ToggleJoin joins the community if not a member, otherwise leaves it.
func (s *Service) ToggleJoin(ctx context.Context, userID, slug string) (bool, error) {
	c, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return false, err
	}
	member, err := s.repo.IsMember(ctx, c.ID, userID)
	if err != nil {
		return false, err
	}
	if member {
		return false, s.repo.Leave(ctx, c.ID, userID)
	}
	return true, s.repo.Join(ctx, c.ID, userID)
}

func (s *Service) ListPosts(ctx context.Context, slug string) (*domain.Community, []domain.Post, error) {
	c, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, nil, err
	}
	posts, err := s.repo.ListPosts(ctx, c.ID, 100)
	if err != nil {
		return nil, nil, err
	}
	return c, posts, nil
}

// PostsByTag lists a community's posts filtered by a tag.
func (s *Service) PostsByTag(ctx context.Context, slug, tag string) (*domain.Community, []domain.Post, error) {
	c, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, nil, err
	}
	posts, err := s.repo.ListPostsByTag(ctx, c.ID, tag, 100)
	if err != nil {
		return nil, nil, err
	}
	return c, posts, nil
}

// CreatePost auto-joins the author so they become a member, then posts. Tags are
// normalized (trimmed, lowercased, de-duplicated) and attached to the post.
func (s *Service) CreatePost(ctx context.Context, userID, slug, title, body string, tags []string) (*domain.Post, error) {
	if strings.TrimSpace(title) == "" {
		return nil, ValidationError{"title is required"}
	}
	c, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	_ = s.repo.Join(ctx, c.ID, userID)
	clean := normalizeTags(tags)
	p := &domain.Post{CommunityID: c.ID, AuthorID: userID, Title: title, Body: body, Tags: clean}
	if err := s.repo.CreatePost(ctx, p); err != nil {
		return nil, err
	}
	if len(clean) > 0 {
		if err := s.repo.SetPostTags(ctx, p.ID, clean); err != nil {
			return nil, err
		}
	}
	s.publish(ctx, domain.EventPostCreated, p.ID, map[string]any{
		"post_id": p.ID, "community_id": c.ID, "author_id": userID,
	})
	return p, nil
}

func normalizeTags(tags []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, t := range tags {
		t = strings.ToLower(strings.TrimSpace(t))
		if t == "" || seen[t] {
			continue
		}
		seen[t] = true
		out = append(out, t)
	}
	return out
}

// Tags lists the tags in use across a community with their post counts.
func (s *Service) Tags(ctx context.Context, slug string) ([]domain.Tag, error) {
	c, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	return s.repo.ListTags(ctx, c.ID)
}

func (s *Service) AddComment(ctx context.Context, userID, postID, body string) (*domain.Comment, error) {
	if strings.TrimSpace(body) == "" {
		return nil, ValidationError{"comment body is required"}
	}
	if _, err := s.repo.GetPost(ctx, postID); err != nil {
		return nil, err
	}
	c := &domain.Comment{PostID: postID, AuthorID: userID, Body: body}
	if err := s.repo.AddComment(ctx, c); err != nil {
		return nil, err
	}
	s.publish(ctx, domain.EventCommentAdded, c.ID, map[string]any{"post_id": postID, "author_id": userID})
	return c, nil
}

func (s *Service) Comments(ctx context.Context, postID string) ([]domain.Comment, error) {
	return s.repo.ListComments(ctx, postID)
}

func (s *Service) ToggleReaction(ctx context.Context, userID, postID string) (bool, error) {
	if _, err := s.repo.GetPost(ctx, postID); err != nil {
		return false, err
	}
	return s.repo.ToggleReaction(ctx, postID, userID, "like")
}

// --- polls ---

// CreatePoll attaches a poll (question + 2-6 options) to an existing post.
// Only the post's author may add its poll, and a post may have at most one.
func (s *Service) CreatePoll(ctx context.Context, userID, postID, question string, options []string) (*domain.Poll, error) {
	if strings.TrimSpace(question) == "" {
		return nil, domain.ErrInvalidPoll
	}
	opts := normalizeOptions(options)
	if len(opts) < 2 || len(opts) > 6 {
		return nil, domain.ErrInvalidPoll
	}
	post, err := s.repo.GetPost(ctx, postID)
	if err != nil {
		return nil, err
	}
	if post.AuthorID != userID {
		return nil, ValidationError{"only the post author can add a poll"}
	}
	if existing, err := s.repo.GetPollByPost(ctx, postID); err == nil && existing != nil {
		return nil, domain.ErrPollExists
	}
	poll := &domain.Poll{PostID: postID, Question: question}
	for _, label := range opts {
		poll.Options = append(poll.Options, domain.PollOption{Label: label})
	}
	if err := s.repo.CreatePoll(ctx, poll); err != nil {
		return nil, err
	}
	return poll, nil
}

func normalizeOptions(options []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, o := range options {
		o = strings.TrimSpace(o)
		if o == "" || seen[strings.ToLower(o)] {
			continue
		}
		seen[strings.ToLower(o)] = true
		out = append(out, o)
	}
	return out
}

// GetPoll returns a poll with its options and current vote counts.
func (s *Service) GetPoll(ctx context.Context, pollID string) (*domain.Poll, error) {
	return s.repo.GetPoll(ctx, pollID)
}

// Vote records the caller's (single) vote for an option in a poll. Voting again
// moves the vote to the new option.
func (s *Service) Vote(ctx context.Context, userID, pollID, optionID string) (*domain.Poll, error) {
	poll, err := s.repo.GetPoll(ctx, pollID)
	if err != nil {
		return nil, err
	}
	valid := false
	for _, o := range poll.Options {
		if o.ID == optionID {
			valid = true
			break
		}
	}
	if !valid {
		return nil, domain.ErrOptionNotInPoll
	}
	if err := s.repo.Vote(ctx, pollID, optionID, userID); err != nil {
		return nil, err
	}
	s.publish(ctx, domain.EventPollVoted, pollID, map[string]any{"poll_id": pollID, "user_id": userID})
	return s.repo.GetPoll(ctx, pollID)
}

// --- moderation / reporting ---

// ReportPost files a moderation report against a post. Any authenticated user may report.
func (s *Service) ReportPost(ctx context.Context, userID, postID, reason string) (*domain.Report, error) {
	if strings.TrimSpace(reason) == "" {
		return nil, ValidationError{"a reason is required"}
	}
	if _, err := s.repo.GetPost(ctx, postID); err != nil {
		return nil, err
	}
	rep := &domain.Report{ReporterID: userID, TargetType: "post", TargetID: postID, Reason: reason, Status: "open"}
	if err := s.repo.CreateReport(ctx, rep); err != nil {
		return nil, err
	}
	s.publish(ctx, domain.EventPostReported, postID, map[string]any{"post_id": postID, "reporter_id": userID})
	return rep, nil
}

// requireModerator resolves a community by slug and verifies the caller moderates it.
func (s *Service) requireModerator(ctx context.Context, userID, slug string) (*domain.Community, error) {
	c, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	role, err := s.repo.MemberRole(ctx, c.ID, userID)
	if err != nil {
		return nil, err
	}
	if role != domain.RoleModerator {
		return nil, domain.ErrNotModerator
	}
	return c, nil
}

// CommunityReports lists open reports for a community. Caller must be a moderator.
func (s *Service) CommunityReports(ctx context.Context, userID, slug string) ([]domain.Report, error) {
	c, err := s.requireModerator(ctx, userID, slug)
	if err != nil {
		return nil, err
	}
	return s.repo.ListOpenReports(ctx, c.ID)
}

// HidePost removes a post. Caller must be a moderator of the post's community.
func (s *Service) HidePost(ctx context.Context, userID, slug, postID string) error {
	c, err := s.requireModerator(ctx, userID, slug)
	if err != nil {
		return err
	}
	post, err := s.repo.GetPost(ctx, postID)
	if err != nil {
		return err
	}
	if post.CommunityID != c.ID {
		return domain.ErrNotFound
	}
	if err := s.repo.DeletePost(ctx, postID); err != nil {
		return err
	}
	s.publish(ctx, domain.EventPostHidden, postID, map[string]any{"post_id": postID, "moderator_id": userID})
	return nil
}

// CreateCommunity creates a new community with the caller as moderator.
func (s *Service) CreateCommunity(ctx context.Context, creatorUserID, name, slug, description, category string) (*domain.Community, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ValidationError{"name is required"}
	}
	if strings.TrimSpace(slug) == "" {
		return nil, ValidationError{"slug is required"}
	}
	slug = strings.ToLower(strings.TrimSpace(slug))
	slug = strings.ReplaceAll(slug, " ", "-")

	c := &domain.Community{
		Slug:        slug,
		Name:        name,
		Description: description,
		Category:    category,
	}

	if err := s.repo.CreateCommunity(ctx, c, creatorUserID); err != nil {
		return nil, err
	}

	return c, nil
}

