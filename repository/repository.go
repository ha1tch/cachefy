
// File: repository.go

package repository

// CacheEntry represents a single cache entry in the repository.
type CacheEntry struct {
	Key       string      `json:"key"`        // Cache key
	Value     interface{} `json:"value"`      // Cache value
	ExpiresAt int64       `json:"expires_at"` // Expiration timestamp (Unix time)
}

// Repository defines the interface for cache persistence operations.
type Repository interface {
	Get(key string) (*CacheEntry, error)
	Set(entry *CacheEntry) error
	Delete(key string) error
	Clear() error
	Paginate(offset, limit int) ([]*CacheEntry, error)
}
