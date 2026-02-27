package middleware

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type rlEntry struct {
	lim      *rate.Limiter
	lastSeen time.Time
}

type RLStore struct {
	mu    sync.Mutex
	m     map[string]*rlEntry
	r     rate.Limit
	burst int
	ttl   time.Duration
}

func NewRLStore(r rate.Limit, burst int, ttl time.Duration) *RLStore {
	return &RLStore{
		m:     make(map[string]*rlEntry),
		r:     r,
		burst: burst,
		ttl:   ttl,
	}
}

func (s *RLStore) get(key string) *rate.Limiter {
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()

	if e, ok := s.m[key]; ok {
		e.lastSeen = now
		return e.lim
	}

	lim := rate.NewLimiter(s.r, s.burst)
	s.m[key] = &rlEntry{lim: lim, lastSeen: now}
	return lim
}

func (s *RLStore) Cleanup() {
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()

	for k, e := range s.m {
		if now.Sub(e.lastSeen) > s.ttl {
			delete(s.m, k)
		}
	}
}

func RateLimit(store *RLStore) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			k, ok := APIKeyFromContext(r.Context())
			if !ok {
				http.Error(w, "missing API key", http.StatusUnauthorized)
				return
			}
			lim := store.get(k)
			if !lim.Allow() {
				w.Header().Set("Retry-After", "1")
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
