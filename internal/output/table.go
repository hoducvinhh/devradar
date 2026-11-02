package output

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/x/term"
	"github.com/devradar/devradar/pkg/models"
)

var (
	// Styling palettes
	primaryColor = lipgloss.Color("#6366f1") // Indigo
	successColor = lipgloss.Color("#22c55e") // Emerald Green
	warningColor = lipgloss.Color("#f59e0b") // Amber Yellow
	dangerColor  = lipgloss.Color("#ef4444") // Coral Red
	infoColor    = lipgloss.Color("#06b6d4") // Cyan
	mutedColor   = lipgloss.Color("#94a3b8") // Slate Muted
	headerColor  = lipgloss.Color("#f8fafc") // Bright White

	badgeBeginner = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#14532d")).
			Background(lipgloss.Color("#86efac")).
			Render("BEGINNER")

	badgeIntermediate = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#78350f")).
				Background(lipgloss.Color("#fde68a")).
				Render("INTERMEDIATE")

	badgeAdvanced = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#701a75")).
			Background(lipgloss.Color("#f5d0fe")).
			Render("ADVANCED")

	badgeUnknown = lipgloss.NewStyle().
			Foreground(mutedColor).
			Render("UNKNOWN")
)

// getTerminalWidth returns current terminal columns or a safe default.
func getTerminalWidth() int {
	if w, _, err := term.GetSize(os.Stdout.Fd()); err == nil && w > 40 {
		return w
	}
	if w, _, err := term.GetSize(os.Stderr.Fd()); err == nil && w > 40 {
		return w
	}
	if w, _, err := term.GetSize(os.Stdin.Fd()); err == nil && w > 40 {
		return w
	}
	return 90
}

// RenderTable renders a styled Lip Gloss table for the list of issues with responsive sizing.
func RenderTable(w io.Writer, issues []models.Issue) {
	if len(issues) == 0 {
		emptyMsg := lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true).
			Render("No matching issues found. Try broadening your search filters.")
		fmt.Fprintln(w, emptyMsg)
		return
	}

	termWidth := getTerminalWidth() - 2
	if termWidth > 130 {
		termWidth = 130
	} else if termWidth < 75 {
		termWidth = 75
	}

	showComments := termWidth >= 115
	showUpdated := termWidth >= 85

	var headers []string
	if showComments {
		headers = []string{"#", "SCORE", "LEVEL", "REPOSITORY", "ISSUE TITLE", "UPDATED", "COMMENTS"}
	} else if showUpdated {
		headers = []string{"#", "SCORE", "LEVEL", "REPOSITORY", "ISSUE TITLE", "UPDATED"}
	} else {
		headers = []string{"#", "SCORE", "LEVEL", "REPOSITORY", "ISSUE TITLE"}
	}

	maxRepoLen := 22
	maxTitleLen := 34
	if termWidth >= 115 {
		maxRepoLen = 28
		maxTitleLen = 46
	} else if termWidth >= 95 {
		maxRepoLen = 24
		maxTitleLen = 38
	}

	var rows [][]string

	for idx, issue := range issues {
		numStr := fmt.Sprintf("%d", idx+1)

		// Score & Tier
		scoreStr := "-"
		if issue.Ranking != nil {
			scoreColor := successColor
			if issue.Ranking.TotalScore < 50 {
				scoreColor = dangerColor
			} else if issue.Ranking.TotalScore < 75 {
				scoreColor = warningColor
			}

			scoreStyled := lipgloss.NewStyle().Bold(true).Foreground(scoreColor).Render(fmt.Sprintf("%d", issue.Ranking.TotalScore))
			scoreStr = fmt.Sprintf("%s/100", scoreStyled)
		}

		// Level Badge
		levelBadge := formatLevelBadge(issue.Difficulty)

		// Repo with stars (truncated safely to prevent overflow)
		repoDisplay := issue.RepoFullName
		if issue.RepoStars > 0 {
			repoDisplay = fmt.Sprintf("%s ★ %s", issue.RepoFullName, formatNumber(issue.RepoStars))
		}
		repoDisplay = truncate(repoDisplay, maxRepoLen)

		// Title truncated dynamically based on screen width
		title := truncate(issue.Title, maxTitleLen)

		// Updated time
		updated := formatRelativeTime(issue.UpdatedAt)

		// Comments
		comments := fmt.Sprintf("%d", issue.CommentsCount)

		var row []string
		if showComments {
			row = []string{numStr, scoreStr, levelBadge, repoDisplay, title, updated, comments}
		} else if showUpdated {
			row = []string{numStr, scoreStr, levelBadge, repoDisplay, title, updated}
		} else {
			row = []string{numStr, scoreStr, levelBadge, repoDisplay, title}
		}

		rows = append(rows, row)
	}

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(primaryColor)).
		Headers(headers...).
		Rows(rows...).
		Width(termWidth).
		Wrap(false).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return lipgloss.NewStyle().
					Bold(true).
					Foreground(headerColor).
					Background(lipgloss.Color("#312e81")).
					Padding(0, 1)
			}

			st := lipgloss.NewStyle().Padding(0, 1)
			switch col {
			case 0:
				return st.Foreground(mutedColor).Align(lipgloss.Right)
			case 1:
				return st.Align(lipgloss.Center)
			case 2:
				return st.Align(lipgloss.Center)
			case 3:
				return st.Bold(true).Foreground(infoColor)
			case 4:
				return st.Foreground(lipgloss.Color("#e2e8f0"))
			case 5:
				return st.Foreground(mutedColor)
			case 6:
				return st.Foreground(mutedColor).Align(lipgloss.Right)
			default:
				return st
			}
		})

	fmt.Fprintln(w, t)
	fmt.Fprintf(w, "\n%s Tip: Run '%s' to view detailed match breakdown and instructions.\n",
		lipgloss.NewStyle().Foreground(infoColor).Bold(true).Render("💡"),
		lipgloss.NewStyle().Foreground(warningColor).Render("devradar view <issue-url>"))
}

func formatLevelBadge(level models.DifficultyLevel) string {
	switch level {
	case models.DifficultyBeginner:
		return badgeBeginner
	case models.DifficultyIntermediate:
		return badgeIntermediate
	case models.DifficultyAdvanced:
		return badgeAdvanced
	default:
		return badgeUnknown
	}
}

func formatNumber(n int) string {
	if n >= 1000 {
		return fmt.Sprintf("%.1fk", float64(n)/1000.0)
	}
	return fmt.Sprintf("%d", n)
}

func formatRelativeTime(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	case d < 30*24*time.Hour:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	case d < 365*24*time.Hour:
		return fmt.Sprintf("%dmo ago", int(d.Hours()/(24*30)))
	default:
		return fmt.Sprintf("%dy ago", int(d.Hours()/(24*365)))
	}
}

func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-3]) + "..."
}
