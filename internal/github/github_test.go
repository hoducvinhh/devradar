package github_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/devradar/devradar/internal/github"
	"github.com/devradar/devradar/pkg/models"
)

func TestBuildSearchQuery(t *testing.T) {
	// 1. Default
	q := github.BuildSearchQuery(nil)
	if !strings.Contains(q, "is:issue") || !strings.Contains(q, "is:open") {
		t.Errorf("expected is:issue is:open, got: %s", q)
	}

	// 2. Full options
	opts := &models.FilterOptions{
		Language:   "go",
		Label:      "good first issue",
		Unassigned: true,
		Org:        "kubernetes",
		Repo:       "kubernetes/kubectl",
		Keywords:   "crash memory",
	}

	fullQ := github.BuildSearchQuery(opts)
	expectedTokens := []string{
		"is:issue",
		"is:open",
		"language:go",
		"label:\"good first issue\"",
		"no:assignee",
		"org:kubernetes",
		"repo:kubernetes/kubectl",
		"crash memory",
	}

	for _, token := range expectedTokens {
		if !strings.Contains(fullQ, token) {
			t.Errorf("expected query to contain %q, but query was %q", token, fullQ)
		}
	}
}

func TestRateLimit_Tracking(t *testing.T) {
	rl := &github.RateLimit{}

	header := make(http.Header)
	header.Set("X-RateLimit-Limit", "5000")
	header.Set("X-RateLimit-Remaining", "4990")
	header.Set("X-RateLimit-Used", "10")
	header.Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(1*time.Hour).Unix()))

	rl.UpdateFromHeaders(header)

	limit, remaining, _ := rl.Status()
	if limit != 5000 || remaining != 4990 {
		t.Errorf("unexpected rate limit status: %d / %d", remaining, limit)
	}

	if rl.IsExceeded() {
		t.Errorf("rate limit should not be exceeded")
	}

	// Test exceeded
	header.Set("X-RateLimit-Remaining", "0")
	rl.UpdateFromHeaders(header)
	if !rl.IsExceeded() {
		t.Errorf("rate limit should be exceeded")
	}
	if msg := rl.WarningMessage(); msg == "" {
		t.Errorf("expected warning message when quota exhausted")
	}
}

func TestClient_MockAPICalls(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-RateLimit-Limit", "60")
		w.Header().Set("X-RateLimit-Remaining", "58")
		w.Header().Set("Content-Type", "application/json")

		switch {
		case strings.HasPrefix(r.URL.Path, "/search/issues"):
			response := `{
				"total_count": 1,
				"incomplete_results": false,
				"items": [
					{
						"id": 1001,
						"number": 42,
						"title": "Fix memory leak in buffer",
						"body": "Buffer is not freed properly",
						"html_url": "https://github.com/mock/repo/issues/42",
						"repository_url": "https://api.github.com/repos/mock/repo",
						"state": "open",
						"comments": 2,
						"created_at": "2026-01-01T10:00:00Z",
						"updated_at": "2026-01-02T10:00:00Z",
						"user": { "login": "tester" },
						"labels": [ { "name": "bug" } ]
					}
				]
			}`
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(response))

		case r.URL.Path == "/repos/mock/repo":
			response := `{
				"id": 501,
				"name": "repo",
				"full_name": "mock/repo",
				"stargazers_count": 1500,
				"forks_count": 120,
				"language": "Go",
				"pushed_at": "2026-01-02T10:00:00Z"
			}`
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(response))

		case r.URL.Path == "/repos/mock/repo/issues/42":
			response := `{
				"id": 1001,
				"number": 42,
				"title": "Fix memory leak in buffer",
				"body": "Buffer is not freed properly",
				"html_url": "https://github.com/mock/repo/issues/42",
				"state": "open",
				"comments": 2,
				"user": { "login": "tester" },
				"labels": [ { "name": "bug" } ]
			}`
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(response))

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := github.NewClient(
		github.WithBaseURL(server.URL),
		github.WithToken("test-token"),
	)

	ctx := context.Background()

	// 1. Search issues
	issues, err := client.SearchIssues(ctx, &models.FilterOptions{Language: "go"})
	if err != nil {
		t.Fatalf("search issues failed: %v", err)
	}
	if len(issues) != 1 || issues[0].Number != 42 {
		t.Fatalf("unexpected issues response: %+v", issues)
	}
	if issues[0].RepoFullName != "mock/repo" {
		t.Errorf("expected repo full name mock/repo, got %s", issues[0].RepoFullName)
	}

	// 2. Get repo
	repo, err := client.GetRepo(ctx, "mock", "repo")
	if err != nil {
		t.Fatalf("get repo failed: %v", err)
	}
	if repo.StargazersCount != 1500 || repo.Language != "Go" {
		t.Errorf("unexpected repo data: %+v", repo)
	}

	// 3. Get single issue
	issue, err := client.GetIssue(ctx, "mock", "repo", 42)
	if err != nil {
		t.Fatalf("get issue failed: %v", err)
	}
	if issue.Title != "Fix memory leak in buffer" {
		t.Errorf("unexpected issue title: %s", issue.Title)
	}
}
