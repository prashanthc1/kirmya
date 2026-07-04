// Package application implements the AI use cases: resume review, career coach
// (threaded), and skill-gap analysis — orchestrating the LLM port + repository.
package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"workspace-app/internal/ai/domain"
)

var (
	ErrForbidden = errors.New("forbidden")
)

type ValidationError struct{ Msg string }

func (e ValidationError) Error() string { return e.Msg }

type Service struct {
	repo domain.Repository
	llm  domain.LLM
}

func NewService(repo domain.Repository, llm domain.LLM) *Service {
	return &Service{repo: repo, llm: llm}
}

// Available reports whether AI features are usable (provider configured).
func (s *Service) Available() bool { return s.llm != nil && s.llm.Ready() }

func (s *Service) ensureReady() error {
	if !s.Available() {
		return domain.ErrLLMNotReady
	}
	return nil
}

// ---------- Resume Review ----------

const resumeReviewSystem = `You are an expert technical recruiter and ATS (applicant tracking system) specialist.
Analyze the candidate's resume and return strict JSON only — no markdown, no commentary.
JSON shape:
{
  "summary": "2-3 sentence overall assessment",
  "ats_score": <integer 0-100>,
  "keyword_feedback": "how well it matches common ATS keywords and what to add",
  "formatting_feedback": "structure, length, readability feedback",
  "strengths": ["..."],
  "improvements": ["actionable suggestion", "..."]
}`

func (s *Service) ReviewResume(ctx context.Context, userID, resumeText string) (*domain.ResumeReview, error) {
	if err := s.ensureReady(); err != nil {
		return nil, err
	}
	if len(strings.TrimSpace(resumeText)) < 30 {
		return nil, ValidationError{"resume text is too short to review"}
	}
	c, err := s.llm.Complete(ctx, resumeReviewSystem, []domain.LLMMessage{
		{Role: domain.RoleUser, Content: "Resume:\n\n" + resumeText},
	}, 2000)
	if err != nil {
		return nil, err
	}
	_ = s.repo.LogInteraction(ctx, userID, "resume_review", c)

	var review domain.ResumeReview
	if err := decodeJSON(c.Text, &review); err != nil {
		return nil, fmt.Errorf("could not parse AI response: %w", err)
	}
	review.ATSScore = clamp(review.ATSScore)
	return &review, nil
}

// ---------- Skill Gap ----------

const skillGapSystem = `You are a career coach and labor-market analyst helping someone change roles or recover from job loss.
Given a current role, a target role, and current skills, identify the gap and a concrete learning path.
Return strict JSON only — no markdown, no commentary.
JSON shape:
{
  "target_role": "...",
  "summary": "2-3 sentence assessment of the transition",
  "missing_skills": ["..."],
  "suggested_roles": ["adjacent or stepping-stone roles"],
  "learning_path": [{"skill": "...", "resource": "concrete course/cert/practice", "why": "..."}]
}`

func (s *Service) SkillGap(ctx context.Context, userID, currentRole, targetRole string, currentSkills []string) (*domain.SkillGap, error) {
	if err := s.ensureReady(); err != nil {
		return nil, err
	}
	if strings.TrimSpace(targetRole) == "" {
		return nil, ValidationError{"target role is required"}
	}
	prompt := fmt.Sprintf("Current role: %s\nTarget role: %s\nCurrent skills: %s",
		fallback(currentRole, "(unspecified)"), targetRole, strings.Join(currentSkills, ", "))
	c, err := s.llm.Complete(ctx, skillGapSystem, []domain.LLMMessage{
		{Role: domain.RoleUser, Content: prompt},
	}, 2000)
	if err != nil {
		return nil, err
	}
	_ = s.repo.LogInteraction(ctx, userID, "skill_gap", c)

	var gap domain.SkillGap
	if err := decodeJSON(c.Text, &gap); err != nil {
		return nil, fmt.Errorf("could not parse AI response: %w", err)
	}
	if gap.TargetRole == "" {
		gap.TargetRole = targetRole
	}
	return &gap, nil
}

