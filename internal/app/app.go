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
	transportHTTP "marketing/internal/transport/http" // Импортируем наш HTTP транспорт
	"marketing/internal/transport/http/handler"
)

type App struct {
	httpServer  *http.Server
	kafkaReader *kafka.Reader
}

func NewApp(cfg *config.Config) *App {
	// 1. Инициализация Kafka Reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{cfg.Kafka.Broker},
		Topic:    cfg.Kafka.Topic,
		GroupID:  cfg.Kafka.GroupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	// 2. Инициализация HTTP хендлеров и роутера
	userHandler := handler.NewUserHandler()
	router := transportHTTP.NewRouter(userHandler) // Передаем хендлеры в роутер

	// 3. Настройка HTTP-сервера
	srv := &http.Server{
		Addr:              cfg.Server.Addr,
		Handler:           router, // Подставляем собранный роутер
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

	// Запуск Kafka консьюмера в фоне
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Запуск Kafka консьюмера...")
		for {
			msg, err := a.kafkaReader.ReadMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					log.Println("Прекращаем чтение из Kafka...")
					return
				}
				log.Printf("Ошибка чтения из Kafka: %v", err)
				continue
			}
			log.Printf("Получено сообщение: %s", string(msg.Value))
		}
	}()

	// Запуск HTTP сервера в фоне
	go func() {
		log.Printf("Запуск HTTP сервера на %s", a.httpServer.Addr)
		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Ошибка HTTP сервера: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Получен сигнал остановки. Инициируем Graceful Shutdown...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Плавный останов HTTP
	if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Ошибка при остановке HTTP: %v", err)
	}

	// Ожидание завершения обработки текущих сообщений Kafka
	log.Println("Ожидаем завершения работы воркеров Kafka...")
	ch := make(chan struct{})
	go func() {
		wg.Wait()
		close(ch)
	}()

	select {
	case <-ch:
		log.Println("Все воркеры Kafka успешно завершили работу.")
	case <-shutdownCtx.Done():
		log.Println("Превышен лимит времени ожидания воркеров Kafka!")
	}

	// Закрытие соединения с Kafka
	if err := a.kafkaReader.Close(); err != nil {
		log.Printf("Ошибка при закрытии Kafka: %v", err)
	}

	log.Println("Приложение полностью остановлено.")
}
