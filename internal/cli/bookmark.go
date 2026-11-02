package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/x/term"
	"github.com/devradar/devradar/internal/output"
	"github.com/devradar/devradar/pkg/models"
	"github.com/spf13/cobra"
)

var saveCmd = &cobra.Command{
	Use:   "save <issue-url or owner/repo#number>",
	Short: "Save an issue to your local bookmark list",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		owner, repo, number, err := parseIssueIdentifier(args[0])
		if err != nil {
			return err
		}

		issue, err := application.Client.GetIssue(cmd.Context(), owner, repo, number)
		if err != nil {
			return fmt.Errorf("fetching issue: %w", err)
		}

		bm := models.Bookmark{
			IssueID:      issue.ID,
			Number:       issue.Number,
			Title:        issue.Title,
			HTMLURL:      issue.HTMLURL,
			RepoFullName: fmt.Sprintf("%s/%s", owner, repo),
			SavedAt:      time.Now(),
		}

		if err := application.Storage.SaveBookmark(bm); err != nil {
			return fmt.Errorf("saving bookmark: %w", err)
		}

		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#22c55e")).Bold(true).Render(
			fmt.Sprintf("✔ Issue #%d (%s) saved to bookmarks! View anytime with 'devradar saved'.", bm.Number, bm.RepoFullName),
		))
		return nil
	},
}

var savedCmd = &cobra.Command{
	Use:   "saved",
	Short: "List all bookmarked issues",
	RunE: func(cmd *cobra.Command, args []string) error {
		bookmarks, err := application.Storage.LoadBookmarks()
		if err != nil {
			return fmt.Errorf("loading bookmarks: %w", err)
		}

		if flagFormat == "json" {
			return output.RenderJSON(os.Stdout, bookmarks)
		}

		if len(bookmarks) == 0 {
			fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#f59e0b")).Render(
				"No bookmarked issues yet. Use 'devradar save <issue-url>' to track tasks.",
			))
			return nil
		}

		headers := []string{"#", "REPO", "ISSUE", "SAVED AT", "URL"}
		var rows [][]string

		for i, bm := range bookmarks {
			rows = append(rows, []string{
				fmt.Sprintf("%d", i+1),
				bm.RepoFullName,
				fmt.Sprintf("#%d: %s", bm.Number, truncateString(bm.Title, 40)),
				bm.SavedAt.Format("2006-01-02"),
				bm.HTMLURL,
			})
		}

		termWidth := 90
		if w, _, err := term.GetSize(os.Stdout.Fd()); err == nil && w > 40 {
			termWidth = w - 2
			if termWidth > 120 {
				termWidth = 120
			}
		}

		t := table.New().
			Border(lipgloss.RoundedBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#6366f1"))).
			Headers(headers...).
			Rows(rows...).
			Width(termWidth).
			Wrap(false).
			StyleFunc(func(row, col int) lipgloss.Style {
				if row == table.HeaderRow {
					return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f8fafc")).Background(lipgloss.Color("#312e81")).Padding(0, 1)
				}
				st := lipgloss.NewStyle().Padding(0, 1)
				if col == 1 {
					return st.Bold(true).Foreground(lipgloss.Color("#06b6d4"))
				}
				return st
			})

		fmt.Println(t)
		fmt.Printf("\nTip: Remove a bookmark with 'devradar unsave <url-or-number>'\n")
		return nil
	},
}

var unsaveCmd = &cobra.Command{
	Use:   "unsave <issue-url or number>",
	Short: "Remove an issue from your local bookmarks",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		removed, err := application.Storage.RemoveBookmark(args[0])
		if err != nil {
			return fmt.Errorf("removing bookmark: %w", err)
		}

		if !removed {
			fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#ef4444")).Render("Bookmark not found matching: " + args[0]))
			return nil
		}

		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#22c55e")).Bold(true).Render("✔ Bookmark removed."))
		return nil
	},
}

func truncateString(s string, maxLen int) string {
	r := []rune(s)
	if len(r) <= maxLen {
		return s
	}
	return string(r[:maxLen-3]) + "..."
}
