package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App `yaml:"app"`
		Log `yaml:"logger"`
	}

	// App -.
	App struct {
		Name          string `yaml:"name"`
		Version       string `yaml:"version"`
		TelegramToken string `yaml:"telegram_token"`
		WeatherToken  string `yaml:"weather_token"`
	}

	// Log -.
	Log struct {
		Level string `yaml:"log_level"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
