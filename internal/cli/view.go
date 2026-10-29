package cli

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/devradar/devradar/internal/classification"
	"github.com/devradar/devradar/internal/output"
	"github.com/devradar/devradar/internal/ranking"
	"github.com/devradar/devradar/internal/repository"
	"github.com/spf13/cobra"
)

var viewCmd = &cobra.Command{
	Use:   "view <issue-url or owner/repo#number>",
	Short: "View detailed information and match breakdown for an issue",
	Long: `Display rich details, full description, difficulty breakdown, and contribution advice for a specific issue.

Examples:
  devradar view https://github.com/gin-gonic/gin/issues/3421
  devradar view gin-gonic/gin#3421
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		owner, repo, number, err := parseIssueIdentifier(args[0])
		if err != nil {
			return err
		}

		issue, err := application.Client.GetIssue(cmd.Context(), owner, repo, number)
		if err != nil {
			return fmt.Errorf("fetching issue: %w", err)
		}

		// Fetch repository details for health analysis
		r, _ := application.Client.GetRepo(cmd.Context(), owner, repo)
		var health float64 = 70.0
		var healthReason = "Repository metadata unavailable"
		if r != nil {
			issue.RepoStars = r.StargazersCount
			issue.RepoForks = r.ForksCount
			issue.RepoLanguage = r.Language
			issue.RepoPushedAt = r.PushedAt
			issue.RepoArchived = r.Archived

			analyzer := repository.NewAnalyzer()
			health, healthReason = analyzer.Analyze(r)
		}

		// Difficulty classification
		classifier := classification.NewClassifier()
		diff, diffReason := classifier.Classify(issue)
		issue.Difficulty = diff
		issue.DifficultyReason = diffReason

		// Rank against user profile
		prof, _ := application.Profile.Get()
		ranker := ranking.NewRanker()
		ranker.Rank(issue, prof, "", "", health, healthReason)

		if flagFormat == "json" {
			return output.RenderJSON(os.Stdout, issue)
		}

		output.RenderIssueDetail(os.Stdout, issue)
		printRateLimitWarning()
		return nil
	},
}

func parseIssueIdentifier(input string) (string, string, int, error) {
	input = strings.TrimSpace(input)

	// Case 1: Full URL https://github.com/owner/repo/issues/123
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		u, err := url.Parse(input)
		if err != nil {
			return "", "", 0, fmt.Errorf("invalid url: %w", err)
		}
		parts := strings.Split(strings.Trim(u.Path, "/"), "/")
		if len(parts) >= 4 && parts[2] == "issues" {
			num, err := strconv.Atoi(parts[3])
			if err != nil {
				return "", "", 0, fmt.Errorf("invalid issue number: %w", err)
			}
			return parts[0], parts[1], num, nil
		}
	}

	// Case 2: owner/repo#number
	if strings.Contains(input, "#") {
		parts := strings.Split(input, "#")
		if len(parts) == 2 {
			repoParts := strings.Split(parts[0], "/")
			if len(repoParts) == 2 {
				num, err := strconv.Atoi(parts[1])
				if err != nil {
					return "", "", 0, fmt.Errorf("invalid issue number: %w", err)
				}
				return repoParts[0], repoParts[1], num, nil
			}
		}
	}

	return "", "", 0, errors.New("invalid issue identifier. Format must be https://github.com/owner/repo/issues/123 or owner/repo#123")
}
