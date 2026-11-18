package config

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

type Config struct {
	Port            string
	RedisAddr       string
	RedisPassword   string
	RedisDB         int
	CORSOrigin      string
	TypesafeAPIKey  string
	TavilyAPIKey    string
	FirecrawlAPIKey string
}

func LoadConfig() (*Config, error) {
	loadDotEnv(".env")

	cfg := &Config{
		Port:            getEnv("PORT", DefaultPort),
		RedisAddr:       getEnv("REDIS_ADDR", DefaultRedisAddr),
		RedisPassword:   getEnv("REDIS_PASSWORD", ""),
		RedisDB:         0,
		CORSOrigin:      getEnv("CORS_ORIGIN", "*"),
		TypesafeAPIKey:  getEnv("TYPESAFE_API_KEY", ""),
		TavilyAPIKey:    getEnv("TAVILY_API_KEY", ""),
		FirecrawlAPIKey: getEnv("FIRECRAWL_API_KEY", ""),
	}

	if cfg.TypesafeAPIKey == "" {
		return nil, errors.New("TYPESAFE_API_KEY is required")
	}

	if cfg.TavilyAPIKey == "" && cfg.FirecrawlAPIKey == "" {
		return nil, errors.New("at least one search API key (TAVILY_API_KEY or FIRECRAWL_API_KEY) must be provided")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func loadDotEnv(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			k := strings.TrimSpace(parts[0])
			v := strings.TrimSpace(parts[1])
			v = strings.Trim(v, `"'`)
			if os.Getenv(k) == "" {
				os.Setenv(k, v)
			}
		}
	}
}
