package ranking

import (
	"fmt"
	"strings"
	"time"

	"github.com/devradar/devradar/pkg/models"
)

// CalculateSkillScore computes skill match score (0-100) and explanation.
func CalculateSkillScore(issue *models.Issue, profile *models.Profile, targetLang string) (float64, string) {
	issueLang := strings.ToLower(strings.TrimSpace(issue.RepoLanguage))
	target := strings.ToLower(strings.TrimSpace(targetLang))

	// Direct match from search filter
	if target != "" && issueLang == target {
		return 100.0, fmt.Sprintf("Ngôn ngữ %s khớp chính xác với tiêu chí tìm kiếm", issue.RepoLanguage)
	}

	if profile == nil || (len(profile.Languages) == 0 && len(profile.Interests) == 0) {
		if issueLang != "" {
			return 80.0, fmt.Sprintf("Ngôn ngữ %s (chế độ tìm kiếm mở)", issue.RepoLanguage)
		}
		return 70.0, "Không có ngôn ngữ chỉ định"
	}

	// Check profile languages
	for _, lang := range profile.Languages {
		if strings.EqualFold(lang, issueLang) {
			return 100.0, fmt.Sprintf("Ngôn ngữ %s khớp hoàn toàn với profile", issue.RepoLanguage)
		}
	}

	// Check interests in labels, title, or body
	textToScan := strings.ToLower(issue.Title + " " + issue.Body)
	for _, label := range issue.Labels {
		textToScan += " " + strings.ToLower(label.Name)
	}

	interestMatches := 0
	var matchedInterest string
	for _, interest := range profile.Interests {
		if strings.Contains(textToScan, strings.ToLower(interest)) {
			interestMatches++
			if matchedInterest == "" {
				matchedInterest = interest
			}
		}
	}

	if interestMatches > 0 {
		return 75.0, fmt.Sprintf("Khớp chủ đề quan tâm trong profile ('%s')", matchedInterest)
	}

	if len(profile.Languages) > 0 {
		return 30.0, fmt.Sprintf("Ngôn ngữ %s không nằm trong danh sách ưu tiên của profile", issue.RepoLanguage)
	}

	return 70.0, "Ngôn ngữ hợp lệ"
}

// CalculateDifficultyScore computes difficulty alignment score (0-100) and explanation.
func CalculateDifficultyScore(actual models.DifficultyLevel, desired models.DifficultyLevel) (float64, string) {
	if desired == "" || desired == models.DifficultyUnknown {
		return 80.0, fmt.Sprintf("Độ khó %s (không yêu cầu cấp độ cụ thể)", actual)
	}

	if actual == desired {
		return 100.0, fmt.Sprintf("Độ khó %s khớp với mong muốn", actual)
	}

	// Calculate distance between levels
	diffRank := map[models.DifficultyLevel]int{
		models.DifficultyBeginner:     1,
		models.DifficultyIntermediate: 2,
		models.DifficultyAdvanced:     3,
	}

	actualRank, ok1 := diffRank[actual]
	desiredRank, ok2 := diffRank[desired]

	if !ok1 || !ok2 {
		return 60.0, fmt.Sprintf("Độ khó ước lượng %s", actual)
	}

	dist := actualRank - desiredRank
	if dist < 0 {
		dist = -dist
	}

	if dist == 1 {
		return 50.0, fmt.Sprintf("Độ khó %s (chênh lệch 1 bậc so với mong muốn %s)", actual, desired)
	}

	return 15.0, fmt.Sprintf("Độ khó %s (chênh lệch xa so với mong muốn %s)", actual, desired)
}

// CalculateFreshnessScore computes freshness score (0-100) and explanation.
func CalculateFreshnessScore(updatedAt time.Time) (float64, string) {
	if updatedAt.IsZero() {
		return 50.0, "Thời gian cập nhật không xác định"
	}

	duration := time.Since(updatedAt)
	days := int(duration.Hours() / 24)
	hours := int(duration.Hours())

	var timeStr string
	if hours < 24 {
		timeStr = fmt.Sprintf("%d giờ trước", hours)
	} else {
		timeStr = fmt.Sprintf("%d ngày trước", days)
	}

	switch {
	case days <= 3:
		return 100.0, fmt.Sprintf("Cập nhật %s (Rất mới)", timeStr)
	case days <= 7:
		return 90.0, fmt.Sprintf("Cập nhật %s (Trong 1 tuần)", timeStr)
	case days <= 14:
		return 75.0, fmt.Sprintf("Cập nhật %s (Trong 2 tuần)", timeStr)
	case days <= 30:
		return 60.0, fmt.Sprintf("Cập nhật %s (Trong tháng)", timeStr)
	case days <= 90:
		return 40.0, fmt.Sprintf("Cập nhật %s", timeStr)
	default:
		return 20.0, fmt.Sprintf("Cập nhật khá lâu (%s)", timeStr)
	}
}

// CalculateClarityScore evaluates issue description quality (0-100) and explanation.
func CalculateClarityScore(body string) (float64, string) {
	if strings.TrimSpace(body) == "" {
		return 20.0, "Mô tả issue rất ngắn hoặc để trống"
	}

	score := 20.0
	var features []string

	// Check code blocks
	if strings.Contains(body, "```") {
		score += 30.0
		features = append(features, "có code block")
	}

	// Check task checklist
	if strings.Contains(body, "- [ ]") || strings.Contains(body, "- [x]") {
		score += 30.0
		features = append(features, "có checklist nhiệm vụ")
	}

	// Check length and headers
	if len(body) >= 200 {
		score += 10.0
	}
	if strings.Contains(body, "##") || strings.Contains(body, "**") {
		score += 10.0
		features = append(features, "trình bày có cấu trúc")
	}

	if score > 100.0 {
		score = 100.0
	}

	if len(features) > 0 {
		return score, fmt.Sprintf("Mô tả chi tiết (%s)", strings.Join(features, ", "))
	}
	return score, "Mô tả issue ở mức cơ bản"
}

// CalculateUnassignedScore checks if issue is free for taking (0-100) and explanation.
func CalculateUnassignedScore(isAssigned bool) (float64, string) {
	if !isAssigned {
		return 100.0, "Issue đang để trống assignee (Sẵn sàng tiếp nhận)"
	}
	return 0.0, "Issue đã có người được gán xử lý"
}
