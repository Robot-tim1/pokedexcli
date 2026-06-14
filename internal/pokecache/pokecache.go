package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cacheMu *sync.Mutex
	Entries map[string]CacheEntry
}

type CacheEntry struct {
	CreatedAt time.Time
	Val       []byte
}

func NewCache(interval time.Duration) Cache {
	newCache := Cache{
		cacheMu: &sync.Mutex{},
		Entries: make(map[string]CacheEntry),
	}
	go newCache.reapLoop(interval)
	return newCache
}

func (c *Cache) reapLoop(interval time.Duration) {
	timer := time.Tick(interval)
	for range timer {
		c.cacheMu.Lock()
		for key, entry := range c.Entries {
			if time.Since(entry.CreatedAt) >= interval {
				delete(c.Entries, key)
			}
		}
		c.cacheMu.Unlock()
	}
}

func (c *Cache) Add(key string, val []byte) {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()
	c.Entries[key] = CacheEntry{CreatedAt: time.Now(), Val: val}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()
	entry, ok := c.Entries[key]
	return entry.Val, ok
}

func (c *Cache) Delete(key string) {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()
	delete(c.Entries, key)
}
