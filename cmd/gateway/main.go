package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/martinsdevv/aegis/internal/config"
	"github.com/martinsdevv/aegis/internal/gateway/gtwhttp"
	"github.com/martinsdevv/aegis/internal/health"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	healthCheck := health.New()
	router := gtwhttp.NewRouter(healthCheck)

	go func() {
		time.Sleep(time.Second * 2)
		healthCheck.SetReady()
	}()

	listenAddr := []string{":", cfg.AegisListenPort}
	http.ListenAndServe(strings.Join(listenAddr, ""), router)
}
