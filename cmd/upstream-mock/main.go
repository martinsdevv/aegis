package main

import (
	"net/http"
	"strings"
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

	listenAddr := []string{":", "9000"}

	http.ListenAndServe(strings.Join(listenAddr, ""), router)
}
