package gtwhttp

import (
	"net/http"

	"github.com/martinsdevv/aegis/internal/health"
)

func NewRouter(healthCheck *health.Checker) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", health.HealthHandler(healthCheck))
	return mux
}
