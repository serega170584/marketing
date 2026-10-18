package kafka

import (
	"context"
	"fmt"
	"marketing/internal/transport/kafka/message"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"marketing/internal/config"
)

type Consumer struct {
	consumer    *kafka.Consumer
	tracer      Tracer
	readTimeout time.Duration
	logger      Logger
	metadata    *kafka.Metadata
	router      *Router
}

func NewConsumer(
	cfg config.KafkaConsumerConfig,
	logger Logger,
	tracer Tracer,
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
		tracer:      tracer,
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
		map[string]any{
			"topics": topics,
		},
	)

	for {
		select {
		case <-ctx.Done():
			c.logger.Info(
				ctx,
				"received shutdown signal, stopping poll loop...",
				nil,
			)

			return nil
		default:
			func() {
				msg, msgErr := c.consumer.ReadMessage(
					c.readTimeout,
				)
				if msgErr != nil {
					c.logger.Error(
						ctx,
						"message getting error",
						map[string]any{
							"error": msgErr,
						},
					)

					return
				}

				consumerCtx, cancel := context.WithCancel(ctx)

				consumerCtx = otel.GetTextMapPropagator().Extract(consumerCtx, &KafkaConsumerHeadersCarrier{Headers: msg.Headers})

				consumerCtx, span := c.tracer.Start(
					consumerCtx,
					fmt.Sprintf("Kafka CONSUME from %s", *msg.TopicPartition.Topic),
					trace.WithSpanKind(trace.SpanKindConsumer),
				)
				defer span.End()

				defer cancel()

				var topicName string

				if msg.TopicPartition.Topic != nil {
					topicName = *msg.TopicPartition.Topic
				}

				internalMsg := message.Message{
					Topic:   topicName,
					Payload: msg.Value,
				}

				err = c.router.Route(consumerCtx, internalMsg)
				if err != nil {
					c.logger.Error(
						ctx,
						"message routing error",
						map[string]any{
							"error": err,
						},
					)
				}
			}()
		}
	}
}

func (c *Consumer) Close(ctx context.Context) error {
	c.logger.Info(
		ctx,
		"closing consumer and committing offsets...",
		nil,
	)

	return c.consumer.Close()
}

type KafkaConsumerHeadersCarrier struct {
	Headers []kafka.Header
}

func (c *KafkaConsumerHeadersCarrier) Get(key string) string {
	for _, h := range c.Headers {
		if h.Key == key {
			return string(h.Value)
		}
	}

	return ""
}

func (c *KafkaConsumerHeadersCarrier) Set(key string, value string) {}

func (c *KafkaConsumerHeadersCarrier) Keys() []string {
	keys := make([]string, len(c.Headers))
	for i, h := range c.Headers {
		keys[i] = h.Key
	}

	return keys
}
