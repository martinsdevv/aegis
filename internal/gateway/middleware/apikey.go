package middleware

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"net/http"
	"strings"
	"time"
)

type ctxKeyAPIKey struct{}

// APIKey representa uma API Key persistida no banco
type APIKey struct {
	ID           int64
	KeyHash      string
	Name         string
	UpstreamHost string
	Active       bool
	MonthlyQuota int
	CreatedAt    time.Time
}

func APIKeyFromContext(ctx context.Context) (*APIKey, bool) {
	v, ok := ctx.Value(ctxKeyAPIKey{}).(*APIKey)
	return v, ok && v != nil
}

func WithAPIKey(store *APIKeyStore) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			rawKey := strings.TrimSpace(r.Header.Get("X-API-Key"))
			if rawKey == "" {
				http.Error(w, "X-API-Key header is absent", http.StatusUnauthorized)
				return
			}

			hashed := hashKey(rawKey)

			apiKey, err := store.FindByHash(r.Context(), hashed)
			if err != nil {
				if err == ErrAPIKeyNotFound {
					http.Error(w, "invalid api key", http.StatusForbidden)
					return
				}
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			if !apiKey.Active {
				http.Error(w, "api key disabled", http.StatusForbidden)
				return
			}

			r.Header.Del("X-API-Key")

			ctx := context.WithValue(r.Context(), ctxKeyAPIKey{}, apiKey)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func findAPIKey(ctx context.Context, db *sql.DB, hash string) (*APIKey, error) {
	const query = `
		SELECT id, key, name, upstream_host, is_active, monthly_quota, created_at
		FROM api_keys
		WHERE key = $1
		LIMIT 1
	`

	row := db.QueryRowContext(ctx, query, hash)

	var k APIKey
	err := row.Scan(
		&k.ID,
		&k.KeyHash,
		&k.Name,
		&k.UpstreamHost,
		&k.Active,
		&k.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &k, nil
}

func hashKey(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
