
// File: sharded_cache_test.go

package inmemory

import (
	"testing"
	"time"
)

func TestShardedCache(t *testing.T) {
	cache := NewShardedCache(4, 2*time.Second, 100)

	// Test Set and Get
	cache.Set("key1", "value1")
	val, err := cache.Get("key1")
	if err != nil || val != "value1" {
		t.Errorf("Expected value1, got %v, error: %v", val, err)
	}

	// Test expiration
	time.Sleep(3 * time.Second)
	val, err = cache.Get("key1")
	if err == nil {
		t.Errorf("Expected cache miss, got value: %v", val)
	}

	// Test Clear
	cache.Set("key2", "value2")
	cache.Clear()
	_, err = cache.Get("key2")
	if err == nil {
		t.Errorf("Expected cache miss after Clear")
	}
}
