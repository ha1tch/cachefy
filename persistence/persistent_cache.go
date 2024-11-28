// File: persistent_cache.go

package persistence

import (
	"cachefy/interfaces"
	"cachefy/repository"
	"sync"
)

// PersistentCache is a cache that wraps another Cache and persists data using a Repository.
type PersistentCache struct {
	cache interfaces.Cache
	repo  repository.Repository
	mutex sync.Mutex
}

// NewPersistentCache creates a new PersistentCache.
func NewPersistentCache(cache interfaces.Cache, repo repository.Repository) *PersistentCache {
	return &PersistentCache{
		cache: cache,
		repo:  repo,
	}
}

// Set adds or updates a cache entry and persists it.
func (p *PersistentCache) Set(key string, value interface{}) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	err := p.cache.Set(key, value)
	if err != nil {
		return err
	}

	// Persist the entry
	entry := &repository.CacheEntry{
		Key:       key,
		Value:     value,
		ExpiresAt: 0, // TTL not implemented for persistence
	}
	return p.repo.Set(entry)
}

// Get retrieves a value from the cache.
func (p *PersistentCache) Get(key string) (interface{}, error) {
	return p.cache.Get(key)
}

// Delete removes a value from the cache and the repository.
func (p *PersistentCache) Delete(key string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	err := p.cache.Delete(key)
	if err != nil {
		return err
	}

	return p.repo.Delete(key)
}

// Clear removes all entries from the cache and the repository.
func (p *PersistentCache) Clear() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	err := p.cache.Clear()
	if err != nil {
		return err
	}

	return p.repo.Clear()
}
