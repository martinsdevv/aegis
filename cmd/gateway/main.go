package main

import (
	"net/http"
	"time"

	"github.com/martinsdevv/aegis/internal/gateway/gtwhttp"
	"github.com/martinsdevv/aegis/internal/health"
)

func main() {
	healthCheck := health.New()
	router := gtwhttp.NewRouter(healthCheck)

	go func() {
		time.Sleep(time.Second * 2)
		healthCheck.SetReady()
	}()
	http.ListenAndServe(":8000", router)
}
