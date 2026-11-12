package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"marketing/internal/auth"
	"marketing/internal/infrastructure/database"
	"marketing/internal/infrastructure/telemetry"
	"marketing/internal/infrastructure/telemetry/logger"
	"marketing/internal/repository"
	"marketing/internal/service"
	"marketing/internal/usecase"
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
	tracer     *telemetry.Tracer
	metrics    *telemetry.Metrics
}

//nolint:exhaustruct_v5
func NewApp(
	ctx context.Context,
	cfg *config.Config,
) (*App, error) {
	err := auth.LoadPublicKey("public.pem")
	if err != nil {
		return nil, fmt.Errorf("load public key error: %w", err)
	}

	pgConn, err := database.New(
		ctx,
		cfg.DB,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	log.Println("Database connection established successfully")

	err = logger.InitLogger(cfg.Logger)
	if err != nil {
		return nil, fmt.Errorf(
			"logger init failed: %w",
			err,
		)
	}

	tracer, err := telemetry.NewTracer(
		ctx,
		cfg.Tracer,
	)
	if err != nil {
		return nil, fmt.Errorf("tracer init failed: %w", err)
	}

	metrics, err := telemetry.NewMetrics(
		ctx,
		cfg.Metrics.ServiceName,
		cfg.Metrics.ServiceVersion,
		cfg.Metrics.Address,
	)
	if err != nil {
		return nil, fmt.Errorf("metrics init failed: %w", err)
	}

	promoCodeRepo := repository.NewPromoCodeRepo(pgConn)
	promoCodeService := service.NewPromoCodeService(promoCodeRepo)
	promoCodeUseCase := usecase.NewPromoCodeUseCase(promoCodeService)

	router := transportHTTP.GetRouter(
		pgConn,
		promoCodeUseCase,
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
		tracer:     tracer,
		metrics:    metrics,
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

	wg.Go(func() {
		log.Println("Waiting for tracer to finish...")

		err := a.tracer.Shutdown(shutdownCtx)
		if err != nil {
			log.Printf("tracer shutdown failed: %v", err)
		}
	})

	wg.Go(func() {
		log.Println("Waiting for metrics to finish...")

		err := a.metrics.Shutdown(shutdownCtx)
		if err != nil {
			log.Printf("metrics shutdown failed: %v", err)
		}
	})

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
