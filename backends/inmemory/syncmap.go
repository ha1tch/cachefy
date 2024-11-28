// File: syncmap.go

package inmemory

import (
	"sync"
	"time"
)

type SyncMapCache struct {
	data       sync.Map
	defaultTTL time.Duration
}

type syncMapItem struct {
	value     interface{}
	expiresAt time.Time
}

func NewSyncMapCache(defaultTTL time.Duration) *SyncMapCache {
	return &SyncMapCache{
		defaultTTL: defaultTTL,
	}
}

func (c *SyncMapCache) Get(key string) (interface{}, error) {
	item, ok := c.data.Load(key)
	if !ok {
		return nil, ErrCacheMiss
	}

	cachedItem := item.(syncMapItem)
	if time.Now().After(cachedItem.expiresAt) {
		c.data.Delete(key)
		return nil, ErrCacheMiss
	}
	return cachedItem.value, nil
}

func (c *SyncMapCache) Set(key string, value interface{}) error {
	c.data.Store(key, syncMapItem{
		value:     value,
		expiresAt: time.Now().Add(c.defaultTTL),
	})
	return nil
}

func (c *SyncMapCache) Delete(key string) error {
	c.data.Delete(key)
	return nil
}

func (c *SyncMapCache) Clear() error {
	c.data.Range(func(key, value interface{}) bool {
		c.data.Delete(key)
		return true
	})
	return nil
}
