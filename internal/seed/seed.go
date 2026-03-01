package seed

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"log"

	"github.com/golang-migrate/migrate/v4"
)

func RunSeed(ctx context.Context, db *sql.DB) error {
	type Seed struct {
		Name, RawKey, Upstream string
		Quota                  int
	}

	seeds := []Seed{
		{"default-dev", "DEV_KEY_123", "https://httpbin.org", 10000},
		{"internal-test", "TEST_KEY_456", "https://postman-echo.com", 5000},
	}

	for _, s := range seeds {
		hashed := hashKey(s.RawKey)
		_, err := db.ExecContext(ctx, `
			INSERT INTO api_keys (name, key, upstream_host, monthly_quota, is_active)
			VALUES ($1, $2, $3, $4, TRUE)
			ON CONFLICT (key) DO NOTHING;
		`, s.Name, hashed, s.Upstream, s.Quota)
		if err != nil {
			return err
		}
	}

	log.Println("Seed executed")
	return nil
}

func RunMigrations(databaseURL string) error {
	m, err := migrate.New(
		"file://internal/db/migrations",
		databaseURL,
	)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Database migrations applied")
	return nil
}

func hashKey(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
