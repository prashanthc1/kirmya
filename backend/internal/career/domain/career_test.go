package domain

import "testing"

func TestPathsCuratedOperations(t *testing.T) {
	p := Paths("Operations Manager")
	if p.Target != "VP of Operations" {
		t.Fatalf("expected curated operations target, got %q", p.Target)
	}
	if len(p.Steps) != 4 {
		t.Fatalf("expected 4 rungs, got %d", len(p.Steps))
	}
	if !p.Steps[0].Current {
		t.Fatalf("first rung must be marked current")
	}
	if p.Steps[0].PayBand == "" || len(p.GapSkills) == 0 {
		t.Fatalf("expected pay band and gap skills, got %+v", p)
	}
}

func TestPathsCaseInsensitiveAndKeywords(t *testing.T) {
	// Lower-cased and a different-but-related title should both resolve to the
	// operations family.
	if got := Paths("senior supply chain analyst"); got.Target != "VP of Operations" {
		t.Fatalf("expected operations family for supply-chain role, got %q", got.Target)
	}
	if got := Paths("Staff Software Developer"); got.Target != "Staff Engineer" {
		t.Fatalf("expected engineering family, got %q", got.Target)
	}
}

func TestPathsFallbackGeneric(t *testing.T) {
	p := Paths("Underwater Basket Weaver")
	if p.From != "Underwater Basket Weaver" {
		t.Fatalf("expected From echoed, got %q", p.From)
	}
	if len(p.Steps) < 2 || !p.Steps[0].Current {
		t.Fatalf("expected a generic ladder with a current first rung, got %+v", p.Steps)
	}
	if p.Steps[0].Title != "Underwater Basket Weaver" {
		t.Fatalf("expected first rung to be the supplied role, got %q", p.Steps[0].Title)
	}
}
