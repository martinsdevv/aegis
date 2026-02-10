package gtwhttp

import (
	"net/http"

	"github.com/martinsdevv/aegis/internal/gateway/middleware"
	"github.com/martinsdevv/aegis/internal/gateway/proxy"
	"github.com/martinsdevv/aegis/internal/health"
)

func NewRouter(healthCheck *health.Checker) http.Handler {
	mux := http.NewServeMux()
	prx, err := proxy.NewProxy("http://localhost:9000")
	if err != nil {
		http.Error(nil, "the proxy could not be created", http.StatusInternalServerError)
		return nil
	}

	mux.HandleFunc("/healthz", health.HealthHandler(healthCheck))
	mux.HandleFunc("/proxy/", proxy.HandleProxy(prx))
	mux.HandleFunc("/proxy", proxy.HandleProxy(prx))

	var handler http.Handler = mux
	handler = middleware.Logger(handler)

	return handler
}
