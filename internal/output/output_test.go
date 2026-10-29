package output_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/devradar/devradar/internal/output"
	"github.com/devradar/devradar/pkg/models"
)

func sampleIssues() []models.Issue {
	return []models.Issue{
		{
			ID:            101,
			Number:        1,
			Title:         "Sample issue title",
			Body:          "This is sample issue body",
			HTMLURL:       "https://github.com/sample/repo/issues/1",
			State:         "open",
			CommentsCount: 3,
			RepoFullName:  "sample/repo",
			RepoLanguage:  "Go",
			RepoStars:     1200,
			Difficulty:    models.DifficultyBeginner,
			UpdatedAt:     time.Now().Add(-2 * time.Hour),
			Ranking: &models.RankingScore{
				TotalScore:   92,
				MatchTier:    "Excellent Match",
				Explanations: []string{"• [30/30] Skill match", "• [25/25] Difficulty match"},
			},
		},
	}
}

func TestRenderTable(t *testing.T) {
	var buf bytes.Buffer
	output.RenderTable(&buf, sampleIssues())
	res := buf.String()

	if !strings.Contains(res, "Sample issue title") || !strings.Contains(res, "sample/repo") {
		t.Errorf("table missing issue content: %s", res)
	}

	// Empty table
	var emptyBuf bytes.Buffer
	output.RenderTable(&emptyBuf, []models.Issue{})
	if !strings.Contains(emptyBuf.String(), "No matching issues found") {
		t.Errorf("expected empty message")
	}
}

func TestRenderJSON(t *testing.T) {
	var buf bytes.Buffer
	err := output.RenderJSON(&buf, sampleIssues())
	if err != nil {
		t.Fatalf("render json error: %v", err)
	}
	res := buf.String()
	if !strings.Contains(res, `"sample/repo"`) || !strings.Contains(res, `"total_score": 92`) {
		t.Errorf("unexpected json output: %s", res)
	}
}

func TestRenderCSV(t *testing.T) {
	var buf bytes.Buffer
	err := output.RenderCSV(&buf, sampleIssues())
	if err != nil {
		t.Fatalf("render csv error: %v", err)
	}
	res := buf.String()
	if !strings.Contains(res, "sample/repo") || !strings.Contains(res, "92,Excellent Match") {
		t.Errorf("unexpected csv output: %s", res)
	}
}

func TestRenderIssueDetail(t *testing.T) {
	var buf bytes.Buffer
	issues := sampleIssues()
	output.RenderIssueDetail(&buf, &issues[0])
	res := buf.String()
	if !strings.Contains(res, "Sample issue title") || !strings.Contains(res, "Match Breakdown & Score") {
		t.Errorf("unexpected detail view: %s", res)
	}
}
