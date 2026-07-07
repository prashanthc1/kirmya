// Package domain holds the AI bounded context: career-coach threads, the LLM
// port, and the result types for resume review and skill-gap analysis.
package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound    = errors.New("not found")
	ErrLLMNotReady = errors.New("AI provider is not configured")
)

const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
)

// CoachThread is a career-coaching conversation.
type CoachThread struct {
	ID        string
	UserID    string
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
	Messages  []CoachMessage
}

type CoachMessage struct {
	ID        string
	ThreadID  string
	Role      string
	Content   string
	CreatedAt time.Time
}

// ResumeReview is the AI resume-reviewer output.
type ResumeReview struct {
	Summary            string   `json:"summary"`
	ATSScore           int      `json:"ats_score"`
	KeywordFeedback    string   `json:"keyword_feedback"`
	FormattingFeedback string   `json:"formatting_feedback"`
	Strengths          []string `json:"strengths"`
	Improvements       []string `json:"improvements"`
}

// SkillGap is the AI skill-gap engine output.
type SkillGap struct {
	TargetRole     string         `json:"target_role"`
	Summary        string         `json:"summary"`
	MissingSkills  []string       `json:"missing_skills"`
	SuggestedRoles []string       `json:"suggested_roles"`
	LearningPath   []LearningStep `json:"learning_path"`
}

type LearningStep struct {
	Skill    string `json:"skill"`
	Resource string `json:"resource"`
	Why      string `json:"why"`
}

// LLMMessage is one turn in a completion request.
type LLMMessage struct {
	Role    string
	Content string
}

// Completion is the result of an LLM call.
type Completion struct {
	Text         string
	Model        string
	InputTokens  int
	OutputTokens int
}

// LLM is the provider port (Claude primary; OpenAI/others pluggable later).
type LLM interface {
	// Ready reports whether the provider is configured (API key present).
	Ready() bool
	Complete(ctx context.Context, system string, messages []LLMMessage, maxTokens int) (Completion, error)
	StreamComplete(ctx context.Context, system string, messages []LLMMessage, maxTokens int) (chan string, error)
}

// Repository persists coach threads/messages and the interaction log.
type Repository interface {
	CreateThread(ctx context.Context, t *CoachThread) error
	GetThread(ctx context.Context, id string) (*CoachThread, error)
	ListThreads(ctx context.Context, userID string) ([]CoachThread, error)
	AddMessage(ctx context.Context, m *CoachMessage) error
	ListMessages(ctx context.Context, threadID string) ([]CoachMessage, error)
	LogInteraction(ctx context.Context, userID, kind string, c Completion) error
}
