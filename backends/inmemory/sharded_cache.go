// File: sharded_cache.go

package inmemory

import (
	"hash/fnv"
	"time"
)

// ShardedCache is a thread-safe in-memory cache with multiple shards for scalability.
type ShardedCache struct {
	shards        []*RWMutexCache
	shardCount    int
	defaultTTL    time.Duration
	shardCapacity int
}

// NewShardedCache initializes a ShardedCache with the given number of shards and default TTL.
func NewShardedCache(shardCount int, defaultTTL time.Duration, shardCapacity int) *ShardedCache {
	shards := make([]*RWMutexCache, shardCount)
	for i := 0; i < shardCount; i++ {
		shards[i] = NewRWMutexCache(defaultTTL)
	}
	return &ShardedCache{
		shards:        shards,
		shardCount:    shardCount,
		defaultTTL:    defaultTTL,
		shardCapacity: shardCapacity,
	}
}

// hashKey determines the shard index for a given key.
func (c *ShardedCache) hashKey(key string) int {
	hasher := fnv.New32a()
	hasher.Write([]byte(key))
	return int(hasher.Sum32()) % c.shardCount
}

// Get retrieves a value from the cache.
func (c *ShardedCache) Get(key string) (interface{}, error) {
	shard := c.shards[c.hashKey(key)]
	return shard.Get(key)
}

// Set stores a value in the cache.
func (c *ShardedCache) Set(key string, value interface{}) error {
	shard := c.shards[c.hashKey(key)]
	return shard.Set(key, value)
}

// Delete removes a value from the cache.
func (c *ShardedCache) Delete(key string) error {
	shard := c.shards[c.hashKey(key)]
	return shard.Delete(key)
}

// Clear removes all entries from the cache.
func (c *ShardedCache) Clear() error {
	for _, shard := range c.shards {
		shard.Clear()
	}
	return nil
}
