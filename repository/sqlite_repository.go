package repository

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// SQLiteRepository is a repository implementation for SQLite.
type SQLiteRepository struct {
	db *sql.DB
}

// SQL statements as constants
const (
	sqlCreateTable = `
	CREATE TABLE IF NOT EXISTS cache (
		key TEXT PRIMARY KEY,
		value BLOB,
		expires_at INTEGER
	)`

	sqlGetEntry = `
	SELECT value, expires_at FROM cache WHERE key = ?`

	sqlInsertOrUpdateEntry = `
	INSERT INTO cache (key, value, expires_at)
	VALUES (?, ?, ?)
	ON CONFLICT(key) DO UPDATE SET
		value = excluded.value,
		expires_at = excluded.expires_at`

	sqlDeleteEntry = `
	DELETE FROM cache WHERE key = ?`

	sqlClearEntries = `
	DELETE FROM cache`

	sqlPaginateEntries = `
	SELECT key, value, expires_at FROM cache
	ORDER BY key ASC LIMIT ? OFFSET ?`
)

// NewSQLiteRepository creates a new repository instance connected to an SQLite database.
func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	repo := &SQLiteRepository{db: db}
	if err := repo.initTable(); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *SQLiteRepository) initTable() error {
	_, err := r.db.Exec(sqlCreateTable)
	return err
}

func (r *SQLiteRepository) Get(key string) (*CacheEntry, error) {
	row := r.db.QueryRow(sqlGetEntry, key)

	var value []byte
	var expiresAt int64
	err := row.Scan(&value, &expiresAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("key not found")
	} else if err != nil {
		return nil, err
	}

	// Check expiration
	if time.Now().Unix() > expiresAt {
		_ = r.Delete(key) // Automatically clean up expired entries
		return nil, errors.New("key expired")
	}

	return &CacheEntry{Key: key, Value: value, ExpiresAt: expiresAt}, nil
}

func (r *SQLiteRepository) Set(entry *CacheEntry) error {
	_, err := r.db.Exec(sqlInsertOrUpdateEntry, entry.Key, entry.Value, entry.ExpiresAt)
	return err
}

func (r *SQLiteRepository) Delete(key string) error {
	_, err := r.db.Exec(sqlDeleteEntry, key)
	return err
}

func (r *SQLiteRepository) Clear() error {
	_, err := r.db.Exec(sqlClearEntries)
	return err
}

func (r *SQLiteRepository) Paginate(offset, limit int) ([]*CacheEntry, error) {
	rows, err := r.db.Query(sqlPaginateEntries, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*CacheEntry
	for rows.Next() {
		var key string
		var value []byte
		var expiresAt int64

		if err := rows.Scan(&key, &value, &expiresAt); err != nil {
			return nil, err
		}

		entries = append(entries, &CacheEntry{Key: key, Value: value, ExpiresAt: expiresAt})
	}
	return entries, nil
}

