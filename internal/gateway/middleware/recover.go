package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"runtime/debug"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				cid, _ := ContentIDFromContext(r.Context())
				slog.Error("panic recovered",
					"method", r.Method,
					"path", r.URL.Path,
					"X-Content-ID", cid,
					"panic", err,
					"stack", string(debug.Stack()),
				)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error": "internal server error"}`))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func ContentIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxKeyContentID{}).(string)
	return v, ok && v != "missing"
}
