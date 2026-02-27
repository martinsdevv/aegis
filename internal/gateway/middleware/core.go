package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type (
	ctxKeyContentID struct{}
	ctxKeyRequestID struct{}
	Middleware      func(http.Handler) http.Handler
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.wroteHeader = true
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

// NewMiddleware aplica todos os middlewares na ordem correta

func ContentID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cid := r.Header.Get("X-Content-ID")
			if cid == "" {
				cid = uuid.NewString()
				r.Header.Set("X-Content-ID", cid)
			}

			ctx := context.WithValue(r.Context(), ctxKeyContentID{}, cid)
			r = r.WithContext(ctx)

			wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}
			wrapped.Header().Set("X-Content-ID", cid)

			next.ServeHTTP(wrapped, r)
		})
	}
}

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := r.Header.Get("X-Request-ID")
			if reqID == "" {
				reqID = uuid.NewString()
				r.Header.Set("X-Request-ID", reqID)
			}

			ctx := context.WithValue(r.Context(), ctxKeyRequestID{}, reqID)
			r = r.WithContext(ctx)

			wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}
			wrapped.Header().Set("X-Request-ID", reqID)

			next.ServeHTTP(wrapped, r)
		})
	}
}

func ContentIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxKeyContentID{}).(string)
	return v, ok && v != "missing"
}

func RequestIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxKeyRequestID{}).(string)
	return v, ok && v != ""
}
