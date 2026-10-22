package storage

import (
	"path/filepath"
	"sync"
	"time"

	"github.com/devradar/devradar/pkg/models"
)

type CacheItem struct {
	Repo      models.Repository `json:"repo"`
	ExpiresAt time.Time         `json:"expires_at"`
}

type Cache struct {
	filePath string
	items    map[string]CacheItem
	mu       sync.RWMutex
	ttl      time.Duration
}

func NewCache(baseDir string, ttl time.Duration) *Cache {
	return &Cache{
		filePath: filepath.Join(baseDir, CacheFile),
		items:    make(map[string]CacheItem),
		ttl:      ttl,
	}
}
