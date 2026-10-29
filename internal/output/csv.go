package output

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/devradar/devradar/pkg/models"
)

// RenderCSV exports issues into standard CSV format.
func RenderCSV(w io.Writer, issues []models.Issue) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	headers := []string{
		"ID",
		"Number",
		"Repository",
		"Title",
		"Difficulty",
		"MatchScore",
		"MatchTier",
		"State",
		"Comments",
		"Language",
		"Stars",
		"HTMLURL",
		"UpdatedAt",
	}

	if err := writer.Write(headers); err != nil {
		return err
	}

	for _, issue := range issues {
		score := "0"
		tier := ""
		if issue.Ranking != nil {
			score = strconv.Itoa(issue.Ranking.TotalScore)
			tier = issue.Ranking.MatchTier
		}

		updatedAtStr := ""
		if !issue.UpdatedAt.IsZero() {
			updatedAtStr = issue.UpdatedAt.Format(time.RFC3339)
		}

		row := []string{
			strconv.FormatInt(issue.ID, 10),
			strconv.Itoa(issue.Number),
			issue.RepoFullName,
			issue.Title,
			string(issue.Difficulty),
			score,
			tier,
			issue.State,
			strconv.Itoa(issue.CommentsCount),
			issue.RepoLanguage,
			strconv.Itoa(issue.RepoStars),
			issue.HTMLURL,
			updatedAtStr,
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("writing csv row: %w", err)
		}
	}

	return nil
}
