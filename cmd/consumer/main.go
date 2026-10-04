package main

import (
	"fmt"
	"marketing/internal/app"
	"marketing/internal/config"
	"os"

	"github.com/spf13/cobra"
)

//nolint:exhaustruct_v5
func main() {
	var rootCmd = &cobra.Command{
		Use:   "consumer",
		Short: "Consumer tool",
	}

	var consumerCmd = &cobra.Command{
		Use:   "test",
		Short: "Test running",
		Run: func(cmd *cobra.Command, args []string) {
			var cfg *config.Config

			var err error

			var cliApp *app.CliApp

			cfg, err = config.LoadConfig()
			if err != nil {
				panic(fmt.Errorf("config init failed: %w", err))
			}

			cliApp, err = app.NewCliApp(cfg)
			if err != nil {
				panic(fmt.Errorf("cli app init failed: %w", err))
			}

			cliApp.Run()
		},
	}

	rootCmd.AddCommand(consumerCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
