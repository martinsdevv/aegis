package upshttp

import (
	"net/http"

	health "github.com/martinsdevv/aegis/internal/health"
)

func NewRouter(healthChecker *health.Checker) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", HandlePing)
	mux.HandleFunc("/echo", HandleEcho)
	mux.HandleFunc("/healthz", health.HealthHandler(healthChecker))

	return mux
}
