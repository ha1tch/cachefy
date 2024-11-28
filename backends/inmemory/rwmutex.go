// File: rwmutex.go

package inmemory

import (
	"sync"
	"time"
)

// RWMutexCache is an in-memory cache implementation with RWMutex for thread safety.
type RWMutexCache struct {
	data       map[string]cacheItem
	mutex      sync.RWMutex
	defaultTTL time.Duration
}

type cacheItem struct {
	value     interface{}
	expiresAt time.Time
}

// NewRWMutexCache creates a new RWMutexCache instance with the provided default TTL.
func NewRWMutexCache(defaultTTL time.Duration) *RWMutexCache {
	return &RWMutexCache{
		data:       make(map[string]cacheItem),
		defaultTTL: defaultTTL,
	}
}

// Get retrieves the value associated with the given key.
func (c *RWMutexCache) Get(key string) (interface{}, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.data[key]
	if !exists || time.Now().After(item.expiresAt) {
		return nil, ErrCacheMiss
	}
	return item.value, nil
}

// Set stores the value associated with the given key.
func (c *RWMutexCache) Set(key string, value interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = cacheItem{
		value:     value,
		expiresAt: time.Now().Add(c.defaultTTL),
	}
	return nil
}

// Delete removes the value associated with the given key.
func (c *RWMutexCache) Delete(key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.data, key)
	return nil
}

// Clear removes all entries from the cache.
func (c *RWMutexCache) Clear() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data = make(map[string]cacheItem)
	return nil
}
