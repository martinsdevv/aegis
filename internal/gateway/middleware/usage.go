package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

func PublishUsage(redisClient *redis.Client, streamName string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()

			// buffer de resposta
			buf := newResponseBuffer()

			next.ServeHTTP(buf, r)

			apiKey, ok := APIKeyFromContext(r.Context())
			if !ok {
				buf.FlushTo(w)
				return
			}

			reqID, _ := RequestIDFromContext(r.Context())

			event := UsageEvent{
				EventID:    uuid.NewString(),
				EventVer:   1,
				RequestID:  reqID,
				APIKeyID:   strconv.FormatInt(apiKey.ID, 10),
				Upstream:   apiKey.UpstreamHost,
				Path:       r.URL.Path,
				Method:     r.Method,
				StatusCode: buf.status,
				LatencyMS:  time.Since(start).Milliseconds(),
				Timestamp:  time.Now().UTC().Format(time.RFC3339),
			}

			payload, err := json.Marshal(event)
			if err != nil {
				log.Println("marshal error:", err)
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			if redisClient != nil {
				err = redisClient.XAdd(context.Background(), &redis.XAddArgs{
					Stream: streamName,
					Values: map[string]interface{}{"payload": payload},
				}).Err()

				if err != nil {
					log.Println("redis error:", err)
					http.Error(w, "service unavailable", http.StatusServiceUnavailable)
					return
				}
			}

			// só escreve depois do sucesso
			buf.FlushTo(w)
		})
	}
}

type responseBuffer struct {
	header http.Header
	body   []byte
	status int
}

func newResponseBuffer() *responseBuffer {
	return &responseBuffer{
		header: make(http.Header),
		status: http.StatusOK,
	}
}

func (r *responseBuffer) Header() http.Header {
	return r.header
}

func (r *responseBuffer) Write(b []byte) (int, error) {
	r.body = append(r.body, b...)
	return len(b), nil
}

func (r *responseBuffer) WriteHeader(statusCode int) {
	r.status = statusCode
}

func (r *responseBuffer) FlushTo(w http.ResponseWriter) {
	for k, v := range r.header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	w.WriteHeader(r.status)
	w.Write(r.body)
}
