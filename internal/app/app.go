package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"marketing/internal/infrastructure/database"
	"marketing/internal/infrastructure/telemetry"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"marketing/internal/config"
	transportHTTP "marketing/internal/transport/http"

	_ "github.com/jackc/pgx/v5"
)

type App struct {
	httpServer *http.Server
	pgConn     *database.Postgres
	logger     *telemetry.Logger
}

//nolint:exhaustruct_v5
func NewApp(
	ctx context.Context,
	cfg *config.Config,
) (*App, error) {
	pgConn, err := database.New(
		ctx,
		cfg.DB,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	log.Println("Database connection established successfully")

	customLogger := telemetry.NewLogger(cfg.Logger)

	router := transportHTTP.GetRouter(
		pgConn,
		customLogger,
		cfg.Service,
	)

	srv := &http.Server{
		Addr:              cfg.Server.Addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &App{
		httpServer: srv,
		pgConn:     pgConn,
		logger:     customLogger,
	}, nil
}

//nolint:contextcheck
func (a *App) Run(ctx context.Context) {
	runCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("Starting HTTP server on %s", a.httpServer.Addr)

		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	<-runCtx.Done()
	log.Println("Shutdown signal received. Initiating graceful shutdown...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var wg sync.WaitGroup

	wg.Go(
		func() {
			log.Println("Waiting for HTTP server to shutdown...")

			if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
				log.Printf("Error shutting down HTTP server: %v", err)
			}
		},
	)

	wg.Go(
		func() {
			log.Println("Waiting for database connection to close...")

			a.pgConn.Close()
		},
	)

	ch := make(chan struct{})

	go func() {
		wg.Wait()
		close(ch)
	}()

	select {
	case <-ch:
		log.Println("All shutdown tasks finished.")
	case <-shutdownCtx.Done():
		log.Println("Shutdown timeout exceeded!")
	}

	log.Println("Application stopped completely.")
}
