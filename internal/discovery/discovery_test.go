package discovery_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/devradar/devradar/internal/discovery"
	"github.com/devradar/devradar/internal/github"
	"github.com/devradar/devradar/internal/storage"
	"github.com/devradar/devradar/pkg/models"
)

func TestDiscovery_Discover(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/search/issues" {
			res := `{
				"total_count": 2,
				"incomplete_results": false,
				"items": [
					{
						"id": 1,
						"number": 10,
						"title": "Fix typo in doc",
						"body": "Small spelling correction",
						"html_url": "https://github.com/org/repo/issues/10",
						"repository_url": "https://api.github.com/repos/org/repo",
						"state": "open",
						"comments": 1,
						"labels": [ { "name": "good first issue" } ]
					},
					{
						"id": 2,
						"number": 11,
						"title": "Severe deadlock in raft",
						"body": "Distributed consensus deadlock",
						"html_url": "https://github.com/org/repo/issues/11",
						"repository_url": "https://api.github.com/repos/org/repo",
						"state": "open",
						"comments": 18,
						"labels": [ { "name": "concurrency" } ]
					}
				]
			}`
			_, _ = w.Write([]byte(res))
			return
		}
		if r.URL.Path == "/repos/org/repo" {
			res := `{
				"id": 10,
				"name": "repo",
				"full_name": "org/repo",
				"stargazers_count": 800,
				"language": "Go",
				"pushed_at": "2026-10-25T10:00:00Z"
			}`
			_, _ = w.Write([]byte(res))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	tmpDir, _ := os.MkdirTemp("", "devradar-disc-test-*")
	defer os.RemoveAll(tmpDir)

	store, _ := storage.NewStorageWithDir(tmpDir)
	cache := storage.NewCache(tmpDir, 10*time.Minute)
	client := github.NewClient(github.WithBaseURL(server.URL))

	svc := discovery.NewService(client, store, cache)

	opts := &models.FilterOptions{
		Language: "go",
		Limit:    10,
	}

	issues, err := svc.Discover(context.Background(), opts)
	if err != nil {
		t.Fatalf("discovery error: %v", err)
	}

	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(issues))
	}

	// Verify ranking is applied and sorted descending
	if issues[0].Ranking == nil || issues[1].Ranking == nil {
		t.Fatalf("issues missing ranking scores")
	}

	if issues[0].Ranking.TotalScore < issues[1].Ranking.TotalScore {
		t.Errorf("expected issues sorted descending by score")
	}
}
