package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(wrapped, r)

		host := r.Host
		if strings.HasPrefix(r.URL.Path, "/proxy") {
			if ups := wrapped.Header().Get("X-Upstream-Host"); ups != "" {
				host = ups
			}
		}

		slog.Info("request processed",
			"method", r.Method,
			"path", r.URL.Path,
			"host", host,
			"status", wrapped.status,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}
