package github

import (
	"context"
	"fmt"
	"time"

	"github.com/devradar/devradar/pkg/models"
)

// GHRepoResponse represents raw GitHub repository response.
type GHRepoResponse struct {
	ID              int64       `json:"id"`
	Name            string      `json:"name"`
	FullName        string      `json:"full_name"`
	Owner           models.User `json:"owner"`
	Private         bool        `json:"private"`
	HTMLURL         string      `json:"html_url"`
	Description     string      `json:"description"`
	Fork            bool        `json:"fork"`
	CreatedAt       string      `json:"created_at"`
	UpdatedAt       string      `json:"updated_at"`
	PushedAt        string      `json:"pushed_at"`
	Homepage        string      `json:"homepage"`
	Size            int         `json:"size"`
	StargazersCount int         `json:"stargazers_count"`
	WatchersCount   int         `json:"watchers_count"`
	Language        string      `json:"language"`
	ForksCount      int         `json:"forks_count"`
	OpenIssuesCount int         `json:"open_issues_count"`
	Archived        bool        `json:"archived"`
	Disabled        bool        `json:"disabled"`
	DefaultBranch   string      `json:"default_branch"`
}

// GetRepo fetches repository metadata for owner/repo.
func (c *Client) GetRepo(ctx context.Context, owner, repo string) (*models.Repository, error) {
	endpoint := fmt.Sprintf("/repos/%s/%s", owner, repo)

	var raw GHRepoResponse
	if err := c.Get(ctx, endpoint, &raw); err != nil {
		return nil, err
	}

	r := &models.Repository{
		ID:              raw.ID,
		Name:            raw.Name,
		FullName:        raw.FullName,
		Owner:           raw.Owner,
		Private:         raw.Private,
		HTMLURL:         raw.HTMLURL,
		Description:     raw.Description,
		Fork:            raw.Fork,
		CreatedAt:       parseGitHubTime(raw.CreatedAt),
		UpdatedAt:       parseGitHubTime(raw.UpdatedAt),
		PushedAt:        parseGitHubTime(raw.PushedAt),
		Homepage:        raw.Homepage,
		Size:            raw.Size,
		StargazersCount: raw.StargazersCount,
		WatchersCount:   raw.WatchersCount,
		Language:        raw.Language,
		ForksCount:      raw.ForksCount,
		OpenIssuesCount: raw.OpenIssuesCount,
		Archived:        raw.Archived,
		Disabled:        raw.Disabled,
		DefaultBranch:   raw.DefaultBranch,
	}

	return r, nil
}

// GetIssue fetches a single issue by repo owner, repo name, and issue number.
func (c *Client) GetIssue(ctx context.Context, owner, repo string, number int) (*models.Issue, error) {
	endpoint := fmt.Sprintf("/repos/%s/%s/issues/%d", owner, repo, number)

	var item GHIssueItem
	if err := c.Get(ctx, endpoint, &item); err != nil {
		return nil, err
	}

	issue := item.ToModel()
	issue.RepoOwner = owner
	issue.RepoName = repo
	issue.RepoFullName = fmt.Sprintf("%s/%s", owner, repo)
	return &issue, nil
}

func parseGitHubTime(str string) time.Time {
	if str == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return time.Time{}
	}
	return t
}
