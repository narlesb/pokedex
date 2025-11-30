package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	// Add fields for caching Pokémon data as needed
	mapData map[string]cacheEntry
	mux     *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		mapData: make(map[string]cacheEntry),
		mux:     &sync.Mutex{},
	}
	go c.reapLoop(interval)
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.mapData[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	entry, exists := c.mapData[key]
	if !exists {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for now := range ticker.C {
		c.reap(now, interval)
	}
}

func (c *Cache) reap(now time.Time, interval time.Duration) {
	c.mux.Lock()
	defer c.mux.Unlock()

	cutoff := now.Add(-interval)
	for key, entry := range c.mapData {
		if entry.createdAt.Before(cutoff) {
			delete(c.mapData, key)
		}
	}
}
