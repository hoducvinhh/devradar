package models

import (
	"strings"
	"time"
)

// Label represents a GitHub issue label.
type Label struct {
	ID          int64  `json:"id,omitempty"`
	Name        string `json:"name"`
	Color       string `json:"color,omitempty"`
	Description string `json:"description,omitempty"`
}

// User represents a GitHub user.
type User struct {
	Login     string `json:"login"`
	ID        int64  `json:"id,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
	HTMLURL   string `json:"html_url,omitempty"`
}

// Issue represents a GitHub issue with DevRadar enrichment.
type Issue struct {
	ID            int64      `json:"id"`
	Number        int        `json:"number"`
	Title         string     `json:"title"`
	Body          string     `json:"body"`
	HTMLURL       string     `json:"html_url"`
	State         string     `json:"state"`
	CommentsCount int        `json:"comments"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	ClosedAt      *time.Time `json:"closed_at,omitempty"`

	User      User    `json:"user"`
	Assignee  *User   `json:"assignee,omitempty"`
	Assignees []User  `json:"assignees,omitempty"`
	Labels    []Label `json:"labels"`

	// Repository context
	RepoOwner    string    `json:"repo_owner,omitempty"`
	RepoName     string    `json:"repo_name,omitempty"`
	RepoFullName string    `json:"repo_full_name,omitempty"`
	RepoLanguage string    `json:"repo_language,omitempty"`
	RepoStars    int       `json:"repo_stars,omitempty"`
	RepoForks    int       `json:"repo_forks,omitempty"`
	RepoPushedAt time.Time `json:"repo_pushed_at,omitempty"`
	RepoArchived bool      `json:"repo_archived,omitempty"`

	// DevRadar Intelligence
	Difficulty       DifficultyLevel `json:"difficulty"`
	DifficultyReason string          `json:"difficulty_reason,omitempty"`
	Ranking          *RankingScore   `json:"ranking,omitempty"`
}

// IsAssigned returns true if the issue is currently assigned to any contributor.
func (i *Issue) IsAssigned() bool {
	return i.Assignee != nil || len(i.Assignees) > 0
}

// HasLabel returns true if the issue has a specific label (case-insensitive).
func (i *Issue) HasLabel(name string) bool {
	for _, l := range i.Labels {
		if strings.EqualFold(l.Name, name) {
			return true
		}
	}
	return false
}
