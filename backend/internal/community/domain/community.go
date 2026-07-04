// Package domain holds the Communities bounded context.
package domain

import (
	"errors"
	"time"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrNotModerator    = errors.New("you are not a moderator of this community")
	ErrPollExists      = errors.New("this post already has a poll")
	ErrInvalidPoll     = errors.New("a poll needs a question and between 2 and 6 options")
	ErrPollNotFound    = errors.New("poll not found")
	ErrOptionNotInPoll = errors.New("option does not belong to this poll")
)

// Moderation roles stored in community_members.role.
const (
	RoleMember    = "member"
	RoleModerator = "moderator"
)

type Community struct {
	ID          string
	Slug        string
	Name        string
	Description string
	Category    string
	MemberCount int
	CreatedAt   time.Time
}

type Post struct {
	ID            string
	CommunityID   string
	AuthorID      string
	Title         string
	Body          string
	CommentCount  int
	ReactionCount int
	Tags          []string
	CreatedAt     time.Time
}

type Comment struct {
	ID        string
	PostID    string
	AuthorID  string
	Body      string
	CreatedAt time.Time
}

// Poll is an optional vote attached to a post; a user has at most one vote per poll.
type Poll struct {
	ID        string
	PostID    string
	Question  string
	Options   []PollOption
	CreatedAt time.Time
}

type PollOption struct {
	ID        string
	PollID    string
	Label     string
	VoteCount int
}

// Report is a moderation report filed against a post or comment. It is persisted
// in the shared content_reports table (target_type in 'post'|'comment').
type Report struct {
	ID         string
	ReporterID string
	TargetType string
	TargetID   string
	Reason     string
	Status     string
	CreatedAt  time.Time
}

// Tag is a post tag with how many posts in the community carry it.
type Tag struct {
	Name  string
	Count int
}
