package main

import (
	"log"
	"net/http"
	"strings"
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
	router := gtwhttp.NewRouter(healthCheck, cfg, store)

	go func() {
		t := time.NewTicker(5 * time.Minute)
		defer t.Stop()
		for range t.C {
			store.Cleanup(time.Now())
		}
	}()

	go func() {
		time.Sleep(time.Second * 2)
		healthCheck.SetReady()
	}()

	listenAddr := []string{":", cfg.AegisListenPort}
	addr := strings.Join(listenAddr, "")
	log.Fatal(http.ListenAndServe(addr, router))
}
