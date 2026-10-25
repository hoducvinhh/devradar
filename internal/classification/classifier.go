package classification

import (
	"fmt"
	"strings"

	"github.com/devradar/devradar/pkg/models"
)

// Classifier estimates the difficulty of GitHub issues.
type Classifier struct{}

// NewClassifier creates a new Classifier instance.
func NewClassifier() *Classifier {
	return &Classifier{}
}

// Classify determines the difficulty level and rationale for an issue.
func (c *Classifier) Classify(issue *models.Issue) (models.DifficultyLevel, string) {
	if issue == nil {
		return models.DifficultyUnknown, "No issue data"
	}

	labels := make([]string, 0, len(issue.Labels))
	for _, l := range issue.Labels {
		labels = append(labels, strings.ToLower(strings.TrimSpace(l.Name)))
	}

	// 1. Check for explicit Advanced labels
	for _, adv := range AdvancedLabels {
		for _, l := range labels {
			if strings.Contains(l, adv) {
				return models.DifficultyAdvanced, fmt.Sprintf("Has advanced label '%s'", l)
			}
		}
	}

	// 2. Check for explicit Beginner labels
	for _, beg := range BeginnerLabels {
		for _, l := range labels {
			if strings.Contains(l, beg) {
				// Special check: if an issue has 'good first issue' but 25+ comments, it might be contested
				if issue.CommentsCount > 25 {
					return models.DifficultyIntermediate, fmt.Sprintf("Marked as '%s' but has %d comments (complex discussion)", l, issue.CommentsCount)
				}
				return models.DifficultyBeginner, fmt.Sprintf("Has beginner-friendly label '%s'", l)
			}
		}
	}

	// 3. Check for explicit Intermediate labels
	for _, inter := range IntermediateLabels {
		for _, l := range labels {
			if strings.Contains(l, inter) {
				if issue.CommentsCount > 20 {
					return models.DifficultyAdvanced, fmt.Sprintf("Has label '%s' with %d comments (high complexity)", l, issue.CommentsCount)
				}
				return models.DifficultyIntermediate, fmt.Sprintf("Has label '%s'", l)
			}
		}
	}

	// 4. Content and keyword heuristic
	textLower := strings.ToLower(issue.Title + " " + issue.Body)

	for _, kw := range AdvancedKeywords {
		if strings.Contains(textLower, kw) {
			return models.DifficultyAdvanced, fmt.Sprintf("Text mentions complex topic '%s'", kw)
		}
	}

	for _, kw := range BeginnerKeywords {
		if strings.Contains(textLower, kw) {
			return models.DifficultyBeginner, fmt.Sprintf("Text indicates beginner scope '%s'", kw)
		}
	}

	// 5. Heavy discussion threshold
	if issue.CommentsCount >= 15 {
		return models.DifficultyAdvanced, fmt.Sprintf("Extensive debate (%d comments)", issue.CommentsCount)
	}

	// 6. Default fallback
	return models.DifficultyIntermediate, "Standard issue complexity"
}
