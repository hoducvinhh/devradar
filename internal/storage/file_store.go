package storage

import (
	"os"
	"path/filepath"
	"sync"
)

const (
	DirName       = ".devradar"
	ConfigFile    = "config.json"
	ProfileFile   = "profile.json"
	BookmarksFile = "bookmarks.json"
	CacheFile     = "cache.json"
)

type Storage struct {
	baseDir string
	mu      sync.RWMutex
}

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

func NewStorageWithDir(dir string) (*Storage, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}
	return &Storage{baseDir: dir}, nil
}

func (s *Storage) BaseDir() string {
	return s.baseDir
}
