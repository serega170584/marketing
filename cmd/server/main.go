package main

import (
	"context"
	"log"
	"marketing/internal/app"
	"marketing/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed server load config %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application, err := app.NewApp(
		ctx,
		cfg,
	)
	if err != nil {
		log.Printf("failed to create application %v", err)

		return
	}

	application.Run(ctx)
}
