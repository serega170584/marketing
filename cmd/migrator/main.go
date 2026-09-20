package main

import (
	"errors"
	"flag"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source"
	_ "github.com/golang-migrate/migrate/v4/source/iofs"

	"marketing/internal/config"
)

//nolint:gocritic
func main() {
	var action string
	flag.StringVar(&action, "action", "up", "Migration action: up or down")
	flag.Parse()

	cfg := config.LoadConfig()
	if cfg.DB.URL == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	m, err := migrate.New("file://migrations", cfg.DB.URL)
	if err != nil {
		log.Fatalf("Failed to initialize migration engine: %v", err)
	}

	defer func(m *migrate.Migrate) {
		sourceErr, dbErr := m.Close()
		if sourceErr != nil {
			log.Printf("migrarion close error %s", sourceErr.Error())
		}

		if dbErr != nil {
			log.Printf("migrarion close error %s", sourceErr.Error())
		}
	}(m)

	switch action {
	case "up":
		log.Println("Applying database migrations...")

		if err := m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				log.Println("Database schema is already up to date.")
				return
			}

			log.Fatalf("Failed to apply migrations: %v", err)
		}

		log.Println("Migrations applied successfully!")
	case "down":
		log.Println("Rolling back the last database migration...")

		if err := m.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				log.Println("No migrations to roll back.")

				return
			}

			log.Fatalf("Failed to roll back migrations: %v", err)
		}

		log.Println("Migration rolled back successfully!")

	default:
		log.Fatalf("Unknown action: %s. Use 'up' or 'down'", action)
	}
}
