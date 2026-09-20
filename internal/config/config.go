package config

import (
	"os"
)

// Config содержит все настройки приложения
type Config struct {
	Server ServerConfig
	Kafka  KafkaConfig
}

type ServerConfig struct {
	Addr string
}

type KafkaConfig struct {
	Broker  string
	Topic   string
	GroupID string
}

// LoadConfig читает переменные окружения и возвращает заполненную структуру
func LoadConfig() *Config {
	// В реальном проекте здесь обычно вызывают godotenv.Load(),
	// чтобы подгрузить переменные из .env файла.

	return &Config{
		Server: ServerConfig{
			Addr: getEnv("SERVER_ADDR", ":8080"),
		},
		Kafka: KafkaConfig{
			Broker:  getEnv("KAFKA_BROKER", "localhost:9092"),
			Topic:   getEnv("KAFKA_TOPIC", "user-events"),
			GroupID: getEnv("KAFKA_GROUP_ID", "my-service-group"),
		},
	}
}

// Хелпер для установки дефолтных значений, если переменная пустая
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
