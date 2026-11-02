package cli

import (
	"fmt"
	"os"

	"github.com/devradar/devradar/internal/output"
	"github.com/devradar/devradar/pkg/models"
	"github.com/spf13/cobra"
)

var (
	findLang       string
	findLabel      string
	findLevel      string
	findOrg        string
	findRepo       string
	findUnassigned bool
	findLimit      int
	findSort       string
	findOrder      string
	findMinStars   int
	findDetail     bool
)

var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Search and rank open-source issues",
	Long: `Search GitHub issues across projects with smart difficulty classification,
repository activity checks, and personalized skill matching.

Examples:
  devradar find --language go --level beginner
  devradar find --label "good first issue" --unassigned --limit 15
  devradar find --org kubernetes --level intermediate
  devradar find --repo gin-gonic/gin --unassigned
  devradar find --language python --format json
  devradar find --language rust --format csv > results.csv
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := &models.FilterOptions{
			Language:   findLang,
			Label:      findLabel,
			Level:      models.DifficultyLevel(findLevel),
			Org:        findOrg,
			Repo:       findRepo,
			Unassigned: findUnassigned,
			Limit:      findLimit,
			Sort:       findSort,
			Order:      findOrder,
			MinStars:   findMinStars,
		}

		issues, err := application.Discovery.Discover(cmd.Context(), opts)
		if err != nil {
			return fmt.Errorf("discovering issues: %w", err)
		}

		if findDetail || flagFormat == "detail" {
			for i := range issues {
				output.RenderIssueDetail(os.Stdout, &issues[i])
			}
			printRateLimitWarning()
			return nil
		}

		// Handle output format
		format := flagFormat
		switch format {
		case "json":
			if err := output.RenderJSON(os.Stdout, issues); err != nil {
				return err
			}
		case "csv":
			if err := output.RenderCSV(os.Stdout, issues); err != nil {
				return err
			}
		default:
			output.RenderTable(os.Stdout, issues)
		}

		printRateLimitWarning()
		return nil
	},
}

func init() {
	findCmd.Flags().StringVarP(&findLang, "language", "l", "", "Programming language (e.g. go, python, rust)")
	findCmd.Flags().StringVar(&findLabel, "label", "", "Filter by specific GitHub label (e.g. 'good first issue')")
	findCmd.Flags().StringVar(&findLevel, "level", "", "Difficulty level (beginner, intermediate, advanced)")
	findCmd.Flags().StringVar(&findOrg, "org", "", "Search inside a specific organization (e.g. kubernetes)")
	findCmd.Flags().StringVar(&findRepo, "repo", "", "Search inside a specific repository (e.g. gin-gonic/gin)")
	findCmd.Flags().BoolVarP(&findUnassigned, "unassigned", "u", false, "Only return issues with no assignee")
	findCmd.Flags().IntVar(&findLimit, "limit", 15, "Maximum number of issues to return")
	findCmd.Flags().StringVar(&findSort, "sort", "updated", "Sort field: updated, created, comments")
	findCmd.Flags().StringVar(&findOrder, "order", "desc", "Sort order: desc, asc")
	findCmd.Flags().IntVar(&findMinStars, "min-stars", 0, "Minimum number of repository stars")
	findCmd.Flags().BoolVarP(&findDetail, "detail", "d", false, "Display full details, descriptions and tips for all found issues")
}
