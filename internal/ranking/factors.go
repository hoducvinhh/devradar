package ranking

import (
	"fmt"
	"strings"
	"github.com/devradar/devradar/pkg/models"
)

func CalculateSkillScore(issue *models.Issue, profile *models.Profile, targetLang string) (float64, string) {
	issueLang := strings.ToLower(strings.TrimSpace(issue.RepoLanguage))
	target := strings.ToLower(strings.TrimSpace(targetLang))
	if target != "" && issueLang == target {
		return 100.0, fmt.Sprintf("Ngôn ngữ %s khớp chính xác với tiêu chí tìm kiếm", issue.RepoLanguage)
	}
	return 80.0, "Ngôn ngữ hợp lệ"
}
