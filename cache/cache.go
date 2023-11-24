package cache

import (
	"sync"
	"time"
)

type Record[T interface{}] struct {
	Value     T
	CreatedAt time.Time
}

// Cache is a memoization implementation by a map[string]interface{}.
// Cache is safe for concurrent use by multiple goroutines.
type Cache[T interface{}] struct {
	mem      map[string]Record[T]
	mu       sync.RWMutex
	interval time.Duration
}

func NewCache[T interface{}](interval time.Duration) *Cache[T] {
	cache := Cache[T]{
		mem:      make(map[string]Record[T]),
		interval: interval,
	}
	go cache.Watch()
	return &cache
}

func (c *Cache[T]) Get(key string) (T, bool) {
	c.mu.RLock()
	record, ok := c.mem[key]
	c.mu.RUnlock()
	return record.Value, ok
}

func (c *Cache[T]) Add(key string, val T) {
	c.mu.Lock()
	c.mem[key] = Record[T]{Value: val, CreatedAt: time.Now()}
	c.mu.Unlock()
}

func (c *Cache[T]) Watch() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for t := range ticker.C {
		c.mu.Lock()
		for key, record := range c.mem {
			if t.Sub(record.CreatedAt) > c.interval {
				delete(c.mem, key)
			}
		}
		c.mu.Unlock()
	}
}
