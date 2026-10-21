package models

type DifficultyLevel string

const (
	DifficultyBeginner     DifficultyLevel = "beginner"
	DifficultyIntermediate DifficultyLevel = "intermediate"
	DifficultyAdvanced     DifficultyLevel = "advanced"
	DifficultyUnknown      DifficultyLevel = "unknown"
)

const (
	WeightSkill        = 0.30
	WeightDifficulty   = 0.25
	WeightFreshness    = 0.15
	WeightRepoActivity = 0.15
	WeightClarity      = 0.10
	WeightUnassigned   = 0.05
)

type RankingScore struct {
	TotalScore        int      `json:"total_score"`
	SkillScore        float64  `json:"skill_score"`
	DifficultyScore   float64  `json:"difficulty_score"`
	FreshnessScore    float64  `json:"freshness_score"`
	RepoActivityScore float64  `json:"repo_activity_score"`
	ClarityScore      float64  `json:"clarity_score"`
	UnassignedScore   float64  `json:"unassigned_score"`
	MatchTier         string   `json:"match_tier"`
	Explanations      []string `json:"explanations"`
}
