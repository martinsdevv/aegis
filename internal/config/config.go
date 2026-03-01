package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AEGIS_LISTEN_PORT  string
	AEGIS_REDIS_ADDR   string
	AEGIS_DATABASE_URL string
}

func Load() (Config, error) {
	_ = godotenv.Load()

	if err := RequireEnvs("AEGIS_LISTEN_PORT", "AEGIS_DATABASE_URL"); err != nil {
		return Config{}, err
	}
	cfg := Config{
		AEGIS_LISTEN_PORT:  getEnv("AEGIS_LISTEN_PORT", "8000"),
		AEGIS_REDIS_ADDR:   getEnv("AEGIS_REDIS_ADDR", "localhost:6379"),
		AEGIS_DATABASE_URL: getEnv("AEGIS_DATABASE_URL", ""),
	}

	return cfg, nil
}

func getEnv(key string, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func parseList(key string) []string {
	raw := os.Getenv(key)

	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}

	return result
}

func RequireEnvs(keys ...string) error {
	var missing []string
	for _, k := range keys {
		if v, ok := os.LookupEnv(k); !ok || strings.TrimSpace(v) == "" {
			missing = append(missing, k)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required env(s): %s", strings.Join(missing, ", "))
	}
	return nil
}
