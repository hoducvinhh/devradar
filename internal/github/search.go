package github

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/devradar/devradar/pkg/models"
)

// SearchIssuesResult represents GitHub's search issues response.
type SearchIssuesResult struct {
	TotalCount        int             `json:"total_count"`
	IncompleteResults bool            `json:"incomplete_results"`
	Items             []GHIssueItem   `json:"items"`
}

// GHIssueItem maps directly to GitHub's raw search item.
type GHIssueItem struct {
	ID            int64           `json:"id"`
	Number        int             `json:"number"`
	Title         string          `json:"title"`
	Body          string          `json:"body"`
	HTMLURL       string          `json:"html_url"`
	RepositoryURL string          `json:"repository_url"`
	State         string          `json:"state"`
	CommentsCount int             `json:"comments"`
	CreatedAt     string          `json:"created_at"`
	UpdatedAt     string          `json:"updated_at"`
	User          models.User     `json:"user"`
	Assignee      *models.User    `json:"assignee"`
	Assignees     []models.User   `json:"assignees"`
	Labels        []models.Label  `json:"labels"`
	PullRequest   *struct{}       `json:"pull_request,omitempty"` // Present if item is a PR
}

// BuildSearchQuery constructs the GitHub search syntax according to filter options.
func BuildSearchQuery(opts *models.FilterOptions) string {
	var parts []string

	// Fundamental requirements: always open issues, never pull requests
	parts = append(parts, "is:issue", "is:open")

	if opts == nil {
		return strings.Join(parts, " ")
	}

	if opts.Language != "" {
		parts = append(parts, fmt.Sprintf("language:%s", strings.TrimSpace(opts.Language)))
	}

	if opts.Label != "" {
		label := strings.TrimSpace(opts.Label)
		if strings.Contains(label, " ") {
			parts = append(parts, fmt.Sprintf("label:\"%s\"", label))
		} else {
			parts = append(parts, fmt.Sprintf("label:%s", label))
		}
	}

	if opts.Unassigned {
		parts = append(parts, "no:assignee")
	}

	if opts.Org != "" {
		parts = append(parts, fmt.Sprintf("org:%s", strings.TrimSpace(opts.Org)))
	}

	if opts.Repo != "" {
		parts = append(parts, fmt.Sprintf("repo:%s", strings.TrimSpace(opts.Repo)))
	}

	if opts.Keywords != "" {
		parts = append(parts, strings.TrimSpace(opts.Keywords))
	}

	return strings.Join(parts, " ")
}

// SearchIssues calls GitHub API to search issues matching the filter.
func (c *Client) SearchIssues(ctx context.Context, opts *models.FilterOptions) ([]models.Issue, error) {
	q := BuildSearchQuery(opts)

	params := url.Values{}
	params.Set("q", q)

	limit := 15
	if opts != nil && opts.Limit > 0 {
		limit = opts.Limit
	}
	if limit > 100 {
		limit = 100
	}
	params.Set("per_page", strconv.Itoa(limit))

	sortField := "updated"
	if opts != nil && opts.Sort != "" {
		sortField = opts.Sort
	}
	params.Set("sort", sortField)

	orderDir := "desc"
	if opts != nil && opts.Order != "" {
		orderDir = opts.Order
	}
	params.Set("order", orderDir)

	endpoint := fmt.Sprintf("/search/issues?%s", params.Encode())

	var result SearchIssuesResult
	if err := c.Get(ctx, endpoint, &result); err != nil {
		return nil, err
	}

	var issues []models.Issue
	seen := make(map[int64]bool)

	for _, item := range result.Items {
		// Strict check: GitHub search sometimes includes PRs if not filtered out
		if item.PullRequest != nil {
			continue
		}

		if seen[item.ID] {
			continue
		}
		seen[item.ID] = true

		issue := item.ToModel()
		issues = append(issues, issue)
	}

	return issues, nil
}

// ToModel transforms a raw GitHub issue item into models.Issue.
func (item *GHIssueItem) ToModel() models.Issue {
	issue := models.Issue{
		ID:            item.ID,
		Number:        item.Number,
		Title:         item.Title,
		Body:          item.Body,
		HTMLURL:       item.HTMLURL,
		State:         item.State,
		CommentsCount: item.CommentsCount,
		User:          item.User,
		Assignee:      item.Assignee,
		Assignees:     item.Assignees,
		Labels:        item.Labels,
	}

	// Parse timestamps
	issue.CreatedAt = parseGitHubTime(item.CreatedAt)
	issue.UpdatedAt = parseGitHubTime(item.UpdatedAt)

	// Extract repository info from repository_url: https://api.github.com/repos/:owner/:repo
	if item.RepositoryURL != "" {
		parts := strings.Split(item.RepositoryURL, "/repos/")
		if len(parts) == 2 {
			repoParts := strings.Split(parts[1], "/")
			if len(repoParts) >= 2 {
				issue.RepoOwner = repoParts[0]
				issue.RepoName = repoParts[1]
				issue.RepoFullName = fmt.Sprintf("%s/%s", repoParts[0], repoParts[1])
			}
		}
	} else if item.HTMLURL != "" {
		// Fallback: parse from html_url: https://github.com/:owner/:repo/issues/:number
		parts := strings.Split(item.HTMLURL, "github.com/")
		if len(parts) == 2 {
			subParts := strings.Split(parts[1], "/")
			if len(subParts) >= 2 {
				issue.RepoOwner = subParts[0]
				issue.RepoName = subParts[1]
				issue.RepoFullName = fmt.Sprintf("%s/%s", subParts[0], subParts[1])
			}
		}
	}

	return issue
}
