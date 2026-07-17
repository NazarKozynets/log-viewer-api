package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	CORSOrigin string

	LogFiles     []string
	MaxReadBytes int64

	DefaultLimit int
	MaxLimit     int

	AuthMeURL    string
	AuthLoginURL string
	AuthCacheTTL int
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	logFiles := splitCSV(getEnv("LOG_FILES", ""))
	if len(logFiles) == 0 {
		return nil, fmt.Errorf("LOG_FILES is empty")
	}

	maxReadBytes, err := parseInt64(getEnv("MAX_READ_BYTES", "20971520"))
	if err != nil {
		return nil, fmt.Errorf("invalid MAX_READ_BYTES: %w", err)
	}

	defaultLimit, err := parseInt(getEnv("DEFAULT_LIMIT", "200"))
	if err != nil {
		return nil, fmt.Errorf("invalid DEFAULT_LIMIT: %w", err)
	}

	maxLimit, err := parseInt(getEnv("MAX_LIMIT", "1000"))
	if err != nil {
		return nil, fmt.Errorf("invalid MAX_LIMIT: %w", err)
	}

	if defaultLimit > maxLimit {
		return nil, fmt.Errorf("DEFAULT_LIMIT cannot be greater than MAX_LIMIT")
	}

	authMeURL := getEnv("AUTH_ME_URL", "")
	if authMeURL == "" {
		return nil, fmt.Errorf("AUTH_ME_URL is empty")
	}

	authLoginURL := getAuthLoginURL(authMeURL)
	if authLoginURL == "" {
		return nil, fmt.Errorf("AUTH_LOGIN_URL is empty")
	}

	authCacheTTL, err := parseInt(getEnv("AUTH_CACHE_TTL_SECONDS", "30"))
	if err != nil {
		return nil, fmt.Errorf("invalid AUTH_CACHE_TTL_SECONDS: %w", err)
	}

	return &Config{
		Port:       getEnv("PORT", "4001"),
		CORSOrigin: getEnv("CORS_ORIGIN", "http://localhost:5173"),

		LogFiles:     logFiles,
		MaxReadBytes: maxReadBytes,

		DefaultLimit: defaultLimit,
		MaxLimit:     maxLimit,

		AuthMeURL:    authMeURL,
		AuthLoginURL: authLoginURL,
		AuthCacheTTL: authCacheTTL,
	}, nil
}

func (c *Config) Addr() string {
	return ":" + c.Port
}

func getEnv(key, defaultValue string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}

	return value
}

func parseInt(value string) (int, error) {
	return strconv.Atoi(value)
}

func parseInt64(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}

func splitCSV(value string) []string {
	if value == "" {
		return []string{}
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item != "" {
			result = append(result, item)
		}
	}

	return result
}

func getAuthLoginURL(authMeURL string) string {
	if value := getEnv("AUTH_LOGIN_URL", ""); value != "" {
		return value
	}

	if mainAPIURL := strings.TrimRight(getEnv("MAIN_API_URL", ""), "/"); mainAPIURL != "" {
		return mainAPIURL + "/auth/login"
	}

	if strings.HasSuffix(authMeURL, "/auth/me") {
		return strings.TrimSuffix(authMeURL, "/auth/me") + "/auth/login"
	}

	return ""
}
