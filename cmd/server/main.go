package main

import (
	"marketing/internal/app"
	"marketing/internal/config"
)

func main() {
	cfg := config.LoadConfig()
	application := app.NewApp(cfg)
	application.Run()
}
