package models

import "time"

// Repository represents repository metadata returned by GitHub.
type Repository struct {
	ID              int64     `json:"id"`
	NodeID          string    `json:"node_id,omitempty"`
	Name            string    `json:"name"`
	FullName        string    `json:"full_name"`
	Owner           User      `json:"owner"`
	Private         bool      `json:"private"`
	HTMLURL         string    `json:"html_url"`
	Description     string    `json:"description"`
	Fork            bool      `json:"fork"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	PushedAt        time.Time `json:"pushed_at"`
	Homepage        string    `json:"homepage,omitempty"`
	Size            int       `json:"size"`
	StargazersCount int       `json:"stargazers_count"`
	WatchersCount   int       `json:"watchers_count"`
	Language        string    `json:"language"`
	ForksCount      int       `json:"forks_count"`
	OpenIssuesCount int       `json:"open_issues_count"`
	Archived        bool      `json:"archived"`
	Disabled        bool      `json:"disabled"`
	License         *License  `json:"license,omitempty"`
	DefaultBranch   string    `json:"default_branch"`

	// DevRadar Analysis
	HealthScore     float64   `json:"health_score,omitempty"`
	HealthReason    string    `json:"health_reason,omitempty"`
}

// License represents open source license info.
type License struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	SPDXID string `json:"spdx_id"`
	URL    string `json:"url"`
}
