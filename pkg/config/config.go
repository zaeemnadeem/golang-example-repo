package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// Config holds the configuration for both services
type Config struct {
	Env                string `envconfig:"ENV" default:"development"`
	ScreenPort         int    `envconfig:"SCREEN_PORT" default:"8001"` // Screen changed to HTTP standard port
	ContentPort        int    `envconfig:"CONTENT_PORT" default:"50052"`
	ContentServiceAddr string `envconfig:"CONTENT_SERVICE_ADDR" default:"localhost:50052"`
	DatabaseURL        string `envconfig:"DATABASE_URL" default:"postgres://postgres:postgres@localhost:5432/signage?sslmode=disable"`
}

// Load reads config from environment variables
func Load() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Printf("Failed to process envconfig: %v", err)
		return nil, err
	}
	return &cfg, nil
}
