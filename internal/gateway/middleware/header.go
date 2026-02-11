package middleware

import (
	"context"
	"net/http"
)

type (
	ctxKeyContentID struct{}
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

func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

func NewMiddleware(handler http.Handler) http.Handler {
	return Chain(handler, ContentID("contentID"), Logger, Recover)
}

func ContentID(defaultID string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cid := r.Header.Get("X-Content-ID")
			if cid == "" {
				cid = defaultID
				r.Header.Set("X-Content-ID", cid)
			}

			ctx := context.WithValue(r.Context(), ctxKeyContentID{}, cid)
			r = r.WithContext(ctx)

			wrapped := &responseWriter{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			wrapped.Header().Set("X-Content-ID", cid)
			next.ServeHTTP(wrapped, r)
		})
	}
}
