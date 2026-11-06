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
	tracer   *telemetry.Tracer
	metrics  *telemetry.Metrics
}

func NewCliApp(
	ctx context.Context,
	cfg *config.Config,
) (
	*CliApp,
	error,
) {
	var err error

	var consumer *kafka.Consumer

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

	consumer, err = kafka.NewConsumer(
		cfg.KafkaConsumer,
		tracer,
	)
	if err != nil {
		return nil, fmt.Errorf("consumer init failed: %w", err)
	}

	return &CliApp{
		consumer: consumer,
		tracer:   tracer,
		metrics:  metrics,
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

	wg.Go(func() {
		log.Println("Waiting for tracer to finish...")

		err = a.tracer.Shutdown(shutdownCtx)
		if err != nil {
			log.Printf("tracer shutdown failed: %v", err)
		}
	})

	wg.Go(func() {
		log.Println("Waiting for metrics to finish...")

		err = a.metrics.Shutdown(shutdownCtx)
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
