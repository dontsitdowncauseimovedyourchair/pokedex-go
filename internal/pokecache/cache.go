package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	data      []byte
}

type Cache struct {
	cache map[string]cacheEntry
	mu    sync.Mutex
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		data:      val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	out, ok := c.cache[key]
	return out.data, ok
}

// Deletes items stored in cache that were created before a certain time
func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, cacheItem := range c.cache {
			if !now.Before(cacheItem.createdAt.Add(interval)) {
				delete(c.cache, key)
			}
		}
		c.mu.Unlock()
	}
}

func NewCache(duration time.Duration) *Cache {
	out := Cache{
		cache: make(map[string]cacheEntry),
		mu:    sync.Mutex{},
	}
	go out.reapLoop(duration)
	return &out
}
