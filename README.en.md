# Cachefy

Cachefy is a flexible caching library for Go, supporting in-memory and persistent caches. It offers various backends, including sharded in-memory caches and persistence layers using SQLite or Postgres.

## Features

- In-memory caching with sync.Map, RWMutex, and sharded backends.
- Optional persistence using SQLite (binary blob storage) or Postgres (JSON storage).
- Configurable time-to-live (TTL) for cache entries.
- Pagination support for efficient cache retrieval from persistent stores.

## Installation

bash
go get github.com/yourusername/cachefy



## Usage

### Basic In-Memory Cache

go
config := CacheConfig{
    DefaultTTL: 5 * time.Minute,
    Backend:    "rwmutex",
}

cache, err := NewCache(config)
if err != nil {
    log.Fatalf("Failed to initialize cache: %v", err)
}

cache.Set("key1", "value1")
value, _ := cache.Get("key1")
fmt.Println(value) // Output: value1



### Sharded In-Memory Cache

go
config := CacheConfig{
    DefaultTTL:    5 * time.Minute,
    Backend:       "sharded",
    Shards:        4,
    ShardCapacity: 100,
}

cache, err := NewCache(config)
if err != nil {
    log.Fatalf("Failed to initialize sharded cache: %v", err)
}



### Persistent Cache (SQLite)

go
config := CacheConfig{
    DefaultTTL:        5 * time.Minute,
    Backend:           "rwmutex",
    EnablePersistence: true,
    DatabaseType:      "sqlite",
    DatabaseDSN:       "cache.db",
}

cache, err := NewCache(config)
if err != nil {
    log.Fatalf("Failed to initialize persistent cache: %v", err)
}



### Persistent Cache (Postgres)

go
config := CacheConfig{
    DefaultTTL:        5 * time.Minute,
    Backend:           "rwmutex",
    EnablePersistence: true,
    DatabaseType:      "postgres",
    DatabaseDSN:       "user=your_user password=your_password dbname=your_db sslmode=disable",
}

cache, err := NewCache(config)
if err != nil {
    log.Fatalf("Failed to initialize persistent cache: %v", err)
}



## Testing

Run the tests with:

bash
go test ./... -v
