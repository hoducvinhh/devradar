package models

// DifficultyLevel indicates the estimated issue difficulty.
type DifficultyLevel string

const (
	DifficultyBeginner     DifficultyLevel = "beginner"
	DifficultyIntermediate DifficultyLevel = "intermediate"
	DifficultyAdvanced     DifficultyLevel = "advanced"
	DifficultyUnknown      DifficultyLevel = "unknown"
)

// FactorWeights constants defining the weight of each ranking dimension.
const (
	WeightSkill        = 0.30
	WeightDifficulty   = 0.25
	WeightFreshness    = 0.15
	WeightRepoActivity = 0.15
	WeightClarity      = 0.10
	WeightUnassigned   = 0.05
)

// RankingScore holds the detailed factor scores and weighted total for an issue.
type RankingScore struct {
	TotalScore       int      `json:"total_score"`       // Normalized 0 to 100
	SkillScore       float64  `json:"skill_score"`       // Max 30 pts (raw 0-100 * 0.30)
	DifficultyScore  float64  `json:"difficulty_score"`  // Max 25 pts (raw 0-100 * 0.25)
	FreshnessScore   float64  `json:"freshness_score"`   // Max 15 pts (raw 0-100 * 0.15)
	RepoActivityScore float64 `json:"repo_activity_score"` // Max 15 pts (raw 0-100 * 0.15)
	ClarityScore     float64  `json:"clarity_score"`     // Max 10 pts (raw 0-100 * 0.10)
	UnassignedScore  float64  `json:"unassigned_score"`  // Max 5 pts (raw 0-100 * 0.05)
	MatchTier        string   `json:"match_tier"`        // Excellent, Good, Fair, Low
	Explanations     []string `json:"explanations"`      // Detailed explanation lines
}

// ComputeTier returns a human-friendly match tier based on total score.
func ComputeTier(score int) string {
	switch {
	case score >= 85:
		return "Excellent Match"
	case score >= 70:
		return "Good Match"
	case score >= 50:
		return "Fair Match"
	default:
		return "Low Match"
	}
}
