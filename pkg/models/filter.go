package models

// FilterOptions contains all search query parameters for finding issues.
type FilterOptions struct {
	Language   string          `json:"language,omitempty"`
	Label      string          `json:"label,omitempty"`
	Level      DifficultyLevel `json:"level,omitempty"`
	Org        string          `json:"org,omitempty"`
	Repo       string          `json:"repo,omitempty"`
	Unassigned bool            `json:"unassigned,omitempty"`
	Limit      int             `json:"limit,omitempty"`
	Format     string          `json:"format,omitempty"` // "table", "json", "csv"
	Sort       string          `json:"sort,omitempty"`   // "updated", "created", "comments"
	Order      string          `json:"order,omitempty"`  // "desc", "asc"
	MinStars   int             `json:"min_stars,omitempty"`
	Keywords   string          `json:"keywords,omitempty"`
}

// DefaultFilterOptions returns sensible defaults for issue searches.
func DefaultFilterOptions() *FilterOptions {
	return &FilterOptions{
		Limit:  15,
		Format: "table",
		Sort:   "updated",
		Order:  "desc",
	}
}
