package repository

import (
	"fmt"
	"time"

	"github.com/devradar/devradar/pkg/models"
)

// Analyzer inspects repository activity and computes a health score.
type Analyzer struct{}

// NewAnalyzer creates a new Analyzer instance.
func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

// Analyze evaluates the repository's activity, maintenance status, and community vitality.
// Returns health score (0-100) and rationale.
func (a *Analyzer) Analyze(repo *models.Repository) (float64, string) {
	if repo == nil {
		return 50.0, "Repository metadata unavailable"
	}

	// 1. If archived or disabled, repo is inactive/read-only
	if repo.Archived || repo.Disabled {
		return 0.0, "Repository is archived or disabled (read-only)"
	}

	// 2. Commit / Push freshness
	pushedAt := repo.PushedAt
	if pushedAt.IsZero() {
		pushedAt = repo.UpdatedAt
	}

	daysSincePush := int(time.Since(pushedAt).Hours() / 24)
	if daysSincePush < 0 {
		daysSincePush = 0
	}

	var score float64
	var reason string

	switch {
	case daysSincePush <= 14:
		score = 100.0
		reason = fmt.Sprintf("High activity (pushed %d days ago)", daysSincePush)
	case daysSincePush <= 30:
		score = 90.0
		reason = fmt.Sprintf("Active (pushed %d days ago)", daysSincePush)
	case daysSincePush <= 60:
		score = 75.0
		reason = fmt.Sprintf("Moderate activity (pushed %d days ago)", daysSincePush)
	case daysSincePush <= 90:
		score = 60.0
		reason = fmt.Sprintf("Low recent activity (pushed %d days ago)", daysSincePush)
	case daysSincePush <= 180:
		score = 40.0
		reason = fmt.Sprintf("Infrequent updates (pushed %d days ago)", daysSincePush)
	default:
		score = 10.0
		reason = fmt.Sprintf("Stale repository (no push for %d days)", daysSincePush)
	}

	// Bonus / penalty adjustments
	if repo.StargazersCount > 500 {
		// Popular repo with established maintainers
		score = min(100.0, score+5.0)
	}

	return score, reason
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
