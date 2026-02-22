package infra

import (
	"os"
	"strconv"
	"strings"
)

// LoadEnv reads .env from path and sets env vars. Skips if file not found.
func LoadEnv(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.Index(line, "=")
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		val = strings.Trim(val, "\"")
		if key != "" && os.Getenv(key) == "" {
			os.Setenv(key, val)
		}
	}
}

type Config struct {
	HTTPPort        int
	CacheTTLSeconds int
	DB              DBConfig
}

// db context url holding the connection
type DBConfig struct {
	URL string
}

// GatewayConfig holds gateway/load-balancer config from env.
type GatewayConfig struct {
	Port       int
	BackendURLs string
}

// Load config from env variables (server).
func Load() *Config {
	return &Config{
		HTTPPort:        getEnvInt("HTTP_PORT", 8080),
		CacheTTLSeconds: getEnvInt("CACHE_TTL_SECONDS", 60),
		DB: DBConfig{
			URL: getEnv("DATABASE_URL", "postgres://eventflow:eventflow@localhost:5432/eventflow?sslmode=disable"),
		},
	}
}

// LoadGateway loads gateway config from env variables.
func LoadGateway() *GatewayConfig {
	urls := getEnv("BACKEND_URLS", "http://localhost:8081,http://localhost:8082,http://localhost:8083")
	port := getEnvInt("GATEWAY_PORT", 0)
	if port == 0 {
		port = getEnvInt("HTTP_PORT", 8080)
	}
	return &GatewayConfig{
		Port:        port,
		BackendURLs: urls,
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
