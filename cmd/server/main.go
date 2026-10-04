package main

import (
	"log"
	"marketing/internal/app"
	"marketing/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed server load config %v", err)
	}

	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("failed to create application %v", err)
	}

	application.Run()
}
