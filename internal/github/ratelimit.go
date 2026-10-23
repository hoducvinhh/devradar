package github

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// RateLimit holds current GitHub API rate limit status.
type RateLimit struct {
	Limit     int       `json:"limit"`
	Remaining int       `json:"remaining"`
	ResetTime time.Time `json:"reset_time"`
	Used      int       `json:"used"`
	Resource  string    `json:"resource"`
	mu        sync.RWMutex
}

// UpdateFromHeaders extracts rate limit information from GitHub HTTP response headers.
func (rl *RateLimit) UpdateFromHeaders(header http.Header) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if limitStr := header.Get("X-RateLimit-Limit"); limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil {
			rl.Limit = val
		}
	}
	if remStr := header.Get("X-RateLimit-Remaining"); remStr != "" {
		if val, err := strconv.Atoi(remStr); err == nil {
			rl.Remaining = val
		}
	}
	if usedStr := header.Get("X-RateLimit-Used"); usedStr != "" {
		if val, err := strconv.Atoi(usedStr); err == nil {
			rl.Used = val
		}
	}
	if resetStr := header.Get("X-RateLimit-Reset"); resetStr != "" {
		if sec, err := strconv.ParseInt(resetStr, 10, 64); err == nil {
			rl.ResetTime = time.Unix(sec, 0)
		}
	}
	if res := header.Get("X-RateLimit-Resource"); res != "" {
		rl.Resource = res
	}
}

// Status returns a copy of the current rate limit info.
func (rl *RateLimit) Status() (limit int, remaining int, resetTime time.Time) {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	return rl.Limit, rl.Remaining, rl.ResetTime
}

// IsExceeded returns true if quota is completely exhausted.
func (rl *RateLimit) IsExceeded() bool {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	if rl.Limit > 0 && rl.Remaining <= 0 {
		return time.Now().Before(rl.ResetTime)
	}
	return false
}

// WarningMessage returns a warning message if remaining quota is running low (<= 5).
func (rl *RateLimit) WarningMessage() string {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	if rl.Limit == 0 {
		return ""
	}
	if rl.Remaining <= 0 {
		waitMins := int(time.Until(rl.ResetTime).Minutes()) + 1
		return fmt.Sprintf("⚠️ GitHub API rate limit reached (%d/%d). Resets in %d minutes. Provide a token via 'devradar auth set-token' for 5,000 req/hr.", rl.Used, rl.Limit, waitMins)
	}
	if rl.Remaining <= 5 {
		return fmt.Sprintf("⚠️ GitHub API rate limit low: only %d requests remaining (resets at %s).", rl.Remaining, rl.ResetTime.Format("15:04:05"))
	}
	return ""
}
