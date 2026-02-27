package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type QuotaManager struct {
	client      *redis.Client
	fallbackMap sync.Map
	limit       int
}

func NewQuotaManager(client *redis.Client, limit int) *QuotaManager {
	return &QuotaManager{client: client, limit: limit}
}

func (qm *QuotaManager) Enforce(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		k, _ := APIKeyFromContext(r.Context())
		if k == "" {
			http.Error(w, "missing API key", http.StatusUnauthorized)
			return
		}

		month := time.Now().UTC().Format("2006-01")
		key := "quota:" + k + ":" + month
		allowed := false

		if qm.client != nil {
			incr := qm.client.Incr(context.Background(), key)
			qm.client.Expire(context.Background(), key, 30*24*time.Hour)
			if val, err := incr.Result(); err == nil && val <= int64(qm.limit) {
				allowed = true
			}
		}

		if !allowed {
			v, _ := qm.fallbackMap.LoadOrStore(key, int64(0))
			cnt := v.(int64) + 1
			qm.fallbackMap.Store(key, cnt)
			if cnt <= int64(qm.limit) {
				allowed = true
			}
		}

		if !allowed {
			http.Error(w, "quota exceeded", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
