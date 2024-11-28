package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

type PostgresRepository struct {
	db *sql.DB
}

// SQL statements as constants
const (
	sqlCreateTablePostgres = `
	CREATE TABLE IF NOT EXISTS cache (
		key TEXT PRIMARY KEY,
		value JSONB,
		expires_at BIGINT
	)`

	sqlGetEntryPostgres = `
	SELECT value, expires_at FROM cache WHERE key = $1`

	sqlInsertOrUpdateEntryPostgres = `
	INSERT INTO cache (key, value, expires_at)
	VALUES ($1, $2, $3)
	ON CONFLICT (key) DO UPDATE
	SET value = EXCLUDED.value,
		expires_at = EXCLUDED.expires_at`

	sqlDeleteEntryPostgres = `
	DELETE FROM cache WHERE key = $1`

	sqlClearEntriesPostgres = `
	DELETE FROM cache`

	sqlPaginateEntriesPostgres = `
	SELECT key, value, expires_at FROM cache
	ORDER BY key ASC LIMIT $1 OFFSET $2`
)

func NewPostgresRepository(dsn string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	repo := &PostgresRepository{db: db}
	if err := repo.initTable(); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *PostgresRepository) initTable() error {
	_, err := r.db.Exec(sqlCreateTablePostgres)
	return err
}

func (r *PostgresRepository) Get(key string) (*CacheEntry, error) {
	row := r.db.QueryRow(sqlGetEntryPostgres, key)

	var jsonValue []byte
	var expiresAt int64
	err := row.Scan(&jsonValue, &expiresAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("key not found")
	} else if err != nil {
		return nil, err
	}

	// Check expiration
	if time.Now().Unix() > expiresAt {
		_ = r.Delete(key)
		return nil, errors.New("key expired")
	}

	var value interface{}
	if err := json.Unmarshal(jsonValue, &value); err != nil {
		return nil, err
	}

	return &CacheEntry{Key: key, Value: value, ExpiresAt: expiresAt}, nil
}

func (r *PostgresRepository) Set(entry *CacheEntry) error {
	jsonValue, err := json.Marshal(entry.Value)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(sqlInsertOrUpdateEntryPostgres, entry.Key, jsonValue, entry.ExpiresAt)
	return err
}

func (r *PostgresRepository) Delete(key string) error {
	_, err := r.db.Exec(sqlDeleteEntryPostgres, key)
	return err
}

func (r *PostgresRepository) Clear() error {
	_, err := r.db.Exec(sqlClearEntriesPostgres)
	return err
}

func (r *PostgresRepository) Paginate(offset, limit int) ([]*CacheEntry, error) {
	rows, err := r.db.Query(sqlPaginateEntriesPostgres, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*CacheEntry
	for rows.Next() {
		var key string
		var jsonValue []byte
		var expiresAt int64

		if err := rows.Scan(&key, &jsonValue, &expiresAt); err != nil {
			return nil, err
		}

		var value interface{}
		if err := json.Unmarshal(jsonValue, &value); err != nil {
			return nil, err
		}

		entries = append(entries, &CacheEntry{Key: key, Value: value, ExpiresAt: expiresAt})
	}
	return entries, nil
}

