package domain

import "context"

// Repository is the persistence port for the Communities context. The concrete
// adapter lives in infrastructure/postgres.
type Repository interface {
	ListCommunities(ctx context.Context) ([]Community, error)
	GetBySlug(ctx context.Context, slug string) (*Community, error)
	Join(ctx context.Context, communityID, userID string) error
	Leave(ctx context.Context, communityID, userID string) error
	IsMember(ctx context.Context, communityID, userID string) (bool, error)
	MemberRole(ctx context.Context, communityID, userID string) (string, error)

	ListPosts(ctx context.Context, communityID string, limit int) ([]Post, error)
	ListPostsByTag(ctx context.Context, communityID, tag string, limit int) ([]Post, error)
	CreatePost(ctx context.Context, p *Post) error
	GetPost(ctx context.Context, id string) (*Post, error)
	DeletePost(ctx context.Context, id string) error

	AddComment(ctx context.Context, c *Comment) error
	ListComments(ctx context.Context, postID string) ([]Comment, error)

	ToggleReaction(ctx context.Context, postID, userID, kind string) (reacted bool, err error)

	SetPostTags(ctx context.Context, postID string, tags []string) error
	ListTags(ctx context.Context, communityID string) ([]Tag, error)

	CreatePoll(ctx context.Context, p *Poll) error
	GetPollByPost(ctx context.Context, postID string) (*Poll, error)
	GetPoll(ctx context.Context, pollID string) (*Poll, error)
	Vote(ctx context.Context, pollID, optionID, userID string) error

	CreateReport(ctx context.Context, rep *Report) error
	ListOpenReports(ctx context.Context, communityID string) ([]Report, error)
}
