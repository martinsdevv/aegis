package gtwhttp

import (
	"log"
	"net/http"

	"github.com/martinsdevv/aegis/internal/config"
	"github.com/martinsdevv/aegis/internal/gateway/middleware"
	"github.com/martinsdevv/aegis/internal/gateway/proxy"
	"github.com/martinsdevv/aegis/internal/health"
)

func NewRouter(healthCheck *health.Checker) http.Handler {
	mux := http.NewServeMux()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	prx, err := proxy.NewProxy(cfg.AegisUpstreamURL)
	if err != nil {
		http.Error(nil, "the proxy could not be created", http.StatusInternalServerError)
		return nil
	}

	mux.HandleFunc("/healthz", health.HealthHandler(healthCheck))
	mux.HandleFunc("/proxy/", proxy.HandleProxy(prx))
	mux.HandleFunc("/proxy", proxy.HandleProxy(prx))
	mux.HandleFunc("/panic", HandleNilPointer)

	var handler http.Handler = mux
	handler = middleware.NewMiddleware(handler)

	return handler
}
