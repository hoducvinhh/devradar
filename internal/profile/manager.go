package profile

import (
	"strings"
	"time"

	"github.com/devradar/devradar/internal/storage"
	"github.com/devradar/devradar/pkg/models"
)

// Manager coordinates profile retrieval and persistence.
type Manager struct {
	store *storage.Storage
}

// NewManager creates a new Profile Manager instance.
func NewManager(store *storage.Storage) *Manager {
	return &Manager{store: store}
}

// Get returns current profile or default if none exists.
func (m *Manager) Get() (*models.Profile, error) {
	return m.store.LoadProfile()
}

// Update sets new profile values and saves them to local storage.
func (m *Manager) Update(level string, languages []string, interests []string, minStars int) (*models.Profile, error) {
	current, err := m.store.LoadProfile()
	if err != nil {
		current = models.DefaultProfile()
	}

	if level != "" {
		current.Level = models.DifficultyLevel(strings.ToLower(level))
	}

	if len(languages) > 0 {
		current.Languages = cleanSlice(languages)
	}

	if len(interests) > 0 {
		current.Interests = cleanSlice(interests)
	}

	if minStars >= 0 {
		current.MinStars = minStars
	}

	current.UpdatedAt = time.Now()
	if err := m.store.SaveProfile(current); err != nil {
		return nil, err
	}
	return current, nil
}

// Reset clears saved profile from local storage.
func (m *Manager) Reset() error {
	return m.store.ResetProfile()
}

// ParseCommaSeparated splits a comma-delimited string and trims whitespace.
func ParseCommaSeparated(input string) []string {
	if strings.TrimSpace(input) == "" {
		return []string{}
	}
	parts := strings.Split(input, ",")
	return cleanSlice(parts)
}

func cleanSlice(items []string) []string {
	var result []string
	seen := make(map[string]bool)
	for _, item := range items {
		trimmed := strings.ToLower(strings.TrimSpace(item))
		if trimmed != "" && !seen[trimmed] {
			seen[trimmed] = true
			result = append(result, trimmed)
		}
	}
	return result
}
