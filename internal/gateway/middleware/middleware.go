package middleware

import (
	"net/http"

	"github.com/martinsdevv/aegis/internal/config"
	"github.com/redis/go-redis/v9"
)

func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

func NewMiddleware(handler http.Handler, cfg config.Config, rlStore *RLStore, quotaMgr *QuotaManager, redisClient *redis.Client, apiKeyStore *APIKeyStore) http.Handler {
	return Chain(handler,
		RequestID(),
		ContentID(),
		Recover,
		WithAPIKey(apiKeyStore),
		RateLimit(rlStore),
		quotaMgr.Enforce,
		Logger,
		PublishUsage(redisClient, "aether.usage.v1"),
	)
}
