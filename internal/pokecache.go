package internal

import (
	"sync"
	"time"
)

type Cache struct {
	cacheEntry map[string]CacheEntry
	mu         sync.Mutex
}

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	newInstance := Cache{}
	newInstance.cacheEntry = map[string]CacheEntry{}

	go newInstance.ReapLoop(interval)
	return &newInstance
}

func (c *Cache) Add(key string, val []byte) {
	newEntry := CacheEntry{createdAt: time.Now(), val: val}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.cacheEntry[key] = newEntry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.cacheEntry[key]
	if exists {
		return entry.val, true
	} else {
		return nil, false
	}
}

func (c *Cache) ReapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.cacheEntry {
			if now.Sub(entry.createdAt) > interval {
				delete(c.cacheEntry, key)
			}
		}
		c.mu.Unlock()
	}
}
