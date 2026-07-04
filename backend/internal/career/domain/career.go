// Package domain holds the Career bounded context. Career paths are curated
// reference data (role ladders, pay bands and the skill gap to the next target),
// so this context owns no tables and needs no infrastructure adapter — the data
// lives here and the application layer reads it directly.
package domain

import "strings"

// Step is one rung on a career ladder.
type Step struct {
	Title   string `json:"title"`
	PayBand string `json:"pay_band"`
	// Current marks the rung the user is starting from.
	Current bool `json:"current,omitempty"`
}

// Path is the ladder of roles reachable from a starting role, the headline
// target rung, and the skills standing between the user and that target.
type Path struct {
	From      string   `json:"from"`
	Steps     []Step   `json:"steps"`
	Target    string   `json:"target"`
	GapSkills []string `json:"gap_skills"`
}

type ladder struct {
	steps     []Step
	target    string
	gapSkills []string
}

// curated is the library of hand-authored ladders, keyed by a coarse role family.
var curated = map[string]ladder{
	"operations": {
		steps: []Step{
			{Title: "Operations Manager", PayBand: "$95k–$130k"},
			{Title: "Director of Operations", PayBand: "$140k–$190k"},
			{Title: "VP of Operations", PayBand: "$200k–$280k"},
			{Title: "Head of Supply Chain", PayBand: "$210k–$300k"},
		},
		target:    "VP of Operations",
		gapSkills: []string{"Board & exec communication", "Org design at scale", "Network strategy", "P&L ownership at scale"},
	},
	"engineering": {
		steps: []Step{
			{Title: "Software Engineer", PayBand: "$110k–$150k"},
			{Title: "Senior Software Engineer", PayBand: "$150k–$200k"},
			{Title: "Staff Engineer", PayBand: "$200k–$260k"},
			{Title: "Engineering Manager", PayBand: "$190k–$250k"},
		},
		target:    "Staff Engineer",
		gapSkills: []string{"System design at scale", "Technical leadership", "Cross-team influence", "Mentoring & hiring"},
	},
	"product": {
		steps: []Step{
			{Title: "Product Manager", PayBand: "$120k–$160k"},
			{Title: "Senior Product Manager", PayBand: "$160k–$210k"},
			{Title: "Director of Product", PayBand: "$200k–$270k"},
			{Title: "VP of Product", PayBand: "$250k–$340k"},
		},
		target:    "Director of Product",
		gapSkills: []string{"Product strategy", "Org design", "Executive storytelling", "Portfolio prioritization"},
	},
	"sales": {
		steps: []Step{
			{Title: "Account Executive", PayBand: "$80k–$140k OTE"},
			{Title: "Senior Account Executive", PayBand: "$140k–$200k OTE"},
			{Title: "Sales Manager", PayBand: "$170k–$230k OTE"},
			{Title: "Director of Sales", PayBand: "$220k–$320k OTE"},
		},
		target:    "Sales Manager",
		gapSkills: []string{"Team leadership", "Forecasting & pipeline ops", "Strategic accounts", "Hiring & coaching"},
	},
}

// familyOf maps a free-text role to a curated ladder family, or "" if unknown.
func familyOf(role string) string {
	r := strings.ToLower(role)
	switch {
	case strings.Contains(r, "operation") || strings.Contains(r, "supply chain") || strings.Contains(r, "logistics"):
		return "operations"
	case strings.Contains(r, "engineer") || strings.Contains(r, "developer") || strings.Contains(r, "software"):
		return "engineering"
	case strings.Contains(r, "product"):
		return "product"
	case strings.Contains(r, "sales") || strings.Contains(r, "account executive"):
		return "sales"
	default:
		return ""
	}
}

// genericLadder builds a role-agnostic progression for roles outside the curated
// library, so every caller still gets a usable answer.
func genericLadder(from string) ladder {
	from = strings.TrimSpace(from)
	if from == "" {
		from = "Your current role"
	}
	return ladder{
		steps: []Step{
			{Title: from, PayBand: "—"},
			{Title: "Senior " + from, PayBand: "—"},
			{Title: "Lead / Manager", PayBand: "—"},
			{Title: "Director", PayBand: "—"},
		},
		target:    "Lead / Manager",
		gapSkills: []string{"Cross-functional leadership", "Strategic planning", "Executive communication", "Stakeholder management"},
	}
}

// Paths returns the career ladder for a starting role. The first rung is marked
// Current. Unknown roles fall back to a generic progression.
func Paths(from string) Path {
	l, ok := curated[familyOf(from)]
	if !ok {
		l = genericLadder(from)
	}
	steps := make([]Step, len(l.steps))
	copy(steps, l.steps)
	if len(steps) > 0 {
		steps[0].Current = true
	}
	return Path{
		From:      strings.TrimSpace(from),
		Steps:     steps,
		Target:    l.target,
		GapSkills: l.gapSkills,
	}
}