// ---------- Career Coach (threaded) ----------

const coachSystem = `You are Kirmya's AI career coach — a mentor, not a recruiter. Your user may have recently lost their job or be navigating a difficult transition.
Acknowledge the setback and help them celebrate the comeback. Be warm, practical, and encouraging without being saccharine. Give specific, actionable guidance on job search, interview prep, skill-building, and career planning. Keep replies focused and concise.`

// Coach appends a user message to a thread (creating one if threadID is empty)
// and returns the assistant reply plus the thread id.
func (s *Service) Coach(ctx context.Context, userID, threadID, message string) (string, string, error) {
	if err := s.ensureReady(); err != nil {
		return "", "", err
	}
	if strings.TrimSpace(message) == "" {
		return "", "", ValidationError{"message is required"}
	}

	var history []domain.CoachMessage
	if threadID == "" {
		t := &domain.CoachThread{UserID: userID, Title: makeTitle(message)}
		if err := s.repo.CreateThread(ctx, t); err != nil {
			return "", "", err
		}
		threadID = t.ID
	} else {
		t, err := s.repo.GetThread(ctx, threadID)
		if err != nil {
			return "", "", err
		}
		if t.UserID != userID {
			return "", "", ErrForbidden
		}
		if history, err = s.repo.ListMessages(ctx, threadID); err != nil {
			return "", "", err
		}
	}

	// Persist the user's message, then build the LLM conversation.
	if err := s.repo.AddMessage(ctx, &domain.CoachMessage{ThreadID: threadID, Role: domain.RoleUser, Content: message}); err != nil {
		return "", "", err
	}

	msgs := make([]domain.LLMMessage, 0, len(history)+1)
	for _, m := range history {
		msgs = append(msgs, domain.LLMMessage{Role: m.Role, Content: m.Content})
	}
	msgs = append(msgs, domain.LLMMessage{Role: domain.RoleUser, Content: message})

	c, err := s.llm.Complete(ctx, coachSystem, msgs, 1500)
	if err != nil {
		return "", "", err
	}
	_ = s.repo.LogInteraction(ctx, userID, "coach", c)

	reply := strings.TrimSpace(c.Text)
	if err := s.repo.AddMessage(ctx, &domain.CoachMessage{ThreadID: threadID, Role: domain.RoleAssistant, Content: reply}); err != nil {
		return "", "", err
	}
	return reply, threadID, nil
}

func (s *Service) ListThreads(ctx context.Context, userID string) ([]domain.CoachThread, error) {
	return s.repo.ListThreads(ctx, userID)
}

func (s *Service) GetThread(ctx context.Context, userID, threadID string) (*domain.CoachThread, error) {
	t, err := s.repo.GetThread(ctx, threadID)
	if err != nil {
		return nil, err
	}
	if t.UserID != userID {
		return nil, ErrForbidden
	}
	if t.Messages, err = s.repo.ListMessages(ctx, threadID); err != nil {
		return nil, err
	}
	return t, nil
}

// ---------- helpers ----------

// decodeJSON tolerates models that wrap JSON in prose or markdown fences by
// extracting the outermost { ... } object.
func decodeJSON(text string, dst any) error {
	t := strings.TrimSpace(text)
	t = strings.TrimPrefix(t, "```json")
	t = strings.TrimPrefix(t, "```")
	t = strings.TrimSuffix(t, "```")
	if i := strings.Index(t, "{"); i >= 0 {
		if j := strings.LastIndex(t, "}"); j > i {
			t = t[i : j+1]
		}
	}
	return json.Unmarshal([]byte(t), dst)
}

func clamp(v int) int {
	if v < 0 {
		return 0
	}
	if v > 100 {
		return 100
	}
	return v
}

func fallback(s, def string) string {
	if strings.TrimSpace(s) == "" {
		return def
	}
	return s
}

func makeTitle(message string) string {
	title := strings.TrimSpace(message)
	if len(title) > 60 {
		title = title[:60] + "…"
	}
	return title
}
