package domain

// Match is a scored job recommendation for a seeker. Score is 0–100.
type Match struct {
	Job           Job
	Score         int
	MatchedSkills []string
	MissingSkills []string
	Reason        string
}
