package kafka

import (
	"context"
	"fmt"
	"marketing/internal/config"
	"marketing/internal/infrastructure/logger"
	"marketing/internal/transport/kafka/message"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Consumer struct {
	consumer    *kafka.Consumer
	readTimeout time.Duration
	logger      *logger.Logger
	metadata    *kafka.Metadata
	router      *Router
}

func NewConsumer(
	cfg config.KafkaConsumerConfig,
	logger *logger.Logger,
) (*Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        cfg.Brokers,
		"group.id":                 cfg.GroupID,
		"auto.offset.reset":        cfg.AutoOffsetReset,
		"enable.auto.offset.store": cfg.EnableAutoOffsetStore,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka consumer: %w", err)
	}

	metadata, metadataErr := c.GetMetadata(
		nil,
		false,
		int(cfg.PingMetadataTimeout.Milliseconds()),
	)
	if metadataErr != nil {
		return nil, fmt.Errorf("kafka broker is unavailable: %w", metadataErr)
	}

	return &Consumer{
		consumer:    c,
		readTimeout: cfg.ReadTimeout,
		logger:      logger,
		metadata:    metadata,
		router:      getRouter(logger),
	}, nil
}

func (c *Consumer) Start(
	ctx context.Context,
) error {
	var err error

	topics := c.router.GetTopics()

	err = c.consumer.SubscribeTopics(
		topics,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to subscribe to topics: %w", err)
	}

	c.logger.Info(
		ctx,
		"subscribed to topics. Starting poll loop...",
		logger.Strings("topics", topics),
	)

	for {
		select {
		case <-ctx.Done():
			c.logger.Info(
				ctx,
				"received shutdown signal, stopping poll loop...",
			)

			return nil
		default:
			msg, msgErr := c.consumer.ReadMessage(
				c.readTimeout,
			)
			if msgErr != nil {
				c.logger.Error(
					ctx,
					"message getting error",
					logger.Error(msgErr),
				)

				continue
			}

			var topicName string

			if msg.TopicPartition.Topic != nil {
				topicName = *msg.TopicPartition.Topic
			}

			internalMsg := message.Message{
				Topic:   topicName,
				Payload: msg.Value,
			}

			err = c.router.Route(ctx, internalMsg)
			if err != nil {
				c.logger.Error(ctx, "message routing error", logger.Error(err))
			}
		}
	}
}

func (c *Consumer) Close(ctx context.Context) error {
	c.logger.Info(
		ctx,
		"closing consumer and committing offsets...",
	)

	return c.consumer.Close()
}
