package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	mu      *sync.Mutex
	entries map[string]cacheEntry
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Time) *Cache {
	entries := make(map[string]cacheEntry)

	newCache := Cache{
		mu:      &sync.Mutex{},
		entries: entries,
	}
	return &newCache
}
