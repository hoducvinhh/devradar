package github

import (
	"fmt"
	"strings"

	"github.com/devradar/devradar/pkg/models"
)

func BuildSearchQuery(opts *models.FilterOptions) string {
	var parts []string
	parts = append(parts, "is:issue", "is:open")
	if opts == nil {
		return strings.Join(parts, " ")
	}
	if opts.Language != "" {
		parts = append(parts, fmt.Sprintf("language:%s", strings.TrimSpace(opts.Language)))
	}
	return strings.Join(parts, " ")
}
