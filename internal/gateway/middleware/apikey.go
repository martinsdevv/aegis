package middleware

import (
	"context"
	"net/http"
	"strings"
)

type ctxKeyAPIKey struct{}

func APIKeyFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxKeyAPIKey{}).(string)
	return v, ok && v != ""
}

func WithAPIKey() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			k := strings.TrimSpace(r.Header.Get("X-API-Key"))
			if k == "" {
				http.Error(w, "X-API-Key header is absent", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), ctxKeyAPIKey{}, k)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Keyring(keys []string) Middleware {
	set := make(map[string]struct{}, len(keys))
	for _, v := range keys {
		v = strings.TrimSpace(v)
		if v != "" {
			set[v] = struct{}{}
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			k, ok := APIKeyFromContext(r.Context())
			if !ok {
				http.Error(w, "missing api key", http.StatusUnauthorized)
				return
			}
			if _, ok := set[k]; !ok {
				http.Error(w, "X-API-Key header is invalid", http.StatusForbidden)
				return
			}
			r.Header.Del("X-API-Key")
			next.ServeHTTP(w, r)
		})
	}
}
