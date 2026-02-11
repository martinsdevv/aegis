package main

import (
	"net/http"
	"time"
	"log"

	"github.com/martinsdevv/aegis/internal/gateway/gtwhttp"
	"github.com/martinsdevv/aegis/internal/health"
	"github.com/martinsdevv/aegis/internal/config"
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		log.Fatal(err)
	}

	for _, keys := range cfg.AEGIS_API_KEYS {
		println(keys)
	}
	healthCheck := health.New()
	router := gtwhttp.NewRouter(healthCheck)

	go func() {
		time.Sleep(time.Second * 2)
		healthCheck.SetReady()
	}()
	http.ListenAndServe(":8000", router)
}
