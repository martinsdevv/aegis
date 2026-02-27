package main

import (
	"log"
	"net/http"
	"time"

	"github.com/martinsdevv/aegis/internal/config"
	"github.com/martinsdevv/aegis/internal/gateway/gtwhttp"
	"github.com/martinsdevv/aegis/internal/gateway/middleware"
	"github.com/martinsdevv/aegis/internal/health"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	healthCheck := health.New()
	store := middleware.NewRLStore(5, 10, 30*time.Minute)
	redisClient := middleware.NewRedisClient(cfg.AegisRedisAddr)

	router := gtwhttp.NewRouter(healthCheck, cfg, store, redisClient)

	go func() {
		t := time.NewTicker(5 * time.Minute)
		defer t.Stop()
		for range t.C {
			store.Cleanup()
		}
	}()

	go func() {
		time.Sleep(2 * time.Second)
		healthCheck.SetReady()
	}()

	addr := ":" + cfg.AegisListenPort
	log.Printf("Aegis listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
