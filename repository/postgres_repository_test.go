
// File: postgres_repository_test.go

package repository

import (
	"testing"
	"time"

	"cachefy/serialization"
)

func TestPostgresRepository(t *testing.T) {
	// Mock DSN for testing (ensure test database setup before running this test)
	dsn := "user=testuser password=testpassword dbname=testdb sslmode=disable"
	serializer := &serialization.JSONSerializer{}
	repo, err := NewPostgresRepository(dsn, serializer)
	if err != nil {
		t.Fatalf("Failed to create Postgres repository: %v", err)
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
