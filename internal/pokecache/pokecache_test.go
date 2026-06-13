package pokecache

import (
	"bytes"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}

func TestConcurrencySafety(t *testing.T) {
	const interval = 5 * time.Second
	val := []byte("testdata")
	cache := NewCache(interval)
	var wg sync.WaitGroup
	workers := 10
	iterations := 5000

	for range workers {
		wg.Go(func() {
			for j := range iterations {
				key := fmt.Sprintf("key-%d", j)
				cache.Add(key, val)
			}
		})
	}

	for range workers {
		wg.Go(func() {
			for j := range iterations {
				key := fmt.Sprintf("key-%d", j)
				cache.Get(key)
			}
		})
	}

	wg.Wait()
}

func TestAddLarge(t *testing.T) {
	const interval = 2 * time.Minute
	val := bytes.Repeat([]byte("a"), 5000*1024*1024)
	cache := NewCache(interval)

	cache.Add("key", val)

	got, ok := cache.Get("key")
	if !ok {
		t.Fatalf("expected to find key")
	}

	if !bytes.Equal(got, val) {
		t.Errorf("data corruption: retrieved bytes do not match original payload")
	}
}
