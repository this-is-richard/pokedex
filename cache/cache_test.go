package cache

import (
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	cache := NewCache[string](5 * time.Millisecond)
	key := "foobar"
	value := "abc"

	cache.Add(key, value)
	val, ok := cache.Get(key)
	if !ok {
		t.Error("failed to lookup by key")
	}
	if val != value {
		t.Errorf("expect %v, actual %v", value, val)
	}
}

func TestConcurrentAddGet(t *testing.T) {
	cache := NewCache[string](5 * time.Millisecond)
	key := "foobar"

	for i := 0; i < 100; i++ {
		go func() {
			cache.Add(key, "abc")
		}()
		go func() {
			cache.Get(key)
		}()
	}
}

func TestCacheTimeout(t *testing.T) {
	cache := NewCache[int](5 * time.Millisecond)
	key := "foobar"
	value := 123

	cache.Add(key, value)
	time.Sleep(3 * time.Millisecond)
	val, ok := cache.Get(key)
	if !ok {
		t.Error("failed to lookup by key")
	}
	if val != value {
		t.Errorf("expect %v, actual %v", value, val)
	}

	time.Sleep(4 * time.Millisecond)
	_, ok = cache.Get(key)
	if ok {
		t.Error("expect key-value removed, but found")
	}
}
