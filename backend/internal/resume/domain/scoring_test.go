package domain

import (
	"strings"
	"testing"
)

func goodResume() string {
	// ~450 words, with contact info, sections, action verbs and quantified wins.
	header := "Jane Doe  jane.doe@example.com  +1 415 555 0199\n\n"
	summary := "Summary: Operations leader with 12 years of experience. "
	exp := "Experience: Led a team of 20 and improved throughput 35%. Managed a $2M budget and reduced costs 18%. Delivered projects ahead of schedule. Built new processes and increased retention 22%. "
	edu := "Education: BSc in Business, graduated 2009. "
	skills := "Skills: leadership, budgeting, analytics, project management. "
	filler := strings.Repeat("Coordinated cross-functional teams and delivered measurable results across regions. ", 20)
	return header + summary + exp + edu + skills + filler
}

func TestScoreText_GoodResume(t *testing.T) {
	s := ScoreText(goodResume())
	if s.Overall < 60 {
		t.Errorf("expected a solid overall score, got %d", s.Overall)
	}
	if s.ATS < 70 {
		t.Errorf("expected high ATS (contact + sections present), got %d", s.ATS)
	}
	if s.Formatting < 60 {
		t.Errorf("expected decent formatting, got %d", s.Formatting)
	}
}

func TestScoreText_EmptyText(t *testing.T) {
	s := ScoreText("")
	if s.Overall != 0 && s.Overall > 30 {
		t.Errorf("empty resume should score very low, got %d", s.Overall)
	}
	if !containsSuggestion(s.Suggestions, "extract text") {
		t.Errorf("expected an extraction suggestion, got %v", s.Suggestions)
	}
}

func TestScoreText_MissingSectionsAndContact(t *testing.T) {
	s := ScoreText("Just a short blurb about myself with no structure.")
	if !containsSuggestion(s.Suggestions, "Skills section") {
		t.Errorf("expected a skills-section suggestion, got %v", s.Suggestions)
	}
	if !containsSuggestion(s.Suggestions, "contact email") {
		t.Errorf("expected a contact-email suggestion, got %v", s.Suggestions)
	}
}

func containsSuggestion(list []string, substr string) bool {
	for _, s := range list {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}
