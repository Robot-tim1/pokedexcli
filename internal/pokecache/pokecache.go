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
	newCache := &Cache{
		Mu:      &sync.Mutex{},
		Entries: make(map[string]CacheEntry),
	}
	go newCache.reapLoop(interval)
	return newCache
}

func (c *Cache) reapLoop(interval time.Duration) {
	timer := time.Tick(interval)
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

func (c *Cache) Delete(key string) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	delete(c.Entries, key)
}
