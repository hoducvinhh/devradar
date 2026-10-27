package ranking_test

import (
	"testing"
	"time"

	"github.com/devradar/devradar/internal/ranking"
	"github.com/devradar/devradar/pkg/models"
)

func TestCalculateSkillScore(t *testing.T) {
	issue := &models.Issue{
		RepoLanguage: "Go",
		Title:        "CLI tool improvements",
	}

	// Exact target match
	score, exp := ranking.CalculateSkillScore(issue, nil, "Go")
	if score != 100.0 {
		t.Errorf("expected 100, got %f (%s)", score, exp)
	}

	// Profile language match
	prof := &models.Profile{
		Languages: []string{"go", "rust"},
	}
	score, exp = ranking.CalculateSkillScore(issue, prof, "")
	if score != 100.0 {
		t.Errorf("expected 100, got %f (%s)", score, exp)
	}

	// Profile interest match
	profInterests := &models.Profile{
		Interests: []string{"cli"},
	}
	score, exp = ranking.CalculateSkillScore(issue, profInterests, "")
	if score < 70.0 {
		t.Errorf("expected >= 70 for interest match, got %f (%s)", score, exp)
	}
}

func TestCalculateDifficultyScore(t *testing.T) {
	// Exact match
	score, _ := ranking.CalculateDifficultyScore(models.DifficultyBeginner, models.DifficultyBeginner)
	if score != 100.0 {
		t.Errorf("expected 100, got %f", score)
	}

	// 1 level difference (Beginner vs Intermediate)
	score, _ = ranking.CalculateDifficultyScore(models.DifficultyIntermediate, models.DifficultyBeginner)
	if score != 50.0 {
		t.Errorf("expected 50, got %f", score)
	}

	// 2 levels difference (Beginner vs Advanced)
	score, _ = ranking.CalculateDifficultyScore(models.DifficultyAdvanced, models.DifficultyBeginner)
	if score != 15.0 {
		t.Errorf("expected 15, got %f", score)
	}

	// Open level
	score, _ = ranking.CalculateDifficultyScore(models.DifficultyIntermediate, "")
	if score != 80.0 {
		t.Errorf("expected 80, got %f", score)
	}
}

func TestCalculateFreshnessScore(t *testing.T) {
	// 1 day old
	score, _ := ranking.CalculateFreshnessScore(time.Now().Add(-24 * time.Hour))
	if score != 100.0 {
		t.Errorf("expected 100, got %f", score)
	}

	// 10 days old
	score, _ = ranking.CalculateFreshnessScore(time.Now().Add(-10 * 24 * time.Hour))
	if score != 75.0 {
		t.Errorf("expected 75, got %f", score)
	}

	// 120 days old
	score, _ = ranking.CalculateFreshnessScore(time.Now().Add(-120 * 24 * time.Hour))
	if score != 20.0 {
		t.Errorf("expected 20, got %f", score)
	}
}

func TestCalculateClarityScore(t *testing.T) {
	// Empty body
	score, _ := ranking.CalculateClarityScore("")
	if score != 20.0 {
		t.Errorf("expected 20 for empty body, got %f", score)
	}

	// Rich description with code and checklist
	richBody := `
## Steps to Reproduce
- [ ] Run the server
- [ ] Send request

` + "```go\nfmt.Println(\"test\")\n```"

	score, exp := ranking.CalculateClarityScore(richBody)
	if score < 80.0 {
		t.Errorf("expected >= 80 for rich body, got %f (%s)", score, exp)
	}
}

func TestCalculateUnassignedScore(t *testing.T) {
	if s, _ := ranking.CalculateUnassignedScore(false); s != 100.0 {
		t.Errorf("expected 100 for unassigned, got %f", s)
	}
	if s, _ := ranking.CalculateUnassignedScore(true); s != 0.0 {
		t.Errorf("expected 0 for assigned, got %f", s)
	}
}

func TestRanker_FullIntegration(t *testing.T) {
	ranker := ranking.NewRanker()
	issue := &models.Issue{
		Title:        "Implement lightweight cache",
		Body:         "## Task description\n- [ ] Add TTL cache\n```go\ntype Cache struct{}\n```",
		RepoLanguage: "Go",
		UpdatedAt:    time.Now().Add(-2 * time.Hour),
		Difficulty:   models.DifficultyBeginner,
	}

	prof := &models.Profile{
		Level:     models.DifficultyBeginner,
		Languages: []string{"go"},
	}

	res := ranker.Rank(issue, prof, "Go", models.DifficultyBeginner, 100.0, "Active repo")
	if res.TotalScore < 90 {
		t.Errorf("expected high match score (>=90), got %d", res.TotalScore)
	}
	if len(res.Explanations) != 6 {
		t.Errorf("expected 6 explanation lines, got %d", len(res.Explanations))
	}
	if res.MatchTier != "Excellent Match" {
		t.Errorf("expected Excellent Match tier, got %s", res.MatchTier)
	}
}
