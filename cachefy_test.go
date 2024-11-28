
// File: cachefy_test.go

package cachefy

import (
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	config := CacheConfig{
		DefaultTTL:        5 * time.Minute,
		Backend:           "rwmutex",
		EnablePersistence: false,
	}

	cache, err := NewCache(config)
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	err = cache.Set("test", "value")
	if err != nil {
		t.Fatalf("Failed to set cache value: %v", err)
	}

	value, err := cache.Get("test")
	if err != nil || value != "value" {
		t.Fatalf("Failed to get cache value: %v", err)
	}
}
