package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	Mu      *sync.Mutex
	Entries map[string]CacheEntry
}

type CacheEntry struct {
	CreatedAt time.Time
	Val       []byte
}

func NewCache(interval time.Duration) *Cache {
	entries := make(map[string]CacheEntry)

	newCache := Cache{
		Mu:      &sync.Mutex{},
		Entries: entries,
	}

	timer := time.Tick(interval)
	go newCache.reapLoop(timer, interval)

	return &newCache
}

func (c *Cache) reapLoop(timer <-chan time.Time, interval time.Duration) {
	for range timer {
		c.Mu.Lock()
		for key, entry := range c.Entries {
			if time.Since(entry.CreatedAt) >= interval {
				delete(c.Entries, key)
			}
		}
		c.Mu.Unlock()
	}
}

func (c *Cache) Add(key string, val []byte) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	c.Entries[key] = CacheEntry{CreatedAt: time.Now(), Val: val}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	if entry, ok := c.Entries[key]; ok {
		return entry.Val, true
	} else {
		return nil, false
	}
}
