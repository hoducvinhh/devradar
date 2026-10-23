package github

import (
	"net/http"
	"strconv"
	"sync"
	"time"
)

type RateLimit struct {
	Limit     int       `json:"limit"`
	Remaining int       `json:"remaining"`
	ResetTime time.Time `json:"reset_time"`
	Used      int       `json:"used"`
	Resource  string    `json:"resource"`
	mu        sync.RWMutex
}

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
