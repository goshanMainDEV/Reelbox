package cache

import (
	"sync"
	"time"

	"reelbox/internal/extractor"
)

type cacheEntry struct {
	result    *extractor.YtDlpResult
	expiresAt time.Time
}

type VideoCache struct {
	mu    sync.Mutex
	items map[string]cacheEntry
	ttl   time.Duration
}

func NewVideoCache(ttl time.Duration) *VideoCache {
	return &VideoCache{
		items: make(map[string]cacheEntry),
		ttl:   ttl,
	}
}

func (c *VideoCache) Get(id string) (*extractor.YtDlpResult, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.items[id]
	if !exists {
		return nil, false
	}

	if time.Now().After(entry.expiresAt) {
		delete(c.items, id)
		return nil, false
	}

	return entry.result, true
}

func (c *VideoCache) Set(id string, result *extractor.YtDlpResult) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[id] = cacheEntry{
		result:    result,
		expiresAt: time.Now().Add(c.ttl),
	}
}