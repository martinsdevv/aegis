package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/martinsdevv/aegis/internal/config"
	"github.com/martinsdevv/aegis/internal/gateway/gtwhttp"
	"github.com/martinsdevv/aegis/internal/gateway/middleware"
	"github.com/martinsdevv/aegis/internal/health"
	"github.com/martinsdevv/aegis/internal/seed"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// DB
	db, err := sql.Open("pgx", cfg.AEGIS_DATABASE_URL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	if err := seed.RunMigrations(cfg.AEGIS_DATABASE_URL); err != nil {
		log.Fatal(err)
	}

	if err := seed.RunSeed(ctx, db); err != nil {
		log.Fatal(err)
	}

	healthCheck := health.New()
	store := middleware.NewRLStore(5, 10, 30*time.Minute)
	redisClient := middleware.NewRedisClient(cfg.AEGIS_REDIS_ADDR)
	apiKeyStore := middleware.NewAPIKeyStore(db, redisClient, 60*time.Second)

	router := gtwhttp.NewRouter(healthCheck, cfg, store, redisClient, apiKeyStore)

	server := &http.Server{
		Addr:    ":" + cfg.AEGIS_LISTEN_PORT,
		Handler: router,
	}

	// Cleanup goroutine
	go func() {
		t := time.NewTicker(5 * time.Minute)
		defer t.Stop()
		for range t.C {
			store.Cleanup()
		}
	}()

	// Ready after boot
	go func() {
		time.Sleep(2 * time.Second)
		healthCheck.SetReady()
	}()

	// Start server
	go func() {
		log.Printf("Aegis listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	log.Println("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown failed: %v\n", err)
	}

	log.Println("Aegis stopped gracefully")
}
