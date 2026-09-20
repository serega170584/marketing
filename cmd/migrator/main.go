package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source"
	_ "github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/spf13/cobra"

	"marketing/internal/config"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "migrator",
		Short: "Database migration CLI tool",
	}

	var upCmd = &cobra.Command{
		Use:   "up",
		Short: "Apply all pending migrations",
		Run: func(cmd *cobra.Command, args []string) {
			m := initMigrator()
			defer m.Close()

			log.Println("Applying database migrations...")
			if err := m.Up(); err != nil {
				if errors.Is(err, migrate.ErrNoChange) {
					log.Println("Database schema is already up to date.")
					return
				}
				log.Fatalf("Failed to apply migrations: %v", err)
			}
			log.Println("Migrations applied successfully!")
		},
	}

	var downCmd = &cobra.Command{
		Use:   "down",
		Short: "Roll back the last migration",
		Run: func(cmd *cobra.Command, args []string) {
			m := initMigrator()
			defer m.Close()

			log.Println("Rolling back the last database migration...")
			if err := m.Down(); err != nil {
				if errors.Is(err, migrate.ErrNoChange) {
					log.Println("No migrations to roll back.")
					return
				}
				log.Fatalf("Failed to roll back migrations: %v", err)
			}
			log.Println("Migration rolled back successfully!")
		},
	}

	var createCmd = &cobra.Command{
		Use:   "create [migration_name]",
		Short: "Create a new pair of up and down migration files",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			dir := "migrations"

			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				log.Fatalf("Failed to create migrations directory: %v", err)
			}

			timestamp := time.Now().Format("20060102150405")
			upFile := filepath.Join(dir, fmt.Sprintf("%s_%s.up.sql", timestamp, name))
			downFile := filepath.Join(dir, fmt.Sprintf("%s_%s.down.sql", timestamp, name))

			createEmptyFile(upFile)
			createEmptyFile(downFile)

			log.Printf("Created migration files:\n  %s\n  %s\n", upFile, downFile)
		},
	}

	rootCmd.AddCommand(upCmd, downCmd, createCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func initMigrator() *migrate.Migrate {
	cfg := config.LoadConfig()
	if cfg.DB.URL == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	m, err := migrate.New("file://migrations", cfg.DB.URL)
	if err != nil {
		log.Fatalf("Failed to initialize migration engine: %v", err)
	}

	return m
}

func createEmptyFile(path string) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("Failed to create file %s: %v", path, err)
	}
	_ = file.Close()
}
