// internal/config/config.go
package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
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

	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found or failed to load")
	}

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
