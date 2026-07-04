package domain

import (
	"regexp"
	"strings"
)

var (
	reEmail   = regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)
	rePhone   = regexp.MustCompile(`(\+?\d[\d\s\-().]{7,}\d)`)
	reNumbers = regexp.MustCompile(`(\d+%|\$\s?\d+|\b\d{2,}\b)`)
	reWord    = regexp.MustCompile(`[A-Za-z][A-Za-z'\-]+`)
)

// sectionKeywords maps a resume section to detection terms.
var sectionKeywords = map[string][]string{
	"experience":     {"experience", "employment", "work history"},
	"education":      {"education", "academic"},
	"skills":         {"skills", "competencies", "technologies"},
	"summary":        {"summary", "objective", "profile"},
	"projects":       {"projects", "portfolio"},
	"certifications": {"certification", "certificate", "licenses"},
}

var actionVerbs = []string{
	"led", "managed", "built", "designed", "developed", "delivered", "launched",
	"improved", "increased", "reduced", "created", "implemented", "achieved",
	"coordinated", "owned", "drove", "optimized", "negotiated", "streamlined",
}

// ScoreText evaluates resume text and returns sub-scores (0-100) plus
// actionable suggestions. Pure function — no I/O.
func ScoreText(text string) Score {
	lower := strings.ToLower(text)
	words := reWord.FindAllString(text, -1)
	wordCount := len(words)

	suggestions := []string{}

	// --- sections present ---
	present := 0
	missing := []string{}
	for _, name := range []string{"summary", "experience", "education", "skills"} {
		if sectionPresent(lower, name) {
			present++
		} else {
			missing = append(missing, name)
		}
	}
	for _, name := range missing {
		suggestions = append(suggestions, "Add a "+sectionLabel(name)+" section.")
	}

	// --- length ---
	lengthScore := lengthScore(wordCount)
	if wordCount < 200 {
		suggestions = append(suggestions, "Your resume looks short — aim for 400+ words with concrete detail.")
	} else if wordCount > 1200 {
		suggestions = append(suggestions, "Your resume is long — tighten it toward 1–2 pages.")
	}

	// Formatting: 60% sections (of 4 core), 40% length.
	formatting := clamp((present*100/4)*60/100 + lengthScore*40/100)

	// --- keywords: action verbs + quantified achievements ---
	verbHits := countAny(lower, actionVerbs)
	quantHits := len(reNumbers.FindAllString(text, -1))
	verbScore := minInt(verbHits*12, 60)   // up to 60
	quantScore := minInt(quantHits*10, 40) // up to 40
	keywords := clamp(verbScore + quantScore)
	if verbHits < 3 {
		suggestions = append(suggestions, "Use more strong action verbs (led, built, improved, delivered).")
	}
	if quantHits < 2 {
		suggestions = append(suggestions, "Quantify achievements with numbers and percentages (e.g., \"cut costs 18%\").")
	}

	// --- ATS: contact + sections + parseable text ---
	hasEmail := reEmail.MatchString(text)
	hasPhone := rePhone.MatchString(text)
	contact := 0
	if hasEmail {
		contact += 18
	}
	if hasPhone {
		contact += 12
	}
	parseable := 0
	if wordCount > 0 {
		parseable = 30
	}
	ats := clamp(contact + (present * 40 / 4) + parseable)
	if !hasEmail {
		suggestions = append(suggestions, "Add a contact email so recruiters and ATS can reach you.")
	}
	if !hasPhone {
		suggestions = append(suggestions, "Add a phone number to your contact details.")
	}
	if wordCount == 0 {
		suggestions = append(suggestions, "We couldn't extract text — upload a text-based PDF or DOCX (not a scanned image).")
	}

	overall := clamp((formatting + keywords + ats) / 3)

	return Score{
		Overall: overall, Formatting: formatting, Keywords: keywords, ATS: ats,
		Suggestions: suggestions,
	}
}

func sectionPresent(lower, name string) bool {
	for _, kw := range sectionKeywords[name] {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

func sectionLabel(name string) string {
	switch name {
	case "summary":
		return "Summary/Objective"
	default:
		return strings.Title(name) //nolint:staticcheck // simple capitalization is fine here
	}
}

func countAny(lower string, terms []string) int {
	total := 0
	for _, t := range terms {
		total += strings.Count(lower, t)
	}
	return total
}

func lengthScore(words int) int {
	switch {
	case words >= 400 && words <= 1000:
		return 100
	case words >= 250 && words < 400:
		return 75
	case words > 1000 && words <= 1200:
		return 80
	case words >= 150 && words < 250:
		return 50
	case words == 0:
		return 0
	default:
		return 30
	}
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

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
