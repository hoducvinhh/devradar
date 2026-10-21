package models

import "time"

// Profile represents local user preferences for personalized issue discovery.
type Profile struct {
	Level       DifficultyLevel `json:"level"`
	Languages   []string        `json:"languages"`
	Interests   []string        `json:"interests"`
	MinStars    int             `json:"min_stars,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// DefaultProfile returns a clean default profile for new users.
func DefaultProfile() *Profile {
	return &Profile{
		Level:     DifficultyBeginner,
		Languages: []string{},
		Interests: []string{},
		MinStars:  0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Bookmark represents a saved issue stored in ~/.devradar/bookmarks.json
type Bookmark struct {
	IssueID      int64     `json:"issue_id"`
	Number       int       `json:"number"`
	Title        string    `json:"title"`
	HTMLURL      string    `json:"html_url"`
	RepoFullName string    `json:"repo_full_name"`
	Language     string    `json:"language"`
	Difficulty   DifficultyLevel `json:"difficulty"`
	SavedAt      time.Time `json:"saved_at"`
	Notes        string    `json:"notes,omitempty"`
}
