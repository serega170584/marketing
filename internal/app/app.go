package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"

	"marketing/internal/config"
	transportHTTP "marketing/internal/transport/http"
	"marketing/internal/transport/http/handler"
)

type App struct {
	httpServer  *http.Server
	kafkaReader *kafka.Reader
}

//nolint:exhaustruct_v5
func NewApp(cfg *config.Config) *App {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{cfg.Kafka.Broker},
		Topic:    cfg.Kafka.Topic,
		GroupID:  cfg.Kafka.GroupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	userHandler := handler.NewUserHandler()
	router := transportHTTP.NewRouter(userHandler)

	srv := &http.Server{
		Addr:              cfg.Server.Addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &App{
		httpServer:  srv,
		kafkaReader: reader,
	}
}

func (a *App) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	wg.Go(func() {
		log.Println("Starting Kafka consumer...")

		for {
			msg, err := a.kafkaReader.ReadMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					log.Println("Stopping Kafka message consumption...")
					return
				}

				log.Printf("Error reading from Kafka: %v", err)

				continue
			}

			log.Printf("Received message: %s", string(msg.Value))
		}
	})

	go func() {
		log.Printf("Starting HTTP server on %s", a.httpServer.Addr)

		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutdown signal received. Initiating graceful shutdown...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error shutting down HTTP server: %v", err)
	}

	log.Println("Waiting for Kafka workers to finish...")

	ch := make(chan struct{})

	go func() {
		wg.Wait()
		close(ch)
	}()

	select {
	case <-ch:
		log.Println("All Kafka workers finished successfully.")
	case <-shutdownCtx.Done():
		log.Println("Kafka workers shutdown timeout exceeded!")
	}

	if err := a.kafkaReader.Close(); err != nil {
		log.Printf("Error closing Kafka reader: %v", err)
	}

	log.Println("Application stopped completely.")
}
