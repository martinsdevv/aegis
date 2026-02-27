package gtwhttp

import (
	"log"
	"net/http"

	"github.com/martinsdevv/aegis/internal/config"
	"github.com/martinsdevv/aegis/internal/gateway/middleware"
	"github.com/martinsdevv/aegis/internal/gateway/proxy"
	"github.com/martinsdevv/aegis/internal/health"
	"github.com/redis/go-redis/v9"
)

func NewRouter(healthCheck *health.Checker, cfg config.Config, store *middleware.RLStore, redisClient *redis.Client) http.Handler {
	mux := http.NewServeMux()

	prx, err := proxy.NewProxy(cfg.AegisUpstreamURL)
	if err != nil {
		log.Fatal(err)
	}

	mux.HandleFunc("/healthz", health.HealthHandler(healthCheck))
	mux.HandleFunc("/proxy/", proxy.HandleProxy(prx))
	mux.HandleFunc("/proxy", proxy.HandleProxy(prx))
	mux.HandleFunc("/panic", HandleNilPointer)
	mux.HandleFunc("/rltest", HandleRLTest)

	quotaMgr := middleware.NewQuotaManager(redisClient, 2)

	var handler http.Handler = mux
	handler = middleware.NewMiddleware(handler, cfg, store, quotaMgr, redisClient)

	return handler
}
