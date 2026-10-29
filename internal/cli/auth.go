package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage GitHub API authentication",
	Long:  `Configure your personal access token to expand rate limits from 60 to 5,000 requests per hour.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAuthStatus(cmd, args)
	},
}

var authSetTokenCmd = &cobra.Command{
	Use:   "set-token <token>",
	Short: "Save GitHub Personal Access Token",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		token := strings.TrimSpace(args[0])
		if token == "" {
			return fmt.Errorf("token cannot be empty")
		}

		cfg := application.Config
		cfg.GitHubToken = token
		if err := application.Storage.SaveConfig(cfg); err != nil {
			return fmt.Errorf("saving token: %w", err)
		}

		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#22c55e")).Bold(true).Render(
			"✔ GitHub token saved securely to ~/.devradar/config.json!",
		))
		return nil
	},
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check GitHub authentication and rate limit status",
	RunE:  runAuthStatus,
}

func runAuthStatus(cmd *cobra.Command, args []string) error {
	cfg := application.Config
	hasToken := application.Client.HasToken()

	labelStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#94a3b8"))
	okStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#22c55e"))
	warnStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f59e0b"))

	fmt.Println(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#6366f1")).Render("\n🔐 GitHub Authentication Status"))

	if hasToken {
		masked := "••••••••"
		if len(cfg.GitHubToken) > 8 {
			masked = cfg.GitHubToken[:4] + "••••" + cfg.GitHubToken[len(cfg.GitHubToken)-4:]
		}
		fmt.Printf("%s: %s (%s)\n", labelStyle.Render("Status"), okStyle.Render("Authenticated (Token Active)"), masked)
		fmt.Printf("%s: %s\n", labelStyle.Render("Rate Limit Quota"), "Up to 5,000 requests/hour")
	} else {
		fmt.Printf("%s: %s\n", labelStyle.Render("Status"), warnStyle.Render("Unauthenticated (Anonymous)"))
		fmt.Printf("%s: %s\n", labelStyle.Render("Rate Limit Quota"), "60 requests/hour")
		fmt.Printf("\nTip: Run '%s' to unlock 5,000 requests/hour.\n",
			lipgloss.NewStyle().Foreground(lipgloss.Color("#38bdf8")).Render("devradar auth set-token <ghp_your_token>"))
	}

	limit, remaining, resetTime := application.Client.RateLimit().Status()
	if limit > 0 {
		fmt.Printf("%s: %d / %d (Resets at %s)\n",
			labelStyle.Render("Current Rate Limit"),
			remaining,
			limit,
			resetTime.Format("15:04:05"))
	}
	fmt.Println()

	return nil
}

var authClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Remove saved GitHub token",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := application.Config
		cfg.GitHubToken = ""
		if err := application.Storage.SaveConfig(cfg); err != nil {
			return fmt.Errorf("clearing token: %w", err)
		}
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#22c55e")).Bold(true).Render("✔ GitHub token removed."))
		return nil
	},
}

func init() {
	authCmd.AddCommand(authSetTokenCmd)
	authCmd.AddCommand(authStatusCmd)
	authCmd.AddCommand(authClearCmd)
}
