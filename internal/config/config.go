package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AegisListenPort   string
	AegisUpstreamURL  string
	AegisAPIKeys      []string
	AegisUpstreamPort string
	AegisRedisAddr    string
}

func Load() (Config, error) {
	_ = godotenv.Load()

	if err := RequireEnvs("AegisUpstreamURL", "AegisAPIKeys", "AegisListenPort", "AegisUpstreamPort"); err != nil {
		return Config{}, err
	}
	cfg := Config{
		AegisListenPort:   getEnv("AegisListenPort", "8000"),
		AegisUpstreamURL:  getEnv("AegisUpstreamURL", "http://localhost:9000"),
		AegisUpstreamPort: getEnv("AegisUpstreamPort", "9000"),
		AegisAPIKeys:      parseList("AegisAPIKeys"),
		AegisRedisAddr:    getEnv("AegisRedisAddr", "localhost:6379"),
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
