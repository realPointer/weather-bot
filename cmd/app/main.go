package main

import (
	"log"

	"github.com/realPointer/weather-bot/config"
	"github.com/realPointer/weather-bot/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
