package app

import (
	"context"
	"fmt"
	"log"
	"marketing/internal/config"
	"marketing/internal/infrastructure/telemetry"
	"marketing/internal/transport/kafka"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type CliApp struct {
	consumer *kafka.Consumer
	logger   *telemetry.Logger
}

func NewCliApp(cfg *config.Config) (*CliApp, error) {
	var err error

	var customLogger *telemetry.Logger

	var consumer *kafka.Consumer

	customLogger = telemetry.NewLogger(cfg.Logger)

	consumer, err = kafka.NewConsumer(
		cfg.KafkaConsumer,
		customLogger,
	)
	if err != nil {
		return nil, fmt.Errorf("consumer init failed: %w", err)
	}

	return &CliApp{
		consumer: consumer,
		logger:   customLogger,
	}, nil
}

func (a *CliApp) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var err error

	err = a.consumer.Start(ctx)
	if err != nil {
		panic(fmt.Errorf("consumer start failed: %w", err))
	}

	var wg sync.WaitGroup

	<-ctx.Done()
	log.Println("Shutdown signal received. Initiating graceful shutdown...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	wg.Go(func() {
		log.Println("Waiting for Kafka consumers to finish...")

		err = a.consumer.Close(shutdownCtx)
		if err != nil {
			log.Printf("consumer close failed: %v", err)
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
