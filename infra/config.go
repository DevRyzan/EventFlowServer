package infra

import (
	"os"
	"strconv"
)

type Config struct {
	HTTPPort      int
	CacheTTLSeconds int
	DB            DBConfig
}

// db context url holding the connection
type DBConfig struct {
	URL string
}

// load config from env variables
func Load() *Config {
	return &Config{
		HTTPPort:        getEnvInt("HTTP_PORT", 8080),
		CacheTTLSeconds: getEnvInt("CACHE_TTL_SECONDS", 60),
		DB: DBConfig{
			URL: getEnv("DATABASE_URL", "postgres://eventflow:eventflow@localhost:5432/eventflow?sslmode=disable"),
		},
	}
}

// get env variable with default value
func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

// get env variable as int with default value
func getEnvInt(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return defaultVal
}
