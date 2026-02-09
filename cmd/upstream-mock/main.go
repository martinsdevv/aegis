package main

import (
	"net/http"
	"time"

	"github.com/martinsdevv/aegis/internal/health"
	mockhttp "github.com/martinsdevv/aegis/internal/upstreammock/upshttp"
)

func main() {
	healthChecker := health.New()
	router := mockhttp.NewRouter(healthChecker)

	go func() {
		time.Sleep(2 * time.Second)
		healthChecker.SetReady()
	}()

	http.ListenAndServe(":9000", router)
}
