// File: async_persistence.go

package persistence

import (
	"log"
	"sync"
	"time"

	"cachefy/repository"
)

type AsyncPersistenceManager struct {
	repo       repository.Repository
	taskQueue  chan *repository.CacheEntry
	retryLimit int
	wg         sync.WaitGroup
}

func NewAsyncPersistenceManager(repo repository.Repository, queueSize, retryLimit int) *AsyncPersistenceManager {
	manager := &AsyncPersistenceManager{
		repo:       repo,
		taskQueue:  make(chan *repository.CacheEntry, queueSize),
		retryLimit: retryLimit,
	}

	// Start background workers
	manager.wg.Add(1)
	go manager.worker()

	return manager
}

// Enqueue adds a cache entry to the persistence queue.
func (m *AsyncPersistenceManager) Enqueue(entry *repository.CacheEntry) {
	m.taskQueue <- entry
}

// worker processes persistence tasks from the queue.
func (m *AsyncPersistenceManager) worker() {
	defer m.wg.Done()

	for entry := range m.taskQueue {
		for attempts := 0; attempts < m.retryLimit; attempts++ {
			err := m.repo.Set(entry)
			if err == nil {
				break
			}
			log.Printf("Persistence failed for key %s: %v (retry %d)", entry.Key, err, attempts+1)
			time.Sleep(2 * time.Second)
		}
	}
}

// Shutdown gracefully stops the manager and waits for pending tasks to complete.
func (m *AsyncPersistenceManager) Shutdown() {
	close(m.taskQueue)
	m.wg.Wait()
}
