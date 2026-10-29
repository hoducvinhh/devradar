package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/devradar/devradar/internal/output"
	"github.com/devradar/devradar/internal/repository"
	"github.com/spf13/cobra"
)

var repoCmd = &cobra.Command{
	Use:   "repo <owner/repo>",
	Short: "Analyze repository health and maintenance status",
	Long: `Evaluate whether a repository is actively maintained and friendly for new contributors.

Examples:
  devradar repo gin-gonic/gin
  devradar repo kubernetes/kubernetes
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		parts := strings.Split(strings.TrimSpace(args[0]), "/")
		if len(parts) != 2 {
			return fmt.Errorf("invalid repository format: %s (expected owner/repo)", args[0])
		}
		owner, repo := parts[0], parts[1]

		r, err := application.Client.GetRepo(cmd.Context(), owner, repo)
		if err != nil {
			return fmt.Errorf("fetching repository info: %w", err)
		}

		analyzer := repository.NewAnalyzer()
		healthScore, healthReason := analyzer.Analyze(r)
		r.HealthScore = healthScore
		r.HealthReason = healthReason

		if flagFormat == "json" {
			return output.RenderJSON(os.Stdout, r)
		}

		titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#38bdf8"))
		labelStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#94a3b8"))
		valStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#f8fafc"))
		healthColor := lipgloss.Color("#22c55e")
		if healthScore < 50 {
			healthColor = lipgloss.Color("#ef4444")
		} else if healthScore < 75 {
			healthColor = lipgloss.Color("#f59e0b")
		}

		fmt.Println(titleStyle.Render(fmt.Sprintf("\n📦 Repository Health: %s", r.FullName)))
		if r.Description != "" {
			fmt.Printf("📝 %s\n\n", r.Description)
		}

		fmt.Printf("%s: %s\n", labelStyle.Render("Health Score"),
			lipgloss.NewStyle().Bold(true).Foreground(healthColor).Render(fmt.Sprintf("%.0f/100 (%s)", healthScore, healthReason)))
		fmt.Printf("%s: %s\n", labelStyle.Render("Primary Language"), valStyle.Render(r.Language))
		fmt.Printf("%s: %s | %s: %s | %s: %s\n",
			labelStyle.Render("Stars"), valStyle.Render(fmt.Sprintf("%d", r.StargazersCount)),
			labelStyle.Render("Forks"), valStyle.Render(fmt.Sprintf("%d", r.ForksCount)),
			labelStyle.Render("Open Issues"), valStyle.Render(fmt.Sprintf("%d", r.OpenIssuesCount)))
		fmt.Printf("%s: %s\n", labelStyle.Render("Last Push"), valStyle.Render(r.PushedAt.Format("2006-01-02 15:04:05")))

		if r.Archived {
			fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#ef4444")).Bold(true).Render("⚠️ WARNING: This repository is ARCHIVED (read-only). Do not submit new PRs."))
		}

		fmt.Printf("\nFind beginner tasks in this repo: 'devradar find --repo %s --level beginner'\n\n", r.FullName)
		printRateLimitWarning()
		return nil
	},
}
