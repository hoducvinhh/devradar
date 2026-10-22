package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/devradar/devradar/pkg/models"
)

const (
	DirName       = ".devradar"
	ConfigFile    = "config.json"
	ProfileFile   = "profile.json"
	BookmarksFile = "bookmarks.json"
	CacheFile     = "cache.json"
)

// Config holds user configuration, including authentication token.
type Config struct {
	GitHubToken string    `json:"github_token,omitempty"`
	DefaultLang string    `json:"default_lang,omitempty"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Storage handles local persistence in ~/.devradar/
type Storage struct {
	baseDir string
	mu      sync.RWMutex
}

// NewStorage initializes storage in the user's home directory.
func NewStorage() (*Storage, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	baseDir := filepath.Join(home, DirName)
	if err := os.MkdirAll(baseDir, 0700); err != nil {
		return nil, err
	}
	return &Storage{baseDir: baseDir}, nil
}

// NewStorageWithDir allows specifying a custom directory (useful for unit tests).
func NewStorageWithDir(dir string) (*Storage, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}
	return &Storage{baseDir: dir}, nil
}

// BaseDir returns the root storage directory.
func (s *Storage) BaseDir() string {
	return s.baseDir
}

// --- Config Management ---

func (s *Storage) configPath() string {
	return filepath.Join(s.baseDir, ConfigFile)
}

func (s *Storage) LoadConfig() (*Config, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := os.ReadFile(s.configPath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Config{}, nil
		}
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (s *Storage) SaveConfig(cfg *Config) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	cfg.UpdatedAt = time.Now()
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.configPath(), data, 0600)
}

// --- Profile Management ---

func (s *Storage) profilePath() string {
	return filepath.Join(s.baseDir, ProfileFile)
}

func (s *Storage) LoadProfile() (*models.Profile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := os.ReadFile(s.profilePath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return models.DefaultProfile(), nil
		}
		return nil, err
	}
	var p models.Profile
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *Storage) SaveProfile(p *models.Profile) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	p.UpdatedAt = time.Now()
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.profilePath(), data, 0600)
}

func (s *Storage) ResetProfile() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := os.Remove(s.profilePath())
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

// --- Bookmarks Management ---

func (s *Storage) bookmarksPath() string {
	return filepath.Join(s.baseDir, BookmarksFile)
}

func (s *Storage) LoadBookmarks() ([]models.Bookmark, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := os.ReadFile(s.bookmarksPath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []models.Bookmark{}, nil
		}
		return nil, err
	}
	var bookmarks []models.Bookmark
	if err := json.Unmarshal(data, &bookmarks); err != nil {
		return nil, err
	}
	return bookmarks, nil
}

func (s *Storage) SaveBookmark(bm models.Bookmark) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	bookmarks, err := s.loadBookmarksUnlocked()
	if err != nil {
		bookmarks = []models.Bookmark{}
	}

	// Update existing or append new
	exists := false
	for i, b := range bookmarks {
		if (bm.IssueID > 0 && b.IssueID == bm.IssueID) || (bm.HTMLURL != "" && b.HTMLURL == bm.HTMLURL) {
			bookmarks[i] = bm
			exists = true
			break
		}
	}
	if !exists {
		if bm.SavedAt.IsZero() {
			bm.SavedAt = time.Now()
		}
		bookmarks = append(bookmarks, bm)
	}

	return s.saveBookmarksUnlocked(bookmarks)
}

func (s *Storage) RemoveBookmark(identifier string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	bookmarks, err := s.loadBookmarksUnlocked()
	if err != nil {
		return false, err
	}

	var updated []models.Bookmark
	removed := false
	for _, b := range bookmarks {
		// match by url, issue number, or repo/number
		if b.HTMLURL == identifier || stringMatchIssue(b, identifier) {
			removed = true
			continue
		}
		updated = append(updated, b)
	}

	if removed {
		if err := s.saveBookmarksUnlocked(updated); err != nil {
			return false, err
		}
	}
	return removed, nil
}

func (s *Storage) loadBookmarksUnlocked() ([]models.Bookmark, error) {
	data, err := os.ReadFile(s.bookmarksPath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []models.Bookmark{}, nil
		}
		return nil, err
	}
	var bookmarks []models.Bookmark
	if err := json.Unmarshal(data, &bookmarks); err != nil {
		return nil, err
	}
	return bookmarks, nil
}

func (s *Storage) saveBookmarksUnlocked(bookmarks []models.Bookmark) error {
	data, err := json.MarshalIndent(bookmarks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.bookmarksPath(), data, 0600)
}

func stringMatchIssue(b models.Bookmark, id string) bool {
	// match exact issue number or html_url suffix
	return b.HTMLURL == id || (b.RepoFullName != "" && b.HTMLURL != "" && (b.HTMLURL == id || filepath.Base(b.HTMLURL) == id))
}
