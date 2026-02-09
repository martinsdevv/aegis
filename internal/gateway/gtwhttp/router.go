package gtwhttp

import (
	"net/http"

	"github.com/martinsdevv/aegis/internal/gateway/proxy"
	"github.com/martinsdevv/aegis/internal/health"
)

func NewRouter(healthCheck *health.Checker) *http.ServeMux {
	mux := http.NewServeMux()
	prx, err := proxy.NewProxy("http://localhost:9000")
	if err != nil {
		panic(err)
	}

	mux.HandleFunc("/healthz", health.HealthHandler(healthCheck))
	mux.HandleFunc("/proxy/", proxy.HandleProxy(prx))
	mux.HandleFunc("/proxy", proxy.HandleProxy(prx))

	return mux
}
