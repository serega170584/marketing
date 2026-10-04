package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server        ServerConfig
	KafkaConsumer KafkaConsumerConfig
	DB            DBConfig
	Logger        LoggerConfig
}

type ServerConfig struct {
	Addr string
}

type KafkaConfig struct {
	Broker  string
	Topic   string
	GroupID string
}

type KafkaConsumerConfig struct {
	PingMetadataTimeout   time.Duration
	AutoOffsetReset       string
	EnableAutoOffsetStore bool
	ReadTimeout           time.Duration
	Brokers               []string
	Topic                 string
	GroupID               string
}

type DBConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type LoggerConfig struct {
	Env string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	maxOpenConns, err := parseInt(getEnv("MAX_OPEN_CONNS", "10"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse MAX_OPEN_CONNS: %w", err)
	}

	maxIdleConns, err := parseInt(getEnv("MAX_IDLE_CONNS", "10"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse MAX_IDLE_CONNS: %w", err)
	}

	connMaxLifetime, err := parseDuration(getEnv("CONN_MAX_LIFETIME", "15m"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse CONN_MAX_LIFETIME: %w", err)
	}

	pingMetadataTimeout, err := parseDuration(getEnv("KAFKA_CONSUMER_PING_METADATA_TIMEOUT", "5s"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse KAFKA_CONSUMER_PING_METADATA_TIMEOUT: %w", err)
	}

	enableAutoOffsetStore, err := strconv.ParseBool(getEnv("KAFKA_CONSUMER_ENABLE_AUTO_OFFSET_STORE", "true"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse KAFKA_CONSUMER_ENABLE_AUTO_OFFSET_STORE: %w", err)
	}

	readTimeout, err := parseDuration(getEnv("KAFKA_CONSUMER_READ_TIMEOUT", "5s"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse KAFKA_CONSUMER_READ_TIMEOUT: %w", err)
	}

	kafkaBrokers := strings.Split(
		getEnv("KAFKA_CONSUMER_BROKERS", "localhost:9092"),
		",",
	)

	return &Config{
		Server: ServerConfig{
			Addr: getEnv("SERVER_ADDR", ":8080"),
		},
		KafkaConsumer: KafkaConsumerConfig{
			Brokers:               kafkaBrokers,
			Topic:                 getEnv("KAFKA_CONSUMER_TOPIC", "user-events"),
			GroupID:               getEnv("KAFKA_CONSUMER_GROUP_ID", "my-service-group"),
			PingMetadataTimeout:   pingMetadataTimeout,
			AutoOffsetReset:       getEnv("KAFKA_CONSUMER_AUTO_OFFSET_RESET", "earliest"),
			EnableAutoOffsetStore: enableAutoOffsetStore,
			ReadTimeout:           readTimeout,
		},
		DB: DBConfig{
			URL:             getEnv("DB_URL", "postgres://postgres:postgres@localhost:5432/marketing?sslmode=disable"),
			MaxOpenConns:    maxOpenConns,
			MaxIdleConns:    maxIdleConns,
			ConnMaxLifetime: connMaxLifetime,
		},
		Logger: LoggerConfig{
			Env: getEnv("LOGGER_ENV", "dev"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func parseDuration(duration string) (time.Duration, error) {
	parsedDuration, err := time.ParseDuration(duration)
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration: %w", err)
	}

	return parsedDuration, nil
}

func parseInt(value string) (int, error) {
	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("failed to parse int: %w", err)
	}

	return parsedValue, nil
}
