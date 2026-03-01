package middleware

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrAPIKeyNotFound = errors.New("api key not found")

type APIKeyStore struct {
	db    *sql.DB
	redis *redis.Client
	ttl   time.Duration
}

func NewAPIKeyStore(db *sql.DB, redisClient *redis.Client, ttl time.Duration) *APIKeyStore {
	return &APIKeyStore{
		db:    db,
		redis: redisClient,
		ttl:   ttl,
	}
}

func (s *APIKeyStore) FindByHash(ctx context.Context, hash string) (*APIKey, error) {

	// Redis
	if s.redis != nil {
		if val, err := s.redis.Get(ctx, s.redisKey(hash)).Result(); err == nil {
			var k APIKey
			if err := json.Unmarshal([]byte(val), &k); err == nil {
				return &k, nil
			}
		}
	}

	// Postgres
	k, err := s.findInDB(ctx, hash)
	if err != nil {
		return nil, err
	}

	// Cache apenas se ativa
	if k != nil && k.Active && s.redis != nil {
		b, _ := json.Marshal(k)
		_ = s.redis.Set(ctx, s.redisKey(hash), b, s.ttl).Err()
	}

	return k, nil
}

func (s *APIKeyStore) DeleteFromCache(ctx context.Context, hash string) error {
	if s.redis == nil {
		return nil
	}
	return s.redis.Del(ctx, s.redisKey(hash)).Err()
}

func (s *APIKeyStore) findInDB(ctx context.Context, hash string) (*APIKey, error) {
	const query = `
		SELECT id, key, name, upstream_host, is_active, monthly_quota, created_at
		FROM api_keys
		WHERE key = $1
		LIMIT 1
	`

	row := s.db.QueryRowContext(ctx, query, hash)

	var k APIKey
	err := row.Scan(
		&k.ID,
		&k.KeyHash,
		&k.Name,
		&k.UpstreamHost,
		&k.Active,
		&k.MonthlyQuota,
		&k.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrAPIKeyNotFound
	}
	if err != nil {
		return nil, err
	}

	return &k, nil
}

func (s *APIKeyStore) redisKey(hash string) string {
	return "aegis:apikey:" + hash
}
