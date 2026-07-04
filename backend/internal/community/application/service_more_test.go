package application

import (
	"context"
	"errors"
	"testing"

	"workspace-app/internal/community/domain"
)

// These tests exercise the read/passthrough use cases and the error branches
// that the original service_test.go does not, reusing its in-package fakeRepo.

func TestListAndGetCommunity(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo, nil)
	ctx := context.Background()

	comms, err := svc.ListCommunities(ctx)
	if err != nil {
		t.Fatalf("list communities: %v", err)
	}
	if len(comms) != 1 || comms[0].Slug != "technology" {
		t.Fatalf("expected the seeded technology community, got %+v", comms)
	}

	c, err := svc.GetCommunity(ctx, "technology")
	if err != nil {
		t.Fatalf("get community: %v", err)
	}
	if c.ID != "comm-1" {
		t.Fatalf("expected comm-1, got %q", c.ID)
	}

	if _, err := svc.GetCommunity(ctx, "nope"); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for unknown slug, got %v", err)
	}
}

func TestListPostsAndByTag(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo, nil)
	ctx := context.Background()

	if _, _, err := svc.ListPosts(ctx, "missing"); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound listing posts for unknown community, got %v", err)
	}

	if _, err := svc.CreatePost(ctx, "u1", "technology", "Tagged", "body", []string{"go", "career"}); err != nil {
		t.Fatalf("seed tagged post: %v", err)
	}
	if _, err := svc.CreatePost(ctx, "u1", "technology", "Untagged", "body", nil); err != nil {
		t.Fatalf("seed untagged post: %v", err)
	}

	c, posts, err := svc.ListPosts(ctx, "technology")
	if err != nil {
		t.Fatalf("list posts: %v", err)
	}
	if c.ID != "comm-1" || len(posts) != 2 {
		t.Fatalf("expected community + 2 posts, got %q / %d", c.ID, len(posts))
	}

	_, tagged, err := svc.PostsByTag(ctx, "technology", "go")
	if err != nil {
		t.Fatalf("posts by tag: %v", err)
	}
	if len(tagged) != 1 || tagged[0].Title != "Tagged" {
		t.Fatalf("expected only the tagged post, got %+v", tagged)
	}

	if _, _, err := svc.PostsByTag(ctx, "missing", "go"); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for unknown community, got %v", err)
	}
}

func TestTagsListing(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo, nil)
	ctx := context.Background()

	if _, err := svc.Tags(ctx, "missing"); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for unknown community, got %v", err)
	}

	if _, err := svc.CreatePost(ctx, "u1", "technology", "A", "", []string{"go"}); err != nil {
		t.Fatalf("seed post a: %v", err)
	}
	if _, err := svc.CreatePost(ctx, "u1", "technology", "B", "", []string{"go", "career"}); err != nil {
		t.Fatalf("seed post b: %v", err)
	}

	tags, err := svc.Tags(ctx, "technology")
	if err != nil {
		t.Fatalf("tags: %v", err)
	}
	counts := map[string]int{}
	for _, tg := range tags {
		counts[tg.Name] = tg.Count
	}
	if counts["go"] != 2 || counts["career"] != 1 {
		t.Fatalf("expected go=2 career=1, got %+v", counts)
	}
}

func TestAddCommentAndList(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo, nil)
	ctx := context.Background()
	p, _ := svc.CreatePost(ctx, "author", "technology", "Title", "", nil)

	if _, err := svc.AddComment(ctx, "u2", p.ID, "   "); err == nil {
		t.Fatal("expected validation error for blank comment body")
	}
	if _, err := svc.AddComment(ctx, "u2", "missing-post", "hi"); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound commenting on unknown post, got %v", err)
	}
	cmt, err := svc.AddComment(ctx, "u2", p.ID, "great post")
	if err != nil {
		t.Fatalf("add comment: %v", err)
	}
	if cmt.PostID != p.ID || cmt.Body != "great post" {
		t.Fatalf("unexpected comment %+v", cmt)
	}
	// Comments is a thin passthrough; just exercise it.
	if _, err := svc.Comments(ctx, p.ID); err != nil {
		t.Fatalf("comments: %v", err)
	}
}

func TestGetPoll(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo, nil)
	ctx := context.Background()
	p, _ := svc.CreatePost(ctx, "author", "technology", "Title", "", nil)
	poll, _ := svc.CreatePoll(ctx, "author", p.ID, "Pick", []string{"a", "b"})

	got, err := svc.GetPoll(ctx, poll.ID)
	if err != nil {
		t.Fatalf("get poll: %v", err)
	}
	if got.Question != "Pick" || len(got.Options) != 2 {
		t.Fatalf("unexpected poll %+v", got)
	}
	if _, err := svc.GetPoll(ctx, "missing"); !errors.Is(err, domain.ErrPollNotFound) {
		t.Fatalf("expected ErrPollNotFound, got %v", err)
	}
}

func TestReportUnknownPostAndHideGuards(t *testing.T) {
	repo := newFakeRepo()
	svc := NewService(repo, nil)
	ctx := context.Background()

	if _, err := svc.ReportPost(ctx, "reporter", "missing", "spam"); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound reporting unknown post, got %v", err)
	}

	// Hiding requires moderator; a non-moderator is rejected before anything else.
	p, _ := svc.CreatePost(ctx, "author", "technology", "Title", "", nil)
	if err := svc.HidePost(ctx, "rando", "technology", p.ID); !errors.Is(err, domain.ErrNotModerator) {
		t.Fatalf("expected ErrNotModerator, got %v", err)
	}

	// A moderator hiding a post that belongs to a different community gets ErrNotFound.
	repo.roles["comm-1|mod"] = domain.RoleModerator
	repo.posts["foreign"] = &domain.Post{ID: "foreign", CommunityID: "other-comm", Title: "Elsewhere"}
	if err := svc.HidePost(ctx, "mod", "technology", "foreign"); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound hiding a foreign-community post, got %v", err)
	}
}
