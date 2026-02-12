package middleware

import (
	"net/http"
)

func Keyring(keys []string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			kh := r.Header.Get("X-API-Key")

			if kh == "" {
				http.Error(w, "X-API-Key header is absent", http.StatusUnauthorized)
				return
			}

			set := make(map[string]struct{}, len(keys))
			for _, v := range keys {
				set[v] = struct{}{}
			}

			if _, ok := set[kh]; !ok {
				http.Error(w, "X-API-Key header is invalid", http.StatusForbidden)
				return
			}
			r.Header.Del("X-API-Key")
			next.ServeHTTP(w, r)
		})
	}
}
