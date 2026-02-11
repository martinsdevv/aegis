package config

import (
	"os"
	"strings"
	"github.com/joho/godotenv"
	"fmt"
)

type Config struct {
	AEGIS_LISTEN_ADDR string
	AEGIS_UPSTREAM_URL string
	AEGIS_API_KEYS []string
}

func Load() (Config, error) {
	_ = godotenv.Load()
	
	if err := RequireEnvs("AEGIS_UPSTREAM_URL", "AEGIS_API_KEYS"); err != nil {
		return Config{}, err
	}
	cfg := Config {
		AEGIS_LISTEN_ADDR: getEnv("AEGIS_LISTEN_ADDR", "8000"),
		AEGIS_UPSTREAM_URL: getEnv("AEGIS_UPSTREAM_URL", "9000"),
		AEGIS_API_KEYS: parseList("AEGIS_API_KEYS"),
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
	if len (missing) > 0 {
		return fmt.Errorf("missing required env(s): %s", strings.Join(missing, ", "))
	}
	return nil
}
