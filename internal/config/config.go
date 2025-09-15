// internal/config/config.go
package config

import (
	"os"
	"time"
)

type Config struct {
	APIKey      string
	Language    string
	Days        int
	Timeout     time.Duration
	EnableCache bool
	CacheTTL    time.Duration
}

func FromEnv() Config {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		// you may decide to return empty and error later if not provided via flag
	}
	return Config{
		APIKey:      apiKey,
		Language:    "es",
		Days:        1,
		Timeout:     10 * time.Second,
		EnableCache: true,
		CacheTTL:    10 * time.Minute,
	}
}
