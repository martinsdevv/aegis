package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start).Milliseconds()
		host := r.Host

		var apiKeyID int64
		var apiKeyName string
		var quotaCount int64
		var quotaLimit int64

		if apiKey, ok := APIKeyFromContext(r.Context()); ok {
			apiKeyID = apiKey.ID
			apiKeyName = apiKey.Name
			if apiKey.UpstreamHost != "" {
				host = apiKey.UpstreamHost
			}
		}

		if _, ok := r.Context().Value("quota_count").(string); ok {
			if client := r.Context().Value("quota_limit"); client != nil {
				if ql, ok := client.(int64); ok {
					quotaCount = 0 // podemos recuperar do Redis se quiser depois
					quotaLimit = ql
				}
			}
		}

		// log estruturado
		slog.Info("request processed",
			"method", r.Method,
			"path", r.URL.Path,
			"host", host,
			"status", wrapped.status,
			"duration_ms", duration,
			"api_key_id", apiKeyID,
			"api_key_name", apiKeyName,
			"quota_count", quotaCount,
			"quota_limit", quotaLimit,
		)
	})
}
