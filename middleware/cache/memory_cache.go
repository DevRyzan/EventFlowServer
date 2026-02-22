package cache

import (
	"log/slog"
	"sync"
	"time"
)

type cacheEntry struct {
	value    interface{}
	expireAt time.Time
}

type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]cacheEntry
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		items: make(map[string]cacheEntry),
	}
}

// Get returns the value for key if it exists and has not expired.
func (c *MemoryCache) Get(key string) (interface{}, bool) {
	// mutex pattern to protect the map from concurrent access
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.items[key]
	if !ok || time.Now().After(entry.expireAt) {
		// if the key is not found or the value has expired, return false
		return nil, false
	}
	slog.Debug("cache hit", "key", key)
	return entry.value, true
}

// Set stores a value with the given TTL.
func (c *MemoryCache) Set(key string, value interface{}, ttl time.Duration) {
	// mutex pattern to protect the map from concurrent access
	c.mu.Lock()
	defer c.mu.Unlock()

	// set the value with the given TTL
	c.items[key] = cacheEntry{
		value:    value,
		expireAt: time.Now().Add(ttl),
	}
	slog.Debug("cache set", "key", key)
}
