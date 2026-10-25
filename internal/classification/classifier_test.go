package classification_test

import (
	"testing"

	"github.com/devradar/devradar/internal/classification"
	"github.com/devradar/devradar/pkg/models"
)

func TestClassifier_Beginner(t *testing.T) {
	c := classification.NewClassifier()

	tests := []struct {
		name     string
		issue    *models.Issue
		expected models.DifficultyLevel
	}{
		{
			name: "Good first issue label",
			issue: &models.Issue{
				Labels: []models.Label{
					{Name: "good first issue"},
				},
				CommentsCount: 2,
			},
			expected: models.DifficultyBeginner,
		},
		{
			name: "First-timers-only label",
			issue: &models.Issue{
				Labels: []models.Label{
					{Name: "first-timers-only"},
				},
				CommentsCount: 0,
			},
			expected: models.DifficultyBeginner,
		},
		{
			name: "Documentation label",
			issue: &models.Issue{
				Labels: []models.Label{
					{Name: "documentation"},
				},
				CommentsCount: 1,
			},
			expected: models.DifficultyBeginner,
		},
		{
			name: "Typo in title",
			issue: &models.Issue{
				Title:         "Fix typo in README",
				Body:          "Small spelling mistake on line 42",
				CommentsCount: 0,
			},
			expected: models.DifficultyBeginner,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			level, reason := c.Classify(tc.issue)
			if level != tc.expected {
				t.Errorf("expected level %s, got %s (reason: %s)", tc.expected, level, reason)
			}
		})
	}
}

func TestClassifier_Advanced(t *testing.T) {
	c := classification.NewClassifier()

	tests := []struct {
		name     string
		issue    *models.Issue
		expected models.DifficultyLevel
	}{
		{
			name: "Security label",
			issue: &models.Issue{
				Labels: []models.Label{
					{Name: "security"},
				},
			},
			expected: models.DifficultyAdvanced,
		},
		{
			name: "Concurrency keyword in body",
			issue: &models.Issue{
				Title: "Goroutine leak under heavy load",
				Body:  "Detected race condition in channel dispatcher",
			},
			expected: models.DifficultyAdvanced,
		},
		{
			name: "High comment count discussion",
			issue: &models.Issue{
				Title:         "Refactor core interface",
				CommentsCount: 28,
			},
			expected: models.DifficultyAdvanced,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			level, reason := c.Classify(tc.issue)
			if level != tc.expected {
				t.Errorf("expected level %s, got %s (reason: %s)", tc.expected, level, reason)
			}
		})
	}
}

func TestClassifier_Intermediate_And_Fallback(t *testing.T) {
	c := classification.NewClassifier()

	issue := &models.Issue{
		Title: "Add endpoint for user preferences",
		Body:  "Implement new handler in router",
		Labels: []models.Label{
			{Name: "enhancement"},
		},
		CommentsCount: 3,
	}

	level, _ := c.Classify(issue)
	if level != models.DifficultyIntermediate {
		t.Errorf("expected intermediate, got %s", level)
	}

	nilLevel, _ := c.Classify(nil)
	if nilLevel != models.DifficultyUnknown {
		t.Errorf("expected unknown for nil issue, got %s", nilLevel)
	}
}
