package interfaces

type Cache interface {
    Get(key string) (interface{}, error)
    Set(key string, value interface{}) error
    Delete(key string) error
    Clear() error
}

type Repository interface {
    Get(key string) (*CacheEntry, error)
    Set(entry *CacheEntry) error
    Delete(key string) error
    Clear() error
    Paginate(offset, limit int) ([]*CacheEntry, error)
}

type CacheEntry struct {
    Key       string
    Value     interface{}
    ExpiresAt int64
}
