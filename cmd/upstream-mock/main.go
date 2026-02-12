package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/martinsdevv/aegis/internal/config"
	"github.com/martinsdevv/aegis/internal/health"
	mockhttp "github.com/martinsdevv/aegis/internal/upstreammock/upshttp"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	healthChecker := health.New()
	router := mockhttp.NewRouter(healthChecker)

	go func() {
		time.Sleep(2 * time.Second)
		healthChecker.SetReady()
	}()

	listenAddr := []string{":", cfg.AegisUpstreamPort}

	http.ListenAndServe(strings.Join(listenAddr, ""), router)
}
