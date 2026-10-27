package ranking

import (
	"fmt"
	"math"

	"github.com/devradar/devradar/pkg/models"
)

// Ranker calculates matching scores and explanations for issues.
type Ranker struct{}

// NewRanker creates a new Ranker instance.
func NewRanker() *Ranker {
	return &Ranker{}
}

// Rank computes the multi-factor matching score for an issue.
func (r *Ranker) Rank(issue *models.Issue, profile *models.Profile, targetLang string, targetLevel models.DifficultyLevel, repoHealth float64, repoHealthReason string) *models.RankingScore {
	// 1. Skill Score (30%)
	skillRaw, skillExp := CalculateSkillScore(issue, profile, targetLang)
	skillScore := skillRaw * models.WeightSkill

	// 2. Difficulty Score (25%)
	diffRaw, diffExp := CalculateDifficultyScore(issue.Difficulty, targetLevel)
	diffScore := diffRaw * models.WeightDifficulty

	// 3. Freshness Score (15%)
	freshRaw, freshExp := CalculateFreshnessScore(issue.UpdatedAt)
	freshScore := freshRaw * models.WeightFreshness

	// 4. Repository Activity Score (15%)
	repoScore := repoHealth * models.WeightRepoActivity
	repoExp := repoHealthReason
	if repoExp == "" {
		repoExp = fmt.Sprintf("Repository hoạt động ổn định (Điểm: %.0f/100)", repoHealth)
	}

	// 5. Issue Clarity Score (10%)
	clarityRaw, clarityExp := CalculateClarityScore(issue.Body)
	clarityScore := clarityRaw * models.WeightClarity

	// 6. Unassigned Score (5%)
	unassignedRaw, unassignedExp := CalculateUnassignedScore(issue.IsAssigned())
	unassignedScore := unassignedRaw * models.WeightUnassigned

	// Compute Total Score
	totalFloat := skillScore + diffScore + freshScore + repoScore + clarityScore + unassignedScore
	total := int(math.Round(totalFloat))
	if total > 100 {
		total = 100
	}
	if total < 0 {
		total = 0
	}

	tier := models.ComputeTier(total)

	explanations := []string{
		fmt.Sprintf("• [%02.0f/30] %s", math.Round(skillScore), skillExp),
		fmt.Sprintf("• [%02.0f/25] %s", math.Round(diffScore), diffExp),
		fmt.Sprintf("• [%02.0f/15] %s", math.Round(freshScore), freshExp),
		fmt.Sprintf("• [%02.0f/15] %s", math.Round(repoScore), repoExp),
		fmt.Sprintf("• [%02.0f/10] %s", math.Round(clarityScore), clarityExp),
		fmt.Sprintf("• [%02.0f/05] %s", math.Round(unassignedScore), unassignedExp),
	}

	rs := &models.RankingScore{
		TotalScore:        total,
		SkillScore:        skillScore,
		DifficultyScore:   diffScore,
		FreshnessScore:    freshScore,
		RepoActivityScore: repoScore,
		ClarityScore:      clarityScore,
		UnassignedScore:   unassignedScore,
		MatchTier:         tier,
		Explanations:      explanations,
	}

	issue.Ranking = rs
	return rs
}
