package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/devradar/devradar/pkg/models"
)

type CacheItem struct {
	Repo      models.Repository `json:"repo"`
	ExpiresAt time.Time         `json:"expires_at"`
}

// Cache provides TTL-based caching for GitHub repository metadata.
type Cache struct {
	filePath string
	items    map[string]CacheItem
	mu       sync.RWMutex
	ttl      time.Duration
}

// NewCache initializes a cache instance with a specified TTL (e.g. 30 minutes).
func NewCache(baseDir string, ttl time.Duration) *Cache {
	c := &Cache{
		filePath: filepath.Join(baseDir, CacheFile),
		items:    make(map[string]CacheItem),
		ttl:      ttl,
	}
	c.load()
	return c
}

func (c *Cache) load() {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := os.ReadFile(c.filePath)
	if err != nil {
		return
	}
	var loaded map[string]CacheItem
	if err := json.Unmarshal(data, &loaded); err != nil {
		return
	}
	now := time.Now()
	for k, v := range loaded {
		if v.ExpiresAt.After(now) {
			c.items[k] = v
		}
	}
}

func (c *Cache) saveUnlocked() {
	data, err := json.Marshal(c.items)
	if err != nil {
		return
	}
	_ = os.WriteFile(c.filePath, data, 0600)
}

// GetRepo retrieves a cached repository if not expired.
func (c *Cache) GetRepo(fullName string) (*models.Repository, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[fullName]
	if !exists {
		return nil, false
	}
	if time.Now().After(item.ExpiresAt) {
		return nil, false
	}
	return &item.Repo, true
}

// SetRepo caches a repository for the configured TTL.
func (c *Cache) SetRepo(fullName string, repo models.Repository) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[fullName] = CacheItem{
		Repo:      repo,
		ExpiresAt: time.Now().Add(c.ttl),
	}
	c.saveUnlocked()
}

// Clear removes all cache entries.
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]CacheItem)
	_ = os.Remove(c.filePath)
}
