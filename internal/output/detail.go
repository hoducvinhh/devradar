package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/devradar/devradar/pkg/models"
)

// RenderIssueDetail prints a rich, formatted view of a single issue, including match explanations and instructions.
func RenderIssueDetail(w io.Writer, issue *models.Issue) {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#38bdf8"))
	repoStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#818cf8"))
	sectionStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f1f5f9")).Underline(true)

	fmt.Fprintf(w, "\n%s #%d: %s\n", repoStyle.Render(issue.RepoFullName), issue.Number, titleStyle.Render(issue.Title))
	fmt.Fprintf(w, "🔗 URL: %s\n", lipgloss.NewStyle().Foreground(lipgloss.Color("#60a5fa")).Render(issue.HTMLURL))
	fmt.Fprintf(w, "🏷️  Labels: %s\n", formatLabels(issue.Labels))
	fmt.Fprintf(w, "📊 Difficulty: %s (%s)\n", formatLevelBadge(issue.Difficulty), issue.DifficultyReason)

	if issue.Ranking != nil {
		fmt.Fprintf(w, "\n%s\n", sectionStyle.Render("Match Breakdown & Score"))
		fmt.Fprintf(w, "Score: %s [%s]\n",
			lipgloss.NewStyle().Bold(true).Foreground(successColor).Render(fmt.Sprintf("%d/100", issue.Ranking.TotalScore)),
			lipgloss.NewStyle().Bold(true).Render(issue.Ranking.MatchTier))

		for _, exp := range issue.Ranking.Explanations {
			fmt.Fprintf(w, "  %s\n", exp)
		}
	}

	fmt.Fprintf(w, "\n%s\n", sectionStyle.Render("Description Preview"))
	body := issue.Body
	if len(body) > 600 {
		body = body[:600] + "\n... (truncated, open URL to read full description)"
	}
	fmt.Fprintf(w, "%s\n\n", lipgloss.NewStyle().Foreground(lipgloss.Color("#cbd5e1")).Render(body))

	fmt.Fprintf(w, "%s Next steps for first-time contribution:\n", lipgloss.NewStyle().Foreground(infoColor).Bold(true).Render("🚀"))
	fmt.Fprintf(w, "  1. Review the issue discussion on GitHub\n")
	fmt.Fprintf(w, "  2. Comment on the issue asking to be assigned (e.g. 'Hi! I'd like to work on this issue.')\n")
	fmt.Fprintf(w, "  3. Save this issue locally: devradar save %s\n\n", issue.HTMLURL)
}

func formatLabels(labels []models.Label) string {
	if len(labels) == 0 {
		return "None"
	}
	var rendered []string
	for _, l := range labels {
		rendered = append(rendered, lipgloss.NewStyle().
			Foreground(lipgloss.Color("#334155")).
			Background(lipgloss.Color("#cbd5e1")).
			Padding(0, 1).
			Render(l.Name))
	}
	return strings.Join(rendered, " ")
}
