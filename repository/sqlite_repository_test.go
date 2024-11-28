
// File: sqlite_repository_test.go

package repository

import (
	"testing"
	"time"

	"cachefy/serialization"
)

func TestSQLiteRepository(t *testing.T) {
	serializer := &serialization.BlobSerializer{}
	repo, err := NewSQLiteRepository(":memory:", serializer)
	if err != nil {
		t.Fatalf("Failed to create SQLite repository: %v", err)
	}

	entry := &CacheEntry{
		Key:       "key1",
		Value:     "value1",
		ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
	}

	// Test Set
	err = repo.Set(entry)
	if err != nil {
		t.Fatalf("Failed to set cache entry: %v", err)
	}

	// Test Get
	retrieved, err := repo.Get("key1")
	if err != nil || retrieved.Value != "value1" {
		t.Fatalf("Failed to get cache entry: %v", err)
	}
}
