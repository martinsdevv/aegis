package middleware

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type QuotaManager struct {
	client       *redis.Client
	fallbackMap  sync.Map
	currentMonth string
	mu           sync.Mutex
}

func NewQuotaManager(client *redis.Client) *QuotaManager {
	return &QuotaManager{
		client:       client,
		currentMonth: time.Now().UTC().Format("2006-01"),
	}
}

func (qm *QuotaManager) Enforce(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey, ok := APIKeyFromContext(r.Context())
		if !ok {
			http.Error(w, "missing API key", http.StatusUnauthorized)
			return
		}

		limit := int64(apiKey.MonthlyQuota)
		if limit <= 0 {
			limit = 10000
		}

		now := time.Now().UTC()
		month := now.Format("2006-01")

		// Reset fallback in-memory se mudou de mês
		qm.mu.Lock()
		if month != qm.currentMonth {
			qm.fallbackMap = sync.Map{}
			qm.currentMonth = month
		}
		qm.mu.Unlock()

		apiKeyID := strconv.FormatInt(apiKey.ID, 10)
		key := "quota:" + apiKeyID + ":" + month

		allowed := false

		// Redis
		if qm.client != nil {
			incr := qm.client.Incr(context.Background(), key)
			qm.client.ExpireAt(context.Background(), key, time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.UTC))

			if val, err := incr.Result(); err == nil && val <= limit {
				allowed = true
			}
		}

		// fallback in-memory
		if !allowed {
			v, _ := qm.fallbackMap.LoadOrStore(key, int64(0))
			cnt := v.(int64) + 1
			qm.fallbackMap.Store(key, cnt)
			if cnt <= limit {
				allowed = true
			}
		}

		if !allowed {
			http.Error(w, "quota exceeded", http.StatusForbidden)
			return
		}

		// Coloca contagem no contexto para Logger
		ctx := context.WithValue(r.Context(), "quota_count", key)
		ctx = context.WithValue(ctx, "quota_limit", limit)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
