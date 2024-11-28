// File: cachefy.go

package cachefy

import (
	"cachefy/backends/inmemory"
	"cachefy/interfaces"
	"cachefy/persistence"
	"cachefy/repository"
	"errors"
	"log"
	"time"
)

type CacheConfig struct {
	DefaultTTL               time.Duration
	Backend                  string
	Shards                   int
	ShardCapacity            int
	EnablePersistence        bool
	PersistenceFilePath      string
	PersistenceFlushInterval time.Duration
	DatabaseType             string // "sqlite" or "postgres"
	DatabaseDSN              string // Database connection string
}

func NewCache(config CacheConfig) (interfaces.Cache, error) {
	// Validate config
	if config.DefaultTTL <= 0 {
		return nil, errors.New("defaultTTL must be greater than zero")
	}

	var cache interfaces.Cache
	switch config.Backend {
	case "syncmap":
		cache = inmemory.NewSyncMapCache(config.DefaultTTL)
	case "rwmutex":
		cache = inmemory.NewRWMutexCache(config.DefaultTTL)
	case "sharded":
		if config.Shards <= 0 || config.ShardCapacity <= 0 {
			return nil, errors.New("shards and shardCapacity must be greater than zero")
		}
		cache = inmemory.NewShardedCache(config.Shards, config.DefaultTTL, config.ShardCapacity)
	default:
		return nil, errors.New("unsupported backend")
	}

	// Add persistence if enabled
	if config.EnablePersistence {
		var repo repository.Repository
		//var serializer serialization.Serializer

		var err error

		switch config.DatabaseType {
		case "sqlite":
			// serializer = &serialization.BlobSerializer{}
			//  repo, err = repository.NewSQLiteRepository(config.DatabaseDSN, serializer)
			repo, err = repository.NewSQLiteRepository(config.DatabaseDSN)
		case "postgres":
			//serializer = &serialization.JSONSerializer{}
			// repo, err = repository.NewPostgresRepository(config.DatabaseDSN, serializer)
			repo, err = repository.NewPostgresRepository(config.DatabaseDSN)
		default:
			return nil, errors.New("unsupported database type")
		}

		if err != nil {
			log.Printf("Failed to initialize persistence repository: %v", err)
			return nil, err
		}

		cache = persistence.NewPersistentCache(cache, repo)
		log.Println("Persistence layer enabled.")
	}

	log.Println("Cache successfully initialized.")
	return cache, nil
}
