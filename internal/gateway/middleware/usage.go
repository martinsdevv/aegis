package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type UsageEvent struct {
	EventID    string `json:"event_id"`
	EventVer   int    `json:"event_version"`
	RequestID  string `json:"request_id"`
	APIKeyID   string `json:"api_key_id"`
	Upstream   string `json:"upstream"`
	Path       string `json:"path"`
	Method     string `json:"method"`
	StatusCode int    `json:"status_code"`
	LatencyMS  int64  `json:"latency_ms"`
	Timestamp  string `json:"timestamp"`
}

func NewRedisClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

func PublishUsage(redisClient *redis.Client, streamName string, upstreamHost string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			latency := time.Since(start)
			reqID, _ := RequestIDFromContext(r.Context())
			apiKey, _ := APIKeyFromContext(r.Context())

			event := UsageEvent{
				EventID:    uuid.NewString(),
				EventVer:   1,
				RequestID:  reqID,
				APIKeyID:   apiKey,
				Upstream:   upstreamHost,
				Path:       r.URL.Path,
				Method:     r.Method,
				StatusCode: wrapped.status,
				LatencyMS:  latency.Milliseconds(),
				Timestamp:  time.Now().UTC().Format(time.RFC3339),
			}

			payload, err := json.Marshal(event)
			if err != nil {
				http.Error(wrapped, "internal server error", http.StatusInternalServerError)
				return
			}

			err = redisClient.XAdd(context.Background(), &redis.XAddArgs{
				Stream: streamName,
				Values: map[string]interface{}{"payload": payload},
			}).Err()

			if err != nil {
				fmt.Println("Redis XAdd error:", err)
				http.Error(w, "service unavailable", http.StatusServiceUnavailable)
				return
			} else {
				fmt.Println("Published to Redis:", string(payload))
			}
		})
	}
}
