package main

import (
	"marketing/internal/app"
	"marketing/internal/config"
)

func main() {
	// 1. Загружаем конфигурацию
	cfg := config.LoadConfig()

	// 2. Передаем её при создании инстанса приложения
	application := app.NewApp(cfg)

	// 3. Запускаем
	application.Run()
}
